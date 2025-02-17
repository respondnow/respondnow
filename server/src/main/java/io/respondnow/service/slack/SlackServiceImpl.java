package io.respondnow.service.slack;

import static io.respondnow.model.incident.ChannelStatus.Operational;

import com.slack.api.Slack;
import com.slack.api.app_backend.views.payload.ViewSubmissionPayload;
import com.slack.api.bolt.App;
import com.slack.api.bolt.AppConfig;
import com.slack.api.bolt.context.builtin.GlobalShortcutContext;
import com.slack.api.bolt.request.builtin.GlobalShortcutRequest;
import com.slack.api.bolt.request.builtin.ViewSubmissionRequest;
import com.slack.api.bolt.socket_mode.SocketModeApp;
import com.slack.api.methods.SlackApiException;
import com.slack.api.methods.request.conversations.ConversationsCreateRequest;
import com.slack.api.methods.request.conversations.ConversationsInviteRequest;
import com.slack.api.methods.response.auth.AuthTestResponse;
import com.slack.api.methods.response.chat.ChatPostMessageResponse;
import com.slack.api.methods.response.conversations.ConversationsCreateResponse;
import com.slack.api.methods.response.conversations.ConversationsInviteResponse;
import com.slack.api.methods.response.users.UsersInfoResponse;
import com.slack.api.methods.response.views.ViewsOpenResponse;
import com.slack.api.model.Conversation;
import com.slack.api.model.block.*;
import com.slack.api.model.block.Blocks;
import com.slack.api.model.block.DividerBlock;
import com.slack.api.model.block.InputBlock;
import com.slack.api.model.block.LayoutBlock;
import com.slack.api.model.block.SectionBlock;
import com.slack.api.model.block.composition.MarkdownTextObject;
import com.slack.api.model.block.composition.OptionObject;
import com.slack.api.model.block.composition.PlainTextObject;
import com.slack.api.model.block.composition.TextObject;
import com.slack.api.model.block.element.BlockElements;
import com.slack.api.model.block.element.ButtonElement;
import com.slack.api.model.block.element.StaticSelectElement;
import com.slack.api.model.block.element.UsersSelectElement;
import com.slack.api.model.event.AppHomeOpenedEvent;
import com.slack.api.model.event.AppMentionEvent;
import com.slack.api.model.event.MemberJoinedChannelEvent;
import com.slack.api.model.view.*;
import com.slack.api.socket_mode.SocketModeClient;
import io.respondnow.dto.incident.CreateRequest;
import io.respondnow.exception.RoleUpdateException;
import io.respondnow.model.incident.*;
import io.respondnow.model.user.UserDetails;
import io.respondnow.service.incident.IncidentService;
import java.io.IOException;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import lombok.extern.slf4j.Slf4j;
import org.apache.logging.log4j.util.Strings;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.stereotype.Service;

@Service
@Slf4j
public class SlackServiceImpl implements SlackService {
  private static final Logger logger = LoggerFactory.getLogger(SlackServiceImpl.class);
  private final Slack slackClient;
  private final SocketModeClient socketModeSlackClient;
  private final App slackApp;
  private final SocketModeApp socketModeApp;
  private final ExecutorService executorService;
  @Autowired private IncidentService incidentService;
  private String botUserId;

  @Value("${slack.botToken}")
  private String botToken;

  @Value("${slack.appToken}")
  private String appToken;

  @Value("${slack.incidentChannelID}")
  private String incidentChannelID;

  @Value("${slack.maxRetryAttempts:3}")
  private int maxRetryAttempts;

  public SlackServiceImpl(
      @Value("${slack.botToken}") String botToken,
      @Value("${slack.appToken}") String appToken,
      @Value("${slack.incidentChannelID}") String incidentChannelID)
      throws Exception {
    if (botToken == null || appToken == null || incidentChannelID == null) {
      throw new IllegalArgumentException(
          "Bot token, App token, and Incident channel ID must not be null");
    }
    this.botToken = botToken;
    this.appToken = appToken;
    this.incidentChannelID = incidentChannelID;

    AppConfig appConfig = AppConfig.builder().singleTeamBotToken(botToken).build();
    this.slackApp = new App(appConfig);
    this.socketModeApp = new SocketModeApp(this.appToken, this.slackApp);
    this.socketModeApp.getClient();
    this.slackClient = Slack.getInstance();
    this.socketModeSlackClient = socketModeApp.getClient();
    this.executorService = Executors.newCachedThreadPool();

    registerEventHandlers();
    registerShortcutHandlers();
    registerViewSubmissionHandlers();
    registerBlockActionHandlers();
  }

  private static Map<String, Object> errorResponse() {
    Map<String, Object> response = new HashMap<>();
    response.put("text", "Failed to open the incident modal. Please try again later.");
    return response;
  }

  @Override
  public String getIncidentChannelID() {
    return incidentChannelID;
  }

  @Override
  public String getBotToken() {
    return botToken;
  }

  @Override
  public String getAppToken() {
    return appToken;
  }

  @Override
  public Slack getSlackClient() {
    return slackClient;
  }

  @Override
  public App getSlackApp() {
    return slackApp;
  }

  /** Register all event handlers for the Slack app. */
  private void registerEventHandlers() throws RuntimeException {
    try {
      registerAppHomeOpenedEvent();
      registerAppMentionEvent();
      registerMemberJoinedChannelEvent();
    } catch (RuntimeException e) {
      throw new RuntimeException(e);
    }
  }

  /** Register all shortcut handlers for the Slack app. */
  private void registerShortcutHandlers() throws RuntimeException {
    try {
      registerListOpenIncidentsShortcut();
      registerListClosedIncidentsShortcut();
      registerCreateIncidentShortcut();
    } catch (RuntimeException e) {
      throw new RuntimeException(e);
    }
  }

  /** Register view submission handlers for the Slack app. */
  private void registerViewSubmissionHandlers() throws RuntimeException {
    try {
      handleCreateIncidentViewSubmission();
      handleIncidentSummaryViewSubmission();
      handleIncidentCommentViewSubmission();
      handleIncidentRolesViewSubmission();
      handleIncidentStatusViewSubmission();
      handleIncidentSeverityViewSubmission();
    } catch (RuntimeException e) {
      throw new RuntimeException(e);
    }
  }

  /** Handles the submission of the "Create Incident" modal. */
  private void handleCreateIncidentViewSubmission() {
    slackApp.viewSubmission(
        "create_incident_modal",
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.info("Received create incident view submission: {}", payload);
                  createIncident(payload);
                  logger.info("Successfully created incident for payload: {}", payload);
                } catch (Exception e) {
                  logger.error("Error creating incident for payload: {}", payload, e);
                  // Optionally, handle the error (e.g., notify the user)
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /** Handles the submission of the "Incident Summary" modal. */
  public void handleIncidentSummaryViewSubmission() {
    slackApp.viewSubmission(
        "incident_summary_modal",
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.debug("Update summary received: {}", payload);
                  handleIncidentSummaryViewSubmission(payload);
                  logger.debug("Successfully updated incident summary for payload: {}", payload);
                } catch (Exception e) {
                  logger.error("Error updating incident summary for payload: {}", payload, e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /** Handles the submission of the "Incident Comment" modal. */
  private void handleIncidentCommentViewSubmission() {
    slackApp.viewSubmission(
        "incident_comment_modal",
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.debug("A new comment received: {}", payload);
                  handleIncidentCommentViewSubmission(payload);
                  logger.debug("Successfully processed incident comment for payload: {}", payload);
                } catch (Exception e) {
                  logger.error("Error processing incident comment for payload: {}", payload, e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /** Handles the submission of the "Incident Roles" modal. */
  private void handleIncidentRolesViewSubmission() {
    slackApp.viewSubmission(
        "incident_roles_modal",
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.debug("Incident roles modal received: {}", payload);
                  handleIncidentRolesViewSubmission(payload);
                  logger.debug("Successfully processed incident roles for payload: {}", payload);
                } catch (Exception e) {
                  logger.error("Error processing incident roles for payload: {}", payload, e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /** Handles the submission of the "Incident Status" modal. */
  private void handleIncidentStatusViewSubmission() {
    slackApp.viewSubmission(
        "incident_status_modal",
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.debug("Incident status modal received: {}", payload);
                  handleIncidentStatusViewSubmission(payload);
                  logger.debug("Successfully processed incident status for payload: {}", payload);
                } catch (Exception e) {
                  logger.error("Error processing incident status for payload: {}", payload, e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /** Handles the submission of the "Incident Severity" modal. */
  private void handleIncidentSeverityViewSubmission() {
    slackApp.viewSubmission(
        "incident_severity_modal",
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.debug("Incident severity modal received: {}", payload);
                  handleIncidentSeverityViewSubmission(payload);
                  logger.debug("Successfully processed incident severity for payload: {}", payload);
                } catch (Exception e) {
                  logger.error("Error processing incident severity for payload: {}", payload, e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /** Register all block action handlers for the Slack app. */
  private void registerBlockActionHandlers() {
    try {
      registerCreateIncidentChannelJoinButton();
      registerCreateIncidentButton();
      registerUpdateIncidentSummaryButton();
      registerUpdateIncidentCommentButton();
      registerUpdateIncidentAssignRolesButton();
      registerUpdateIncidentStatusButton();
      registerUpdateIncidentSeverityButton();
      registerViewIncidentActionHandler();
    } catch (RuntimeException e) {
      throw new RuntimeException(e);
    }
  }

  private void registerCreateIncidentChannelJoinButton() {
    slackApp.blockAction(
        "create_incident_channel_join_channel_button",
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  String value = req.getPayload().getActions().get(0).getValue();
                  logger.info("Button clicked with value: {}", value);

                  if (req.getPayload().getResponseUrl() != null) {
                    ctx.respond(
                        r -> r.text("You've sent \"" + value + "\" by clicking the button!"));
                  }

                  // Additional processing logic can be added here
                } catch (Exception e) {
                  logger.error("Error handling create incident channel join button", e);
                  // Optionally, handle the error (e.g., notify the user via Slack)
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  private void registerUpdateIncidentSummaryButton() {
    slackApp.blockAction(
        "update_incident_summary_button",
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  // Build the modal
                  View modalRequest =
                      View.builder()
                          .type("modal")
                          .privateMetadata(req.getPayload().getActions().get(0).getValue())
                          .callbackId("incident_summary_modal")
                          .title(
                              ViewTitle.builder()
                                  .type("plain_text")
                                  .text("Update Incident Summary")
                                  .build())
                          .blocks(
                              Collections.singletonList(
                                  getSummaryBlock(
                                      "create_incident_modal_summary",
                                      "create_incident_modal_set_summary")))
                          .submit(ViewSubmit.builder().type("plain_text").text("Submit").build())
                          .build();
                  ctx.client()
                      .viewsOpen(
                          r -> r.triggerId(req.getPayload().getTriggerId()).view(modalRequest));

                  logger.info("Opened Incident Summary modal for payload: {}", req.getPayload());
                } catch (Exception e) {
                  logger.error("Error opening Incident Summary modal", e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  private void registerUpdateIncidentCommentButton() {
    slackApp.blockAction(
        "update_incident_comment_button",
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  View modalRequest =
                      View.builder()
                          .type("modal")
                          .privateMetadata(req.getPayload().getActions().get(0).getValue())
                          .callbackId("incident_comment_modal")
                          .title(
                              ViewTitle.builder()
                                  .type("plain_text")
                                  .text("Update Incident Comment")
                                  .build())
                          .blocks(
                              Collections.singletonList(
                                  getCommentBlock(
                                      "update_incident_modal_comment",
                                      "update_incident_modal_set_comment")))
                          .submit(ViewSubmit.builder().type("plain_text").text("Submit").build())
                          .build();
                  ctx.client()
                      .viewsOpen(
                          r -> r.triggerId(req.getPayload().getTriggerId()).view(modalRequest));

                  logger.info("Opened Incident Comment modal for payload: {}", req.getPayload());
                } catch (Exception e) {
                  logger.error("Error opening Incident Comment modal", e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  private void registerUpdateIncidentAssignRolesButton() {
    slackApp.blockAction(
        "update_incident_assign_roles_button",
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  View modalRequest =
                      View.builder()
                          .type("modal")
                          .privateMetadata(req.getPayload().getActions().get(0).getValue())
                          .callbackId("incident_roles_modal")
                          .title(
                              ViewTitle.builder()
                                  .type("plain_text")
                                  .text("Assign Incident Roles")
                                  .build())
                          .blocks(getUpdateRoleBlock())
                          .submit(ViewSubmit.builder().type("plain_text").text("Submit").build())
                          .close(ViewClose.builder().type("plain_text").text("Close").build())
                          .build();
                  ctx.client()
                      .viewsOpen(
                          r -> r.triggerId(req.getPayload().getTriggerId()).view(modalRequest));

                  logger.info("Opened Incident Roles modal for payload: {}", req.getPayload());
                } catch (Exception e) {
                  logger.error("Error opening Incident Roles modal", e);
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  private void registerUpdateIncidentStatusButton() {
    slackApp.blockAction(
        "update_incident_status_button",
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  View modalRequest =
                      View.builder()
                          .type("modal")
                          .privateMetadata(req.getPayload().getActions().get(0).getValue())
                          .callbackId("incident_status_modal")
                          .title(
                              ViewTitle.builder()
                                  .type("plain_text")
                                  .text("Update Incident Status")
                                  .build())
                          .blocks(Collections.singletonList(updateStatus()))
                          .submit(ViewSubmit.builder().type("plain_text").text("Submit").build())
                          .build();
                  ctx.client()
                      .viewsOpen(
                          r -> r.triggerId(req.getPayload().getTriggerId()).view(modalRequest));

                  logger.info("Opened Incident Status modal for payload: {}", req.getPayload());
                } catch (Exception e) {
                  logger.error("Error opening Incident Status modal", e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  private void registerUpdateIncidentSeverityButton() {
    slackApp.blockAction(
        "update_incident_severity_button",
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  View modalRequest =
                      View.builder()
                          .type("modal")
                          .privateMetadata(req.getPayload().getActions().get(0).getValue())
                          .callbackId("incident_severity_modal")
                          .title(
                              ViewTitle.builder()
                                  .type("plain_text")
                                  .text("Update Incident Severity")
                                  .build())
                          .blocks(Collections.singletonList(getSeverityBlock()))
                          .submit(ViewSubmit.builder().type("plain_text").text("Submit").build())
                          .build();

                  ctx.client()
                      .viewsOpen(
                          r -> r.triggerId(req.getPayload().getTriggerId()).view(modalRequest));

                  logger.info("Opened Incident Severity modal for payload: {}", req.getPayload());
                } catch (Exception e) {
                  logger.error("Error opening Incident Severity modal", e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  private void registerViewIncidentActionHandler() {
    String regex = "^view_incident.*";
    Pattern pattern = Pattern.compile(regex);

    slackApp.blockAction(
        pattern,
        (req, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  String actionId = req.getPayload().getActions().get(0).getActionId();
                  String[] parts = actionId.split("_");
                  if (parts.length < 3) {
                    logger.warn("Invalid actionId format: {}", actionId);
                    return;
                  }
                  String incidentIdentifier = parts[2];
                  Incident incident = incidentService.getIncidentByIdentifier(incidentIdentifier);

                  if (incident == null) {
                    logger.warn("Incident not found with identifier: {}", incidentIdentifier);
                    // Optionally, notify the user via Slack
                    return;
                  }

                  View modalRequest =
                      View.builder()
                          .type("modal")
                          .privateMetadata(req.getPayload().getActions().get(0).getValue())
                          .callbackId("incident_details_modal")
                          .title(
                              ViewTitle.builder()
                                  .type("plain_text")
                                  .text("🚨 Incident Details")
                                  .emoji(false)
                                  .build())
                          .blocks(getViewDetailLayoutBlock(incident))
                          .submit(ViewSubmit.builder().type("plain_text").text("Submit").build())
                          .close(ViewClose.builder().type("plain_text").text("Close").build())
                          .build();
                  ctx.client()
                      .viewsPush(
                          r -> r.triggerId(req.getPayload().getTriggerId()).view(modalRequest));

                  logger.info("Opened Incident Details modal for incident: {}", incidentIdentifier);
                } catch (Exception e) {
                  logger.error("Error opening Incident Details modal", e);
                  // Optionally, handle the error
                }
              });

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  /**
   * Generates a list of LayoutBlocks detailing the incident information.
   *
   * @param incident The incident object containing details.
   * @return A list of LayoutBlocks representing the incident details.
   */
  private List<LayoutBlock> getViewDetailLayoutBlock(Incident incident) {
    // Validate the incident object
    if (incident == null) {
      logger.error("Incident object is null.");
      throw new IllegalArgumentException("Incident cannot be null.");
    }

    List<LayoutBlock> layoutBlocks = new ArrayList<>();

    try {
      // Header
      layoutBlocks.add(SlackBlockFactory.createSectionBlock("*Incident Details*", ""));

      // Name
      String nameText =
          String.format(
              ":writing_hand: *Name:* %s", Optional.ofNullable(incident.getName()).orElse("N/A"));
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(nameText, ""));

      // Severity
      String severityText =
          String.format(
              ":vertical_traffic_light: *Severity:* %s",
              Optional.ofNullable(incident.getSeverity())
                  .orElse(Severity.valueOf(String.valueOf(Severity.SEV2))));
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(severityText, ""));

      // Current Status
      String statusText =
          String.format(
              ":eyes: *Current Status:* %s",
              Optional.ofNullable(incident.getStatus())
                  .map(Enum::name)
                  .orElse(String.valueOf(Status.Started)));
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(statusText, ""));

      // Incident Commander
      String incidentCommander = extractUserMention(incident, RoleType.Incident_Commander);
      String commanderText =
          String.format(
              ":firefighter: *Commander:* %s",
              incidentCommander.isEmpty() ? "N/A" : incidentCommander);
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(commanderText, ""));

      // Communications Lead
      String communicationsLead = extractUserMention(incident, RoleType.Communications_Lead);
      String communicationsLeadText =
          String.format(
              ":phone: *Communications Lead:* %s",
              communicationsLead.isEmpty() ? "N/A" : communicationsLead);
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(communicationsLeadText, ""));

      // Summary
      String summaryText =
          String.format(
              ":open_book: *Summary:* %s",
              Optional.ofNullable(incident.getSummary()).orElse("N/A"));
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(summaryText, ""));

      // Started At
      String startedAtText =
          String.format(":clock1: *Started At:* %s", formatUnixTimestamp(incident.getCreatedAt()));
      layoutBlocks.add(SlackBlockFactory.createSectionBlock(startedAtText, ""));

      // Completed At (if resolved)
      if (incident.getStatus() == Status.Resolved && incident.getUpdatedAt() != null) {
        String completedAtText =
            String.format(
                ":checkered_flag: *Completed At:* %s",
                formatUnixTimestamp(incident.getUpdatedAt()));
        layoutBlocks.add(SlackBlockFactory.createSectionBlock(completedAtText, ""));
      }

    } catch (Exception e) {
      logger.error("Error while building layout blocks for incident: {}", incident, e);
      // Depending on your application's needs, you might want to rethrow or handle differently
      throw new RuntimeException("Failed to build incident details layout.", e);
    }

    return layoutBlocks;
  }

  /**
   * Extracts and formats the user mention for a given role type.
   *
   * @param incident The incident containing roles.
   * @param roleType The role type to extract.
   * @return The formatted user mention or empty string if not found.
   */
  private String extractUserMention(Incident incident, RoleType roleType) {
    if (incident.getRoles() == null || incident.getRoles().isEmpty()) {
      logger.warn("Incident roles are null or empty.");
      return "";
    }

    Optional<Role> roleOpt =
        incident.getRoles().stream()
            .filter(role -> roleType.equals(role.getRoleType()))
            .findFirst();

    if (roleOpt.isPresent()) {
      String userId =
          Optional.ofNullable(roleOpt.get().getUserDetails())
              .map(UserDetails::getUserId)
              .orElse("");
      if (!userId.isEmpty()) {
        return String.format("<@%s>", userId);
      }
    }

    logger.info("Role '{}' not found or userId is empty.", roleType);
    return "";
  }

  /**
   * Formats a Unix timestamp (in seconds) to a human-readable date string.
   *
   * @param unixTimestamp The Unix timestamp in seconds.
   * @return The formatted date string or "Invalid Date" if timestamp is invalid.
   */
  private String formatUnixTimestamp(long unixTimestamp) {
    try {
      Date date = new Date(unixTimestamp * 1000L);
      return date.toString(); // Customize the format as needed
    } catch (Exception e) {
      logger.error("Invalid Unix timestamp: {}", unixTimestamp, e);
      return "Invalid Date";
    }
  }

  private void registerListClosedIncidentsShortcut() throws RuntimeException {
    // Handle the "list_closed_incidents_modal" shortcut
    slackApp.globalShortcut(
        "list_closed_incidents_modal",
        (req, ctx) -> {
          ctx.logger.info("Shortcut received: {}", req.getPayload());

          try {
            listIncidents(ctx, SlackIncidentType.Closed);
            return ctx.ack();
          } catch (Exception e) {
            ctx.logger.error("Error listing closed incidents: {}", e.getMessage(), e);
            // Respond with a JSON error message
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("text", "Failed to list closed incidents.");
            return ctx.ackWithJson(errorResponse);
          }
        });
  }

  private void registerListOpenIncidentsShortcut() throws RuntimeException {
    // Handle the "list_open_incidents_modal" shortcut
    slackApp.globalShortcut(
        "list_open_incidents_modal",
        (req, ctx) -> {
          ctx.logger.info("Shortcut list_open_incidents_modal received: {}", req.getPayload());

          try {
            listIncidents(ctx, SlackIncidentType.Open);
            return ctx.ack();
          } catch (Exception e) {
            ctx.logger.error("Error listing open incidents: {}", e.getMessage(), e);
            // Respond with a JSON error message
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("text", "Failed to list open incidents.");
            return ctx.ackWithJson(errorResponse);
          }
        });
  }

  private void registerCreateIncidentShortcut() throws RuntimeException {
    // Handle the "open_incident_modal" shortcut
    slackApp.globalShortcut(
        "open_incident_modal",
        (req, ctx) -> {
          // Log the shortcut payload
          ctx.logger.info("Shortcut open_incident_modal received: {}", req.getPayload());

          // Call the service to create an incident modal view
          try {
            createIncident(req, ctx);
            return ctx.ack();
          } catch (Exception e) {
            ctx.logger.error("Error creating incident modal: {}", e.getMessage(), e);
            // Respond with a JSON error message
            Map<String, Object> errorResponse = new HashMap<>();
            errorResponse.put("text", "Failed to open incident modal. Please try again later.");
            return ctx.ackWithJson(errorResponse);
          }
        });
  }

  private void registerAppHomeOpenedEvent() throws RuntimeException {
    slackApp.event(
        AppHomeOpenedEvent.class,
        (payload, ctx) -> {
          AppHomeOpenedEvent event = payload.getEvent();
          logger.info("Received AppHomeOpenedEvent: userId={}", event.getUser());

          // Handle the event asynchronously
          executorService.submit(
              () -> {
                try {
                  handleAppHome(event);
                } catch (InterruptedException e) {
                  throw new RuntimeException(e);
                }
              });
          return ctx.ack(); // Acknowledge the event
        });
  }

  private void registerAppMentionEvent() throws RuntimeException {
    slackApp.event(
        AppMentionEvent.class,
        (payload, ctx) -> {
          AppMentionEvent event = payload.getEvent();
          logger.info("App mentioned in channel: {}", event.getChannel());

          // Handle the event asynchronously
          executorService.submit(
              () -> {
                try {
                  handleAppMention(event);
                } catch (InterruptedException e) {
                  throw new RuntimeException(e);
                }
              });
          return ctx.ack(); // Acknowledge the event
        });
  }

  private void registerMemberJoinedChannelEvent() throws RuntimeException {
    slackApp.event(
        MemberJoinedChannelEvent.class,
        (payload, ctx) -> {
          MemberJoinedChannelEvent event = payload.getEvent();
          logger.info("New member joined channel: {}", event.getChannel());

          // Handle the event asynchronously
          executorService.submit(
              () -> {
                try {
                  handleMemberJoinedChannelEvent(event);
                } catch (InterruptedException e) {
                  throw new RuntimeException(e);
                }
              });
          return ctx.ack(); // Acknowledge the event
        });
  }

  @Override
  public void startApp() {
    // Ensure the Slack app is started asynchronously
    executorService.submit(
        () -> {
          try {
            Runtime.getRuntime().addShutdownHook(new Thread(this::shutdownApp));
            logger.info("Starting Slack app in Socket Mode...");
            socketModeApp.start(); // This is the blocking call, but it will run asynchronously
          } catch (Exception e) {
            logger.error("Error starting Slack app in Socket Mode", e);
            throw new RuntimeException("Failed to start Slack app", e);
          }
        });
  }

  /** Stop the Slack app and release resources. */
  public void shutdownApp() {
    try {
      logger.info("Shutting down Slack app...");
      if (executorService != null && !executorService.isShutdown()) {
        executorService.shutdown();
        if (!executorService.awaitTermination(5, TimeUnit.SECONDS)) {
          executorService.shutdownNow();
        }
      }
      socketModeApp.stop();
      logger.info("Slack app shut down successfully");
    } catch (Exception e) {
      logger.error("Error during Slack app shutdown", e);
    }
  }

  @Override
  public void setBotUserIDAndName() throws Exception {
    AuthTestResponse authTestResponse =
        slackClient.methods(botToken).authTest(r -> r.token(botToken));
    if (authTestResponse.isOk()) {
      botUserId = authTestResponse.getUserId();
      String botUserName = authTestResponse.getUser();
      logger.info("Bot User ID: {}, Bot User Name: {}", botUserId, botUserName);
    } else {
      throw new Exception("Failed to fetch bot details: " + authTestResponse.getError());
    }
  }

  @Override
  public String getBotUserId() throws Exception {
    if (botUserId == null) {
      setBotUserIDAndName();
    }
    return botUserId;
  }

  @Override
  public void addBotUserToIncidentChannel(String botUserID, String channelID) throws Exception {
    if (!isBotInChannel(botUserID, channelID)) {
      // Ensure the bot is a member of the channel
      var joinResponse = slackClient.methods(botToken).conversationsJoin(r -> r.channel(channelID));
      if (!joinResponse.isOk()) {
        throw new Exception("Failed to join the channel: " + joinResponse.getError());
      }
      logger.info("Bot: {} successfully added to the channel: {}", botUserID, channelID);
    } else {
      logger.info("Bot: {} is already in the channel: {}", botUserID, channelID);
    }
  }

  public boolean isBotInChannel(String botUserID, String channelID) throws Exception {
    logger.info("Checking if bot is in the channel...");

    List<String> members = listAllMembersOfChannel(channelID);

    if (members == null || members.isEmpty()) {
      logger.info("No members found in the channel.");
      return false; // Bot cannot be in a channel with no members
    }

    // Normalize and trim spaces for comparison
    botUserID = botUserID.trim(); // Remove leading/trailing spaces
    members =
        members.stream()
            .map(String::trim) // Trim each member ID
            .collect(Collectors.toList());

    // Check if the bot user ID is in the list of channel members
    return members.contains(botUserID);
  }

  public List<String> listAllMembersOfChannel(String channelId) throws Exception {
    List<String> allMembers = new ArrayList<>();
    String cursor = null;

    do {
      String finalCursor = cursor;
      var response =
          slackClient
              .methods(botToken)
              .conversationsMembers(r -> r.channel(channelId).cursor(finalCursor));

      if (!response.isOk()) {
        throw new Exception("Failed to list members for channel: " + response.getError());
      }

      // Add members from the current page to the list
      allMembers.addAll(response.getMembers());

      // Get the next cursor for pagination
      cursor =
          response.getResponseMetadata() != null
              ? response.getResponseMetadata().getNextCursor()
              : null;

    } while (cursor != null && !cursor.isEmpty()); // Continue if there's more data to fetch

    return allMembers;
  }

  @Override
  public List<String> listUsers(String channelID) throws Exception {
    return slackClient
        .methods(botToken)
        .conversationsMembers(r -> r.channel(channelID))
        .getMembers();
  }

  @Override
  public List<Conversation> listChannels() throws Exception {
    List<Conversation> allChannels = new ArrayList<>();
    String cursor = null;

    do {
      String finalCursor = cursor;
      var response =
          slackClient
              .methods(botToken)
              .conversationsList(r -> r.limit(100).excludeArchived(true).cursor(finalCursor));

      if (!response.isOk()) {
        throw new Exception("Failed to list channels: " + response.getError());
      }

      allChannels.addAll(response.getChannels());
      cursor =
          response.getResponseMetadata() != null
              ? response.getResponseMetadata().getNextCursor()
              : null;
    } while (cursor != null && !cursor.isEmpty());

    return allChannels;
  }

  public void handleMemberJoinedChannelEvent(MemberJoinedChannelEvent event)
      throws InterruptedException {
    if (event == null) {
      logger.error("Null event received for MemberJoinedChannelEvent");
      return;
    }
    logger.info(
        "Member Joined Channel Event for user: {} in channel: {}",
        event.getUser(),
        event.getChannel());
  }

  public void handleAppMention(AppMentionEvent event) throws InterruptedException {
    if (event == null) {
      logger.error("Null event received for AppMentionEvent");
      return;
    }
    logger.info(
        "App Mention Event in channel: {} by user: {}", event.getChannel(), event.getUser());

    boolean success = false;
    int attempt = 0;

    // Retry logic for publishing the home view
    while (!success && attempt < maxRetryAttempts) {
      try {
        // Respond with a message
        String channelId = event.getChannel();
        String message =
            "Hi <@"
                + event.getUser()
                + "> :wave:, how can I assist you? Please use available shortcuts to interact with me!";
        slackApp.client().chatPostMessage(r -> r.channel(channelId).text(message));
        success = true;
      } catch (IOException | SlackApiException e) {
        // Retry logic for rate-limited requests
        if ("too_many_requests".equals(e.getMessage())) {
          logger.warn("App mention events rate-limited by slack. Retrying...");
          TimeUnit.SECONDS.sleep(5);
        } else {
          // Handle Slack API error
          logger.error("Slack API error: {}", e.getMessage(), e);
          break;
        }
      } catch (Exception e) {
        // General exception handling
        logger.error("Unexpected error: {}", e.getMessage(), e);
        break;
      }
      attempt++;
    }

    if (!success) {
      logger.error("Failed to respond to app mention event after {} attempts", maxRetryAttempts);
    }
  }

  @Override
  public void handleAppHome(AppHomeOpenedEvent event) throws InterruptedException {
    if (event == null) {
      logger.error("Null event received for AppHomeOpenedEvent");
      return;
    }

    String userId = event.getUser();
    List<LayoutBlock> slackBlocks = createSlackBlocks(userId, botUserId);

    // Create the View object to send to Slack
    View view = Views.view(v -> v.type("home").blocks(slackBlocks));

    boolean success = false;
    int attempt = 0;

    // Retry logic for publishing the home view
    while (!success && attempt < maxRetryAttempts) {
      try {
        // Publish the app home view
        slackApp.client().viewsPublish(r -> r.userId(userId).view(view));
        logger.info("App home view published successfully for user: {}", userId);
        success = true;
      } catch (IOException | SlackApiException e) {
        // Retry logic for rate-limited requests
        if ("too_many_requests".equals(e.getMessage())) {
          logger.warn("App home events rate-limited by Slack. Retrying...");
          TimeUnit.SECONDS.sleep(5);
        } else {
          // Handle Slack API error
          logger.error("Slack API error: {}", e.getMessage(), e);
          break;
        }
      } catch (Exception e) {
        // General exception handling
        logger.error("Unexpected error: {}", e.getMessage(), e);
        break;
      }
      attempt++;
    }

    if (!success) {
      logger.error("Failed to publish app home view after {} attempts", maxRetryAttempts);
    }
  }

  private List<LayoutBlock> createSlackBlocks(String userId, String botUserId) {
    // Modularize the block creation logic for better readability
    return Arrays.asList(
        SlackBlockFactory.createHeaderBlock(":robot_face: Respond Now", "app_home_resp_header"),
        // TODO: Implement this
        //        SlackBlockFactory.createActionsBlock(
        //            "app_home_resp_create_incident_button",
        //            ButtonElement.builder()
        //                .text(PlainTextObject.builder().text("Start New
        // Incident").emoji(true).build())
        //                .actionId("create_incident_modal")
        //                .value("show_incident_modal")
        //                .style("danger")
        //                .build()),
        SlackBlockFactory.createDividerBlock(),
        SlackBlockFactory.createSectionBlock(
            "*Hi there, <@" + userId + "> :wave:*\n\nI'm your friendly Respond Now...",
            "app_home_resp_intro"),
        SlackBlockFactory.createHeaderBlock(
            ":slack: Adding me to a channel", "app_home_resp_add_to_channel_header"),
        SlackBlockFactory.createDividerBlock(),
        SlackBlockFactory.createSectionBlock(
            "To add me to a new channel, please use <@" + botUserId + ">",
            "app_home_resp_add_to_channel_steps"),
        SlackBlockFactory.createHeaderBlock(
            ":firefighter: Creating New Incidents", "app_home_resp_creating_new_incidents_header"),
        SlackBlockFactory.createDividerBlock(),
        SlackBlockFactory.createSectionBlock(
            "To create a new incident, you can do the following:\n"
                + "- Use the 'Start New Incident' button here\n "
                + "- Search for 'start a new incident' in the slack search bar\n"
                + "- Type _/start_ in an incident slack channel to find my create command and run it.",
            "app_home_resp_create_incident_steps"),
        SlackBlockFactory.createHeaderBlock(
            ":point_right: Documentation and Learning Materials", "app_home_resp_docs_header"),
        SlackBlockFactory.createDividerBlock(),
        SlackBlockFactory.createSectionBlock(
            "I have a lot of features. To check them all out, visit my <https://github.com/respondnow/respondnow/blob/main/README.md|docs>.",
            "app_home_resp_docs_content"));
  }

  public void createIncident(ViewSubmissionRequest payload) {
    Map<String, Map<String, ViewState.Value>> stateValues =
        Optional.ofNullable(payload.getPayload())
            .map(ViewSubmissionPayload::getView)
            .map(View::getState)
            .map(ViewState::getValues)
            .orElse(Collections.emptyMap());

    // Safely extract values from the payload
    String incidentIdentifier =
        Optional.ofNullable(payload.getPayload().getView().getPrivateMetadata())
            .orElse("Unknown Incident Identifier");

    String name =
        Optional.ofNullable(stateValues.get("create_incident_modal_name"))
            .map(inner -> inner.get("create_incident_modal_set_name"))
            .map(ViewState.Value::getValue)
            .orElse(null);

    if (name == null) {
      logger.error("Name field is missing in view submission payload.");
    }

    String incidentType =
        Optional.ofNullable(stateValues.get("incident_type"))
            .map(inner -> inner.get("create_incident_modal_set_incident_type"))
            .map(ViewState.Value::getSelectedOption)
            .map(ViewState.SelectedOption::getValue)
            .orElse(null);

    if (incidentType == null) {
      logger.error("Incident type field is missing in view submission payload.");
    }

    String summary =
        Optional.ofNullable(stateValues.get("create_incident_modal_summary"))
            .map(inner -> inner.get("create_incident_modal_set_summary"))
            .map(ViewState.Value::getValue)
            .orElse(null);

    if (summary == null) {
      logger.error("Summary field is missing in view submission payload.");
    }

    String severity =
        Optional.ofNullable(stateValues.get("incident_severity"))
            .map(inner -> inner.get("create_incident_modal_set_incident_severity"))
            .map(ViewState.Value::getSelectedOption)
            .map(ViewState.SelectedOption::getValue)
            .orElse(null);

    if (severity == null) {
      logger.error("Severity field is missing in view submission payload.");
    }

    String responseChannel =
        Optional.ofNullable(stateValues.get("create_incident_modal_conversation_select"))
            .map(inner -> inner.get("create_incident_modal_select_conversation"))
            .map(ViewState.Value::getSelectedConversation)
            .orElse(null);

    if (responseChannel == null) {
      logger.error("Response channel field is missing in view submission payload.");
    }

    // Get user details (assume you have a method to fetch UserDetails by
    // userId)
    UserDetails userDetails = null;
    try {
      userDetails =
          fetchSlackUserDetails(payload.getPayload().getUser().getId(), ChannelSource.Slack);
    } catch (SlackApiException | IOException e) {
      throw new RuntimeException(e);
    }

    UserDetails finalUserDetails = userDetails;
    @NotNull
    List<Role> roles =
        Optional.ofNullable(stateValues.get("incident_role"))
            .map(inner -> inner.get("create_incident_modal_set_incident_role"))
            .map(ViewState.Value::getSelectedOptions) // Get the list of selected options
            .map(
                selectedOptions ->
                    selectedOptions.stream()
                        .map(
                            ViewState.SelectedOption
                                ::getValue) // Extract the value of each selected option
                        .map(
                            roleString -> {
                              // Map string to RoleType enum
                              RoleType roleType = RoleType.valueOf(roleString);

                              // Create Role object
                              return new Role(roleType, finalUserDetails);
                            })
                        .collect(Collectors.toList())) // Collect the values into a List
            .orElse(Collections.emptyList()); // Return an empty list if no roles are found

    if (roles.isEmpty()) {
      logger.error("No roles are selected in view submission payload.");
    } else {
      logger.info("Selected Roles: {}", roles);
    }

    // Logging extracted values
    logger.info(
        "Incident Identifier: {}, Name: {}, Incident Type: {}, Summary: {}, Severity: {}, Response Channel: {}",
        incidentIdentifier,
        name,
        incidentType,
        summary,
        severity,
        responseChannel);

    // Proceed with further business logic if all required fields are present
    if (name != null
        && incidentType != null
        && summary != null
        && severity != null
        && !roles.isEmpty()
        && responseChannel != null) {
      logger.info("Creating incident with the provided details.");

      long createdAt = Instant.now().getEpochSecond();
      String incidentId = incidentService.generateIncidentIdentifier(createdAt);

      Slack slackClient = getSlackClient();
      // Create Slack channel using conversations.create

      // Sanitize the name to remove invalid characters
      String sanitizedChannelName = sanitizeChannelName("inc-" + name + "-" + createdAt);

      try {
        ConversationsCreateRequest createChannelRequest =
            ConversationsCreateRequest.builder()
                .name(sanitizedChannelName) // Channel name
                .isPrivate(false) // Change to true if you want a private channel
                .teamId(payload.getPayload().getTeam().getId()) // Team ID from the callback
                .build();

        ConversationsCreateResponse createChannelResponse =
            slackClient.methods(botToken).conversationsCreate(createChannelRequest);

        if (!createChannelResponse.isOk()) {
          throw new IOException(
              "Failed to create Slack channel: " + createChannelResponse.getError());
        }

        String channelId = createChannelResponse.getChannel().getId();
        log.info("Successfully created an incident channel: {}", channelId);

        // Invite users to the channel
        try {
          ConversationsInviteRequest conversationsInviteRequest =
              ConversationsInviteRequest.builder()
                  .channel(channelId)
                  .users(
                      Collections.singletonList(
                          payload.getPayload().getUser().getId())) // User ID from the callback
                  .build();

          ConversationsInviteResponse conversationsInviteResponse =
              slackClient.methods(botToken).conversationsInvite(conversationsInviteRequest);

          if (!conversationsInviteResponse.isOk()) {
            throw new IOException(
                "Failed to invite user to Slack channel: "
                    + conversationsInviteResponse.getError());
          }

          // set incident channel
          IncidentChannel incidentChannel = new IncidentChannel();
          io.respondnow.model.incident.Slack slack = new io.respondnow.model.incident.Slack();
          slack.setChannelId(incidentChannelID);
          slack.setTeamDomain(payload.getPayload().getTeam().getDomain());
          slack.setTeamId(payload.getPayload().getTeam().getId());
          //          slack.setTeamName(payload.getPayload().getTeam().get);
          incidentChannel.setType(IncidentChannelType.Slack);
          incidentChannel.setSlack(slack);

          // set channels
          List<Channel> channels = new ArrayList<>();
          Channel channel1 =
              new Channel(
                  channelId,
                  payload.getPayload().getTeam().getId(),
                  createChannelResponse.getChannel().getName(),
                  ChannelSource.Slack,
                  String.format(
                      "https://%s.slack.com/archives/%s",
                      payload.getPayload().getTeam().getId(), channelId),
                  Operational);
          channels.add(channel1);

          // Create incident record in the database
          CreateRequest createRequest = new CreateRequest();
          createRequest.setIdentifier(incidentId);
          createRequest.setName(name);
          createRequest.setType(Type.valueOf(incidentType));
          createRequest.setStatus(Status.Started);
          createRequest.setRoles(roles);
          createRequest.setSeverity(Severity.valueOf(severity));
          createRequest.setSummary(summary);
          createRequest.setIncidentChannel(incidentChannel);
          createRequest.setChannels(channels);
          Incident incident = incidentService.createIncident(createRequest, finalUserDetails);

          // Post messages in Slack
          sendCreateIncidentResponseMsg(
              payload.getPayload().getTeam().getDomain(),
              responseChannel,
              channelId,
              sanitizedChannelName,
              incident);

        } catch (SlackApiException | IOException e) {
          throw new RuntimeException(e);
        }

      } catch (SlackApiException | IOException e) {
        throw new RuntimeException(e);
      }
    } else {
      logger.error("Failed to create incident due to missing required fields.");
    }
  }

  // Method to fetch user details (you'll need to implement this based on your data source)
  private UserDetails fetchSlackUserDetails(String userId, ChannelSource source)
      throws SlackApiException, IOException {
    UserDetails userDetails = new UserDetails();
    if (source == ChannelSource.Slack) {
      UsersInfoResponse slackUser = getSlackUserDetails(userId);

      userDetails.setUserId(slackUser.getUser().getId());
      userDetails.setUserName(slackUser.getUser().getName());
      userDetails.setName(slackUser.getUser().getProfile().getRealName());
      userDetails.setEmail(slackUser.getUser().getProfile().getEmail());
    }
    if (source != null) {
      userDetails.setSource(source);
    }
    return userDetails;
  }

  // Function to sanitize channel name
  private String sanitizeChannelName(String channelName) {
    // Replace any non-alphanumeric character or space with a hyphen
    String sanitized = channelName.replaceAll("[^a-zA-Z0-9-]", "-").toLowerCase();

    // Ensure the name is not too long (max length: 80 characters)
    if (sanitized.length() > 80) {
      sanitized = sanitized.substring(0, 80);
    }

    // Ensure the name doesn't start or end with a hyphen
    sanitized = sanitized.replaceAll("^-|-$", "");

    return sanitized;
  }

  private SectionBlock createIncidentNameAndSeveritySection(Incident newIncident) {
    return SectionBlock.builder()
        .fields(
            List.of(
                buildMarkdownTextObject(newIncident.getName(), "Name", ":writing_hand:"),
                buildMarkdownTextObject(
                    newIncident.getSeverity().toString(), "Severity", ":vertical_traffic_light:")))
        .build();
  }

  private SectionBlock createIncidentStatusAndCommanderSection(
      Incident newIncident, String incidentCommander) {
    List<TextObject> blocks = new ArrayList<TextObject>();
    blocks.add(
        buildMarkdownTextObject(newIncident.getStatus().toString(), "Current Status", ":eyes:"));
    if (!incidentCommander.isEmpty()) {
      blocks.add(
          buildMarkdownTextObject(
              "<@" + incidentCommander + ">", "Incident Commander", ":firefighter:"));
    }
    return SectionBlock.builder().fields(blocks).build();
  }

  private MarkdownTextObject buildMarkdownTextObject(String value, String label, String emoji) {
    return MarkdownTextObject.builder()
        .text(String.format("%s *%s*\n _%s_", emoji, label, value))
        .build();
  }

  public void sendCreateIncidentResponseMsg(
      String teamDomain,
      String channelID,
      String joinChannelID,
      String joinChannelName,
      Incident newIncident) {
    // Prepare blocks
    List<LayoutBlock> blocks = new ArrayList<>();
    try {
      // Find the commander
      String incidentCommander =
          newIncident.getRoles().stream()
              .filter(role -> role.getRoleType() == RoleType.Incident_Commander)
              .map(role -> role.getUserDetails().getUserId())
              .findFirst()
              .orElse("");

      // Find the communications lead
      String communicationsLead =
          newIncident.getRoles().stream()
              .filter(role -> role.getRoleType() == RoleType.Communications_Lead)
              .map(role -> role.getUserDetails().getUserId())
              .findFirst()
              .orElse("");

      blocks.add(createHeaderBlock(":fire: :mega: New Incident"));

      // Name, Severity Section
      blocks.add(createIncidentNameAndSeveritySection(newIncident));

      // Summary Section
      blocks.add(createIncidentDetailsSection(newIncident.getSummary(), "Summary", ":open_book:"));

      // Divider
      blocks.add(new DividerBlock());

      // Status and Commander
      blocks.add(createIncidentStatusAndCommanderSection(newIncident, incidentCommander));

      if (!communicationsLead.isEmpty()) {
        blocks.add(
            createIncidentDetailsSection(
                "<@" + communicationsLead + ">", "Communications Lead", ":phone:"));
      }

      boolean shouldBePinned = true;

      // Action Buttons
      blocks.add(createActionButtons(newIncident.getIdentifier()));

      // Created At Information
      blocks.add(createCreatedAtBlock(newIncident));

      // Send Message to newly created incident channel
      ChatPostMessageResponse response =
          slackApp.getClient().chatPostMessage(r -> r.channel(joinChannelID).blocks(blocks));

      if (!response.isOk()) {
        logger.error("Error sending message to Slack: {}", response.getError());
      }

      if (shouldBePinned) {
        // Pin the message if required
        slackApp
            .getClient()
            .pinsAdd(r -> r.channel(joinChannelID).timestamp(response.getMessage().getTs()));
      }

      if (joinChannelID != null
          && !joinChannelID.isEmpty()
          && joinChannelName != null
          && !joinChannelName.isEmpty()) {
        shouldBePinned = false;
        blocks.add(createJoinChannelButton(teamDomain, joinChannelID, joinChannelName));
      }
      // Post another response in the incident channel
      ChatPostMessageResponse incidentChannelResponse =
          slackApp.client().chatPostMessage(r -> r.channel(channelID).blocks(blocks));
      if (shouldBePinned) {
        // Pin the message if required
        slackApp
            .getClient()
            .pinsAdd(
                r -> r.channel(channelID).timestamp(incidentChannelResponse.getMessage().getTs()));
      }

      // Notify users based on role
      for (Role role : newIncident.getRoles()) {
        String userId = role.getUserDetails().getUserId();

        // Send notification to the user
        sendRoleNotificationToUser(userId, role, teamDomain, joinChannelID);
      }
    } catch (IOException | SlackApiException e) {
      logger.error("Failed to send incident response message: {}", e.getMessage(), e);
    }
  }

  private void sendRoleNotificationToUser(
      String userId, Role role, String teamDomain, String channelID) {
    try {
      // Create the notification blocks based on the role
      List<LayoutBlock> blocks = getUserRoleNotificationBlocks(role, teamDomain, channelID);

      // Send the notification message
      ChatPostMessageResponse notificationResponse =
          slackApp.getClient().chatPostMessage(r -> r.channel(userId).blocks(blocks));

      if (!notificationResponse.isOk()) {
        logger.error(
            "Error sending notification to user {}: {}", userId, notificationResponse.getError());
      }
    } catch (SlackApiException | IOException e) {
      logger.error("Failed to send role notification: {}", e.getMessage(), e);
    }
  }

  private List<LayoutBlock> getUserRoleNotificationBlocks(
      Role role, String teamDomain, String channelID) {
    List<LayoutBlock> slackBlocks = new ArrayList<>();

    // Header block with the role message
    slackBlocks.add(
        SectionBlock.builder()
            .text(
                MarkdownTextObject.builder()
                    .text(
                        String.format(
                            ":wave: You have been elected as the *%s* for an incident.",
                            role.getRoleType()))
                    .build())
            .blockId("user_role_notification_header")
            .build());

    // Role-specific description blocks
    if (role.getRoleType() == RoleType.Incident_Commander) {
      slackBlocks.add(
          SectionBlock.builder()
              .text(
                  MarkdownTextObject.builder().text(getIncidentCommanderRoleDescription()).build())
              .blockId("incident_commander_role_description")
              .build());
    } else if (role.getRoleType() == RoleType.Communications_Lead) {
      slackBlocks.add(
          SectionBlock.builder()
              .text(
                  MarkdownTextObject.builder().text(getCommunicationsLeadRoleDescription()).build())
              .blockId("communications_lead_role_description")
              .build());
    }

    // Divider and channel info
    slackBlocks.add(DividerBlock.builder().build());
    slackBlocks.add(
        SectionBlock.builder()
            .text(
                MarkdownTextObject.builder()
                    .text(String.format("Please join the channel here: <#%s>", channelID))
                    .build())
            .blockId("user_role_notification_channel")
            .build());

    return slackBlocks;
  }

  private String getIncidentCommanderRoleDescription() {
    return "The Incident Commander is the decision maker during a major incident, delegating tasks and listening to input from subject matter experts in order to bring the incident to resolution. They become the highest ranking individual on any major incident call, regardless of their day-to-day rank. Their decisions made as commander are final.\n\nYour job as an Incident Commander is to listen to the call and to watch the incident Slack room in order to provide clear coordination, recruiting others to gather context and details. You should not be performing any actions or remediations, checking graphs, or investigating logs. Those tasks should be delegated.\n\nAn IC should also be considering next steps and backup plans at every opportunity, in an effort to avoid getting stuck without any clear options to proceed and to keep things moving towards resolution.\n\nMore information: https://response.pagerduty.com/training/incident_commander/";
  }

  private String getCommunicationsLeadRoleDescription() {
    return "The purpose of the Communications Liaison is to be the primary individual in charge of notifying our customers of the current conditions, and informing the Incident Commander of any relevant feedback from customers as the incident progresses.\n\nIt's important for the rest of the command staff to be able to focus on the problem at hand, rather than worrying about crafting messages to customers.\n\nYour job as Communications Liaison is to listen to the call, watch the incident Slack room, and track incoming customer support requests, keeping track of what's going on and how far the incident is progressing (still investigating vs close to resolution).\n\nThe Incident Commander will instruct you to notify customers of the incident and keep them updated at various points throughout the call. You will be required to craft the message, gain approval from the IC, and then disseminate that message to customers.\n\nMore information: https://response.pagerduty.com/training/customer_liaison/";
  }

  private SectionBlock createIncidentDetailsSection(String value, String label, String emoji) {
    return SectionBlock.builder()
        .text(
            MarkdownTextObject.builder()
                .text(String.format("%s *%s*\n _%s_", emoji, label, value))
                .build())
        .build();
  }

  private HeaderBlock createHeaderBlock(String text) {
    return HeaderBlock.builder()
        .text(PlainTextObject.builder().text(text).build())
        .blockId("create_incident_channel_resp_header")
        .build();
  }

  private ActionsBlock createJoinChannelButton(
      String teamDomain, String joinChannelID, String joinChannelName) {
    return ActionsBlock.builder()
        .elements(
            List.of(
                ButtonElement.builder()
                    .text(
                        PlainTextObject.builder()
                            .text(String.format(":slack: Join %s", joinChannelName))
                            .build())
                    .url(
                        String.format(
                            "https://%s.slack.com/archives/%s", teamDomain, joinChannelID))
                    .style("primary")
                    .build()))
        .blockId("create_incident_channel_join_channel")
        .build();
  }

  private ActionsBlock createActionButtons(String incidentId) {
    return ActionsBlock.builder()
        .elements(
            List.of(
                createActionButton("Update Summary", "update_incident_summary_button", incidentId),
                createActionButton("Add a comment", "update_incident_comment_button", incidentId),
                createActionButton(
                    "Assign Roles", "update_incident_assign_roles_button", incidentId),
                createActionButton("Update Status", "update_incident_status_button", incidentId),
                createActionButton(
                    "Update Severity", "update_incident_severity_button", incidentId)))
        .blockId("incident_action_buttons")
        .build();
  }

  private ButtonElement createActionButton(String text, String actionId, String value) {
    return ButtonElement.builder()
        .text(PlainTextObject.builder().text(text).build())
        .actionId(actionId)
        .value(value)
        .build();
  }

  private ContextBlock createCreatedAtBlock(Incident incident) {
    return ContextBlock.builder()
        .elements(
            List.of(
                PlainTextObject.builder()
                    .text(String.format(":clock1: Started At: %s", incident.getCreatedAt()))
                    .build(),
                PlainTextObject.builder()
                    .text(
                        String.format(
                            ":man: Started By: %s", incident.getCreatedBy().getUserName()))
                    .build()))
        .blockId("create_incident_channel_resp_createdAt")
        .build();
  }

  public void handleIncidentSummaryViewSubmission(ViewSubmissionRequest payload)
      throws SlackApiException, IOException {
    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    String updatedSummary =
        payload
            .getPayload()
            .getView()
            .getState()
            .getValues()
            .get("create_incident_modal_summary")
            .get("create_incident_modal_set_summary")
            .getValue();
    logger.info("Incident Identifier: {}, Updated Summary: {}", incidentIdentifier, updatedSummary);

    // Create and call the service to update the incident summary
    updateIncidentSummary(
        incidentIdentifier,
        updatedSummary,
        fetchSlackUserDetails(payload.getPayload().getUser().getId(), ChannelSource.Slack));
  }

  private void updateIncidentRoles(
      String incidentIdentifier, List<Role> roles, UserDetails userDetails) {
    try {
      Incident updatedIncident =
          incidentService.updateIncidentRoles(incidentIdentifier, roles, userDetails);
      sendUpdateRoleResponseMsg(
          updatedIncident.getChannels().get(0).getId(), updatedIncident, roles);
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

  private void updateIncidentSummary(
      String incidentIdentifier, String updatedSummary, UserDetails user) {
    // Simulate a service call to update the incident summary
    try {
      Incident updatedIncident =
          incidentService.updateSummary(incidentIdentifier, updatedSummary, user);

      // Send confirmation message to Slack
      sendUpdateSummaryResponseMsg(
          updatedIncident.getChannels().get(0).getId(), updatedIncident, updatedSummary);
    } catch (Exception e) {
      logger.error("Failed to update incident summary: {}", e.getMessage(), e);
    }
  }

  private void updateIncidentSeverity(
      String incidentIdentifier, String severity, UserDetails user) {
    Severity newSeverity = Severity.valueOf(severity);
    try {
      Incident updatedIncident =
          incidentService.updateIncidentSeverity(incidentIdentifier, newSeverity, user);

      sendUpdateSeverityResponseMsg(
          updatedIncident.getChannels().get(0).getId(), updatedIncident, severity);
    } catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

  private void updateIncidentStatus(String incidentIdentifier, String status, UserDetails user) {
    Status updatedStatus = Status.valueOf(status);
    // Simulate a service call to add a new incident comment
    try {
      Incident updatedIncident =
          incidentService.updateStatus(incidentIdentifier, updatedStatus, user);

      // Send confirmation message to Slack
      sendUpdateStatusResponseMsg(
          updatedIncident.getChannels().get(0).getId(), updatedIncident, status);
    } catch (Exception e) {
      logger.error("Failed to add a new incident comment: {}", e.getMessage(), e);
    }
  }

  private void updateIncidentComment(String incidentIdentifier, String comment, UserDetails user) {
    // Simulate a service call to add a new incident comment
    try {
      Incident updatedIncident = incidentService.addComment(incidentIdentifier, comment, user);

      // Send confirmation message to Slack
      sendAddCommentResponseMsg(
          updatedIncident.getChannels().get(0).getId(), updatedIncident, comment);
    } catch (Exception e) {
      logger.error("Failed to add a new incident comment: {}", e.getMessage(), e);
    }
  }

  private UsersInfoResponse getSlackUserDetails(String userId)
      throws SlackApiException, IOException {
    // Fetch user info from Slack using userId
    return slackApp.client().usersInfo(r -> r.user(userId));
  }

  private void sendUpdateRoleResponseMsg(
      String channelID, Incident updatedIncident, List<Role> newRoles) {
    try {
      // Fetch user info from Slack using userId
      UsersInfoResponse slackUserInfo =
          getSlackUserDetails(updatedIncident.getUpdatedBy().getUserId());
      if (slackUserInfo == null || slackUserInfo.getUser() == null) {
        logger.error(
            "Role update: failed to fetch Slack user info for userId: {}",
            updatedIncident.getUpdatedBy().getUserId());
        throw new IllegalStateException("Unable to fetch Slack user info");
      }

      // Get the Slack handle (username)
      List<String> roleUpdates = new ArrayList<>();
      for (Role newRole : newRoles) {
        if (newRole.getUserDetails().getUserId() != null) {
          roleUpdates.add(
              String.format(
                  "*%s*: <@%s>",
                  newRole.getRoleType().getDisplayValue(), newRole.getUserDetails().getUserId()));
        }
      }

      String messageText =
          String.format(
              ":firefighter: *Roles Updated*\nThe following roles have been updated:\n%s",
              Strings.join(roleUpdates, '\n'));

      // Send the added comment response message back to the Slack channel
      ChatPostMessageResponse response =
          slackApp.client().chatPostMessage(r -> r.channel(channelID).text(messageText));

      if (!response.isOk()) {
        String errorMessage = "Failed to send update role error message: " + response.getError();
        logger.error(errorMessage);
        throw new RuntimeException(errorMessage);
      }

      logger.info("Update roles confirmation successfully posted to channel: {}", channelID);
    } catch (IOException e) {
      logger.error("IOException occurred while posting message to Slack: {}", e.getMessage(), e);
      // Optionally, add a retry mechanism or send a failure notification to users
    } catch (SlackApiException e) {
      logger.error("Slack API error occurred while posting message: {}", e.getMessage(), e);
      // Handle specific Slack API exceptions if needed (e.g., retry on rate limit errors)
    } catch (IllegalStateException e) {
      logger.error("Error in fetching Slack user info: {}", e.getMessage(), e);
    } catch (Exception e) {
      logger.error("Unexpected error occurred: {}", e.getMessage(), e);
      // Optionally, send a failure notification or alert to a monitoring system
    }
  }

  private void sendAddCommentResponseMsg(
      String channelID, Incident updatedIncident, String newComment) {
    try {
      // Fetch user info from Slack using userId
      UsersInfoResponse slackUserInfo =
          getSlackUserDetails(updatedIncident.getUpdatedBy().getUserId());
      if (slackUserInfo == null || slackUserInfo.getUser() == null) {
        logger.error(
            "Add comment: failed to fetch Slack user info for userId: {}",
            updatedIncident.getUpdatedBy().getUserId());
        throw new IllegalStateException("Unable to fetch Slack user info");
      }

      // Get the Slack handle (username)
      String slackHandle = slackUserInfo.getUser().getName();

      // Prepare the message text
      String messageText =
          String.format(
              ":speech_balloon: *Comment Added*\n <@%s> added a new comment:\n> _%s_",
              slackHandle, newComment);

      // Send the added comment response message back to the Slack channel
      ChatPostMessageResponse response =
          slackApp.client().chatPostMessage(r -> r.channel(channelID).text(messageText));

      if (!response.isOk()) {
        String errorMessage = "Failed to send add comment error message: " + response.getError();
        logger.error(errorMessage);
        throw new RuntimeException(errorMessage);
      }

      logger.info("Comment addition confirmation successfully posted to channel: {}", channelID);
    } catch (IOException e) {
      logger.error("IOException occurred while posting message to Slack: {}", e.getMessage(), e);
      // Optionally, add a retry mechanism or send a failure notification to users
    } catch (SlackApiException e) {
      logger.error("Slack API error occurred while posting message: {}", e.getMessage(), e);
      // Handle specific Slack API exceptions if needed (e.g., retry on rate limit errors)
    } catch (IllegalStateException e) {
      logger.error("Error in fetching Slack user info: {}", e.getMessage(), e);
    } catch (Exception e) {
      logger.error("Unexpected error occurred: {}", e.getMessage(), e);
      // Optionally, send a failure notification or alert to a monitoring system
    }
  }

  private void sendUpdateStatusResponseMsg(
      String channelID, Incident updatedIncident, String newStatus) {
    try {
      // Fetch user info from Slack using userId
      UsersInfoResponse slackUserInfo =
          getSlackUserDetails(updatedIncident.getUpdatedBy().getUserId());
      if (slackUserInfo == null || slackUserInfo.getUser() == null) {
        logger.error(
            "Update status: failed to fetch Slack user info for userId: {}",
            updatedIncident.getUpdatedBy().getUserId());
        throw new IllegalStateException("Unable to fetch Slack user info");
      }

      // Get the Slack handle (username)
      String slackHandle = slackUserInfo.getUser().getName();

      // Prepare the message text
      String messageText =
          String.format(
              ":eyes: *Status Updated*\n <@%s> updated the status to: _%s_",
              slackHandle, newStatus);

      // Send the added comment response message back to the Slack channel
      ChatPostMessageResponse response =
          slackApp.client().chatPostMessage(r -> r.channel(channelID).text(messageText));

      if (!response.isOk()) {
        String errorMessage = "Failed to send add comment error message: " + response.getError();
        logger.error(errorMessage);
        throw new RuntimeException(errorMessage);
      }

      logger.info("Comment addition confirmation successfully posted to channel: {}", channelID);
    } catch (IOException e) {
      logger.error("IOException occurred while posting message to Slack: {}", e.getMessage(), e);
      // Optionally, add a retry mechanism or send a failure notification to users
    } catch (SlackApiException e) {
      logger.error("Slack API error occurred while posting message: {}", e.getMessage(), e);
      // Handle specific Slack API exceptions if needed (e.g., retry on rate limit errors)
    } catch (IllegalStateException e) {
      logger.error("Error in fetching Slack user info: {}", e.getMessage(), e);
    } catch (Exception e) {
      logger.error("Unexpected error occurred: {}", e.getMessage(), e);
      // Optionally, send a failure notification or alert to a monitoring system
    }
  }

  private void sendUpdateSeverityResponseMsg(
      String channelID, Incident updatedIncident, String newSeverity) {
    try {
      // Fetch user info from Slack using userId
      UsersInfoResponse slackUserInfo =
          getSlackUserDetails(updatedIncident.getUpdatedBy().getUserId());
      if (slackUserInfo == null || slackUserInfo.getUser() == null) {
        logger.error(
            "Failed to fetch Slack user info for userId: {}",
            updatedIncident.getUpdatedBy().getUserId());
        throw new IllegalStateException("Unable to fetch Slack user info");
      }

      // Get the Slack handle (username)
      String slackHandle = slackUserInfo.getUser().getName();

      // Prepare the message text
      String messageText =
          String.format(
              ":vertical_traffic_light: *Severity Updated*\n <@%s> updated the severity to: _%s_",
              slackHandle, newSeverity);

      // Send the update summary response message back to the Slack channel
      ChatPostMessageResponse response =
          slackApp.client().chatPostMessage(r -> r.channel(channelID).text(messageText));

      if (!response.isOk()) {
        String errorMessage = "Failed to send message: " + response.getError();
        logger.error(errorMessage);
        throw new RuntimeException(errorMessage);
      }

      logger.info("Summary update confirmation successfully posted to channel: {}", channelID);
    } catch (IOException e) {
      logger.error("IOException occurred while posting message to Slack: {}", e.getMessage(), e);
      // Optionally, add a retry mechanism or send a failure notification to users
    } catch (SlackApiException e) {
      logger.error("Slack API error occurred while posting message: {}", e.getMessage(), e);
      // Handle specific Slack API exceptions if needed (e.g., retry on rate limit errors)
    } catch (IllegalStateException e) {
      logger.error("Error in fetching Slack user info: {}", e.getMessage(), e);
    } catch (Exception e) {
      logger.error("Unexpected error occurred: {}", e.getMessage(), e);
      // Optionally, send a failure notification or alert to a monitoring system
    }
  }

  private void sendUpdateSummaryResponseMsg(
      String channelID, Incident updatedIncident, String newSummary) {
    try {
      // Fetch user info from Slack using userId
      UsersInfoResponse slackUserInfo =
          getSlackUserDetails(updatedIncident.getUpdatedBy().getUserId());
      if (slackUserInfo == null || slackUserInfo.getUser() == null) {
        logger.error(
            "Failed to fetch Slack user info for userId: {}",
            updatedIncident.getUpdatedBy().getUserId());
        throw new IllegalStateException("Unable to fetch Slack user info");
      }

      // Get the Slack handle (username)
      String slackHandle = slackUserInfo.getUser().getName();

      // Prepare the message text
      String messageText =
          String.format(
              ":memo: *Summary Updated*\n <@%s> updated the summary:\n> _%s_",
              slackHandle, newSummary);

      // Send the update summary response message back to the Slack channel
      ChatPostMessageResponse response =
          slackApp.client().chatPostMessage(r -> r.channel(channelID).text(messageText));

      if (!response.isOk()) {
        String errorMessage = "Failed to send message: " + response.getError();
        logger.error(errorMessage);
        throw new RuntimeException(errorMessage);
      }

      logger.info("Summary update confirmation successfully posted to channel: {}", channelID);
    } catch (IOException e) {
      logger.error("IOException occurred while posting message to Slack: {}", e.getMessage(), e);
      // Optionally, add a retry mechanism or send a failure notification to users
    } catch (SlackApiException e) {
      logger.error("Slack API error occurred while posting message: {}", e.getMessage(), e);
      // Handle specific Slack API exceptions if needed (e.g., retry on rate limit errors)
    } catch (IllegalStateException e) {
      logger.error("Error in fetching Slack user info: {}", e.getMessage(), e);
    } catch (Exception e) {
      logger.error("Unexpected error occurred: {}", e.getMessage(), e);
      // Optionally, send a failure notification or alert to a monitoring system
    }
  }

  public void handleIncidentCommentViewSubmission(ViewSubmissionRequest payload)
      throws SlackApiException, IOException {
    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    String updatedComment =
        payload
            .getPayload()
            .getView()
            .getState()
            .getValues()
            .get("update_incident_modal_comment")
            .get("update_incident_modal_set_comment")
            .getValue();
    logger.info("Incident Identifier: {}, Updated Comment: {}", incidentIdentifier, updatedComment);

    updateIncidentComment(
        incidentIdentifier,
        updatedComment,
        fetchSlackUserDetails(payload.getPayload().getUser().getId(), ChannelSource.Slack));
  }

  public void handleIncidentRolesViewSubmission(ViewSubmissionRequest payload)
      throws SlackApiException, IOException {

    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    ViewState state = payload.getPayload().getView().getState();
    Map<String, Map<String, ViewState.Value>> stateValues = state.getValues();

    // Extract roles based on RoleType enum
    List<String> roles =
        Arrays.stream(RoleType.values()).map(RoleType::getValue).collect(Collectors.toList());

    Map<String, UserDetails> roleUserDetails = new HashMap<>();

    for (Map<String, ViewState.Value> viewState : stateValues.values()) {
      for (String role : roles) {
        String viewStateKey = "create_incident_modal_set_" + role;
        ViewState.Value viewStateValue = viewState.get(viewStateKey);
        if (viewStateValue != null) {
          String selectedUser = viewStateValue.getSelectedUser();
          if (selectedUser != null && !selectedUser.isEmpty()) {
            roleUserDetails.put(role, fetchSlackUserDetails(selectedUser, ChannelSource.Slack));
          } else {
            roleUserDetails.put(role, new UserDetails()); // Or handle as per your business logic
          }
        }
      }
    }

    // Convert roleUserDetails to a list of Role objects
    List<Role> roleList =
        roleUserDetails.entrySet().stream()
            .map(
                entry -> {
                  Role role = new Role();
                  role.setRoleType(RoleType.fromValue(entry.getKey()));
                  role.setUserDetails(entry.getValue());
                  return role;
                })
            .collect(Collectors.toList());

    try {
      Incident updatedIncident =
          incidentService.updateIncidentRoles(
              incidentIdentifier,
              roleList,
              fetchSlackUserDetails(payload.getPayload().getUser().getId(), ChannelSource.Slack));

      sendUpdateRoleResponseMsg(
          updatedIncident.getChannels().get(0).getId(), updatedIncident, roleList);
      //      sendSlackMessage(
      //          payload.getPayload().getUser().getId(), "Incident roles have been successfully
      // updated.");
    } catch (RoleUpdateException e) {
      logger.error("Role update failed for incident ID: {}", incidentIdentifier, e);
      sendSlackMessage(
          payload.getPayload().getUser().getId(),
          "Failed to update incident roles: " + e.getMessage());
    } catch (Exception e) {
      logger.error(
          "Unexpected error during role update for incident ID: {}", incidentIdentifier, e);
      sendSlackMessage(
          payload.getPayload().getUser().getId(),
          "An unexpected error occurred while updating roles. Please try again.");
    }
  }

  /**
   * Sends a message to a user via Slack.
   *
   * @param userId The Slack user ID.
   * @param message The message to send.
   */
  private void sendSlackMessage(String userId, String message) {
    try {
      slackApp.client().chatPostMessage(r -> r.channel(userId).text(message));
      logger.info("Sent message to user {}: {}", userId, message);
    } catch (Exception e) {
      logger.error("Failed to send message to user {}: {}", userId, e.getMessage(), e);
    }
  }

  public void handleIncidentStatusViewSubmission(ViewSubmissionRequest payload)
      throws SlackApiException, IOException {
    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    String status =
        payload
            .getPayload()
            .getView()
            .getState()
            .getValues()
            .get("incident_status")
            .get("create_incident_modal_set_incident_status")
            .getSelectedOption()
            .getValue();
    updateIncidentStatus(
        incidentIdentifier,
        status,
        fetchSlackUserDetails(payload.getPayload().getUser().getId(), ChannelSource.Slack));
    logger.info("Incident Identifier: {}, Updated Status: {}", incidentIdentifier, status);
  }

  public void handleIncidentSeverityViewSubmission(ViewSubmissionRequest payload)
      throws SlackApiException, IOException {
    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    String severity =
        payload
            .getPayload()
            .getView()
            .getState()
            .getValues()
            .get("incident_severity")
            .get("create_incident_modal_set_incident_severity")
            .getSelectedOption()
            .getValue();

    updateIncidentSeverity(
        incidentIdentifier,
        severity,
        fetchSlackUserDetails(payload.getPayload().getUser().getId(), ChannelSource.Slack));
    logger.info("Incident Identifier: {}, Updated Severity: {}", incidentIdentifier, severity);
  }

  /** Registers the "Create Incident" button handler. */
  private void registerCreateIncidentButton() {
    slackApp.blockAction(
        "create_incident_modal", // The action ID for the button
        (payload, ctx) -> {
          // Perform processing asynchronously
          CompletableFuture.runAsync(
              () -> {
                try {
                  logger.info("Handling 'Create Incident' button click: {}", payload);

                  // Build the modal view
                  List<LayoutBlock> blocks = new ArrayList<>();
                  blocks.add(
                      SectionBlock.builder()
                          .text(
                              MarkdownTextObject.builder()
                                  .text(
                                      "This will start a new incident channel, and you will "
                                          + "be invited to it. From there, please use our incident "
                                          + "management process to run the incident or coordinate "
                                          + "with others to do so.")
                                  .build())
                          .build());
                  blocks.add(getNameBlock());
                  blocks.add(getTypeBlock());
                  blocks.add(
                      getSummaryBlock(
                          "create_incident_modal_summary", "create_incident_modal_set_summary"));
                  blocks.add(getSeverityBlock());
                  blocks.add(getRoleBlock());
                  blocks.add(getChannelSelectBlock());

                  View modalView =
                      Views.view(
                          v ->
                              v.type("modal")
                                  .callbackId("create_incident_modal")
                                  .title(
                                      Views.viewTitle(
                                          title ->
                                              title
                                                  .type("plain_text")
                                                  .text("Start a new incident")))
                                  .blocks(blocks)
                                  .submit(
                                      Views.viewSubmit(
                                          submit -> submit.type("plain_text").text("Start"))));

                  // Open the modal view
                  ViewsOpenResponse response =
                      ctx.client()
                          .viewsOpen(
                              r ->
                                  r.triggerId(payload.getPayload().getTriggerId()).view(modalView));

                  if (response.isOk()) {
                    logger.info("Create Incident modal opened successfully.");
                  } else {
                    logger.error("Failed to open Create Incident modal: {}", response.getError());
                    // Optionally, notify the user about the failure
                    slackApp
                        .client()
                        .chatPostMessage(
                            r ->
                                r.channel(payload.getPayload().getUser().getId())
                                    .text(
                                        "Failed to open the incident creation modal. Please try again."));
                  }
                } catch (Exception e) {
                  logger.error("Exception while handling 'Create Incident' button click", e);
                  // Optionally, notify the user about the exception
                  try {
                    slackApp
                        .client()
                        .chatPostMessage(
                            r ->
                                r.channel(payload.getPayload().getUser().getId())
                                    .text(
                                        "An error occurred while creating the incident. Please try again."));
                  } catch (IOException | SlackApiException ex) {
                    throw new RuntimeException(ex);
                  }
                }
              },
              executorService);

          // Acknowledge immediately
          return ctx.ack();
        });
  }

  public void createIncident(GlobalShortcutRequest req, GlobalShortcutContext ctx) {
    try {
      List<LayoutBlock> blocks = new ArrayList<>();
      blocks.add(
          SectionBlock.builder()
              .text(
                  MarkdownTextObject.builder()
                      .text(
                          "This will start a new incident channel, and you will "
                              + "be invited to it. From there, please use our incident "
                              + "management process to run the incident or coordinate "
                              + "with others to do so.")
                      .build())
              .build());
      blocks.add(getNameBlock());
      blocks.add(getTypeBlock());
      blocks.add(
          getSummaryBlock("create_incident_modal_summary", "create_incident_modal_set_summary"));
      blocks.add(getSeverityBlock());
      blocks.add(getRoleBlock());
      blocks.add(getChannelSelectBlock());

      View modalView =
          Views.view(
              v ->
                  v.type("modal")
                      .callbackId("create_incident_modal")
                      .title(
                          Views.viewTitle(
                              title -> title.type("plain_text").text("Start a new incident")))
                      .blocks(blocks)
                      .submit(Views.viewSubmit(submit -> submit.type("plain_text").text("Start"))));

      ViewsOpenResponse response =
          ctx.client().viewsOpen(r -> r.triggerId(req.getPayload().getTriggerId()).view(modalView));

      logger.info("CreateIncident view opened successfully: {}", response.isOk());
      ctx.ack();
    } catch (Exception e) {
      logger.error("Failed to open create incident view modal: {}", e.getMessage(), e);
      ctx.ackWithJson(errorResponse());
    }
  }

  private InputBlock getNameBlock() {
    return Blocks.input(
        i ->
            i.blockId("create_incident_modal_name")
                .label(new PlainTextObject(":writing_hand: Incident Name", true))
                .element(
                    BlockElements.plainTextInput(
                        pt ->
                            pt.actionId("create_incident_modal_set_name")
                                .placeholder(new PlainTextObject("IAM service is down", false))
                                .maxLength(76))));
  }

  private InputBlock getTypeBlock() {
    List<OptionObject> options = getIncidentTypes(); // Fetch supported incident types
    return Blocks.input(
        i ->
            i.blockId("incident_type")
                .label(new PlainTextObject(":fire: Incident Type", true))
                .element(
                    BlockElements.staticSelect(
                        s ->
                            s.actionId("create_incident_modal_set_incident_type")
                                .placeholder(new PlainTextObject("Select incident type...", false))
                                .options(options))));
  }

  private InputBlock getSummaryBlock(String blockId, String actionId) {
    return Blocks.input(
        i ->
            i.blockId(blockId)
                .label(new PlainTextObject(":memo: Summary", true))
                .element(
                    BlockElements.plainTextInput(
                        pt ->
                            pt.actionId(actionId)
                                .multiline(true)
                                .placeholder(
                                    new PlainTextObject(
                                        "A brief description of the problem.", false)))));
  }

  private InputBlock getCommentBlock(String blockId, String actionId) {
    return Blocks.input(
        i ->
            i.blockId(blockId)
                .label(new PlainTextObject(":speech_balloon: Comment", true))
                .element(
                    BlockElements.plainTextInput(
                        pt ->
                            pt.actionId(actionId)
                                .multiline(true)
                                .placeholder(new PlainTextObject("Add a comment", false)))));
  }

  private InputBlock getSeverityBlock() {
    List<OptionObject> options = getIncidentSeverities(); // Fetch severities
    return Blocks.input(
        i ->
            i.blockId("incident_severity")
                .label(new PlainTextObject(":vertical_traffic_light: Severity", true))
                .element(
                    BlockElements.staticSelect(
                        s ->
                            s.actionId("create_incident_modal_set_incident_severity")
                                .placeholder(
                                    new PlainTextObject(
                                        "Select severity of the incident...", false))
                                .options(options))));
  }

  private InputBlock getRoleBlock() {
    List<OptionObject> options = getIncidentRoles(); // Fetch roles
    return Blocks.input(
        i ->
            i.blockId("incident_role")
                .label(new PlainTextObject(":firefighter: Assign role to yourself", true))
                .element(
                    BlockElements.multiStaticSelect(
                        ms ->
                            ms.actionId("create_incident_modal_set_incident_role")
                                .placeholder(
                                    new PlainTextObject("Assign role to yourself...", false))
                                .options(options))));
  }

  private List<LayoutBlock> getUpdateRoleBlock() {
    List<LayoutBlock> blocks = new ArrayList<>();
    for (RoleType role : RoleType.values()) {
      TextObject roleText =
          PlainTextObject.builder().text(role.getDisplayValue()).emoji(true).build();

      UsersSelectElement userSelect =
          BlockElements.usersSelect(
              r ->
                  r.placeholder(PlainTextObject.builder().text("Select a user").emoji(true).build())
                      .actionId("create_incident_modal_set_" + role));

      SectionBlock section = SectionBlock.builder().text(roleText).accessory(userSelect).build();

      blocks.add(section);
    }
    return blocks;
  }

  public InputBlock updateStatus() {
    List<Status> supportedIncidentStatuses = Arrays.asList(Status.values());

    if (supportedIncidentStatuses.isEmpty()) {
      throw new IllegalStateException("No incident statuses available.");
    }
    List<OptionObject> incidentStatusOptions = new ArrayList<>();

    for (Status incidentStatus : supportedIncidentStatuses) {
      OptionObject option =
          OptionObject.builder()
              .text(PlainTextObject.builder().text(incidentStatus.getValue()).emoji(true).build())
              .value(incidentStatus.getValue())
              .build();
      incidentStatusOptions.add(option);
    }

    StaticSelectElement selectElement =
        StaticSelectElement.builder()
            .actionId("create_incident_modal_set_incident_status")
            .placeholder(
                PlainTextObject.builder()
                    .text("Select status of the incident...")
                    .emoji(true)
                    .build())
            .options(incidentStatusOptions)
            .build();

    PlainTextObject label =
        PlainTextObject.builder().text(":arrows_counterclockwise: Status").emoji(true).build();

    InputBlock statusInputBlock =
        Blocks.input(block -> block.blockId("incident_status").element(selectElement).label(label));

    return statusInputBlock;
  }

  private InputBlock getChannelSelectBlock() {
    return Blocks.input(
        i ->
            i.blockId("create_incident_modal_conversation_select")
                .label(new PlainTextObject("Select a channel to post the incident details", false))
                .element(
                    BlockElements.conversationsSelect(
                        cs ->
                            cs.actionId("create_incident_modal_select_conversation")
                                .responseUrlEnabled(true)
                                .placeholder(new PlainTextObject("Select a channel...", false)))));
  }

  /**
   * Fetches the list of incident types as Slack option objects.
   *
   * @return List of Slack OptionObjects representing incident types.
   */
  private List<OptionObject> getIncidentTypes() {
    List<OptionObject> options = new ArrayList<>();

    // Dynamically add options based on the Type enum
    for (Type type : Type.values()) {
      options.add(
          OptionObject.builder()
              .text(new PlainTextObject(type.name(), false))
              .value(type.name())
              .build());
    }

    return options;
  }

  /**
   * Fetches the list of incident severities as Slack option objects.
   *
   * @return List of Slack OptionObjects representing incident severities.
   */
  private List<OptionObject> getIncidentSeverities() {
    List<OptionObject> options = new ArrayList<>();

    // Dynamically add options based on the Severity enum
    for (Severity severity : Severity.values()) {
      options.add(
          OptionObject.builder()
              .text(
                  new PlainTextObject(
                      severity.getDescription(), false)) // Using description from enum
              .value(severity.name()) // Using name of the enum as value
              .build());
    }

    return options;
  }

  /**
   * Fetches the list of incident roles as Slack option objects.
   *
   * @return List of Slack OptionObjects representing incident roles.
   */
  private List<OptionObject> getIncidentRoles() {
    List<OptionObject> options = new ArrayList<>();

    // Dynamically add options based on the RoleType enum
    for (RoleType role : RoleType.values()) {
      options.add(
          OptionObject.builder()
              .text(
                  new PlainTextObject(
                      role.name().replace("_", " "),
                      false)) // Converting enum name to human-readable text
              .value(role.name()) // Using name of the enum as value
              .build());
    }

    return options;
  }

  public void listIncidents(GlobalShortcutContext ctx, SlackIncidentType slackIncidentType)
      throws Exception {
    try {
      Query query;
      if (slackIncidentType.equals(SlackIncidentType.Open)) {
        Criteria criteria = Criteria.where("status").ne("Resolved");
        query = new Query(criteria);
      } else {
        Criteria criteria = Criteria.where("status").is("Resolved");
        query = new Query(criteria);
      }
      List<Incident> listIncidents = incidentService.listIncidents(query);

      // Build blocks to display in the modal
      List<LayoutBlock> blocks = new ArrayList<>();

      if (listIncidents.isEmpty()) {
        // No incidents
        String text = ":information_source: No incidents found.";
        blocks.add(createSectionBlock(text));
      } else {
        // Incident found, adding details
        for (Incident incident : listIncidents) {
          String commander = getCommander(incident);
          String text =
              String.format(
                  ":writing_hand: *Name:* %s\n:vertical_traffic_light: *Severity:* %s\n:firefighter: *Commander:* %s\n:eyes: *Current Status:* %s\n\n",
                  incident.getName(), incident.getSeverity(), commander, incident.getStatus());

          SectionBlock sectionBlock = createSectionBlockWithButton(text, incident.getIdentifier());
          blocks.add(sectionBlock);
          blocks.add(new DividerBlock());
        }
      }

      // Build the modal view
      View modalView =
          View.builder()
              .type("modal")
              .callbackId("incident_list_modal")
              .title(ViewTitle.builder().type("plain_text").text("📋 Incident List").build())
              .blocks(blocks)
              .build();

      // Open the modal
      ViewsOpenResponse response =
          ctx.client().viewsOpen(r -> r.triggerId(ctx.getTriggerId()).view(modalView));

      if (response.isOk()) {
        logger.info("ListIncidents view opened successfully.");
      } else {
        logger.error("Failed to open the list of incidents: {}", response.getError());
      }

      ctx.ack();
    } catch (Exception e) {
      logger.error("Failed to open list incidents view modal: {}", e.getMessage(), e);
      //      ctx.ackWithJson(errorResponse("some error occurred"));
    }
  }

  private List<Incident> getIncidentsForSlackView(SlackIncidentType slackIncidentType) {
    // This method should fetch the incidents based on the slackIncidentType
    // For now, returning an empty list or mocked data as an example
    return new ArrayList<>();
  }

  private String getCommander(Incident incident) {
    // Extract the incident commander from the incident roles
    for (Role role : incident.getRoles()) {
      if (role.getRoleType().equals(RoleType.Incident_Commander)) {
        return "<@" + role.getUserDetails().getUserId() + ">";
      }
    }
    return "N/A";
  }

  private SectionBlock createSectionBlock(String text) {
    MarkdownTextObject markdownText = MarkdownTextObject.builder().text(text).build();
    return SectionBlock.builder().text(markdownText).build();
  }

  private SectionBlock createSectionBlockWithButton(String text, String incidentId) {
    MarkdownTextObject markdownText = MarkdownTextObject.builder().text(text).build();

    PlainTextObject buttonText =
        PlainTextObject.builder().text("🔍 View Details").emoji(true).build();

    ButtonElement buttonElement =
        ButtonElement.builder()
            .actionId("view_incident_" + incidentId)
            .text(buttonText)
            .value(incidentId)
            .build();

    return SectionBlock.builder().text(markdownText).accessory(buttonElement).build();
  }
}
