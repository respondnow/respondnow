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
import com.slack.api.methods.response.conversations.ConversationsCreateResponse;
import com.slack.api.methods.response.conversations.ConversationsInviteResponse;
import com.slack.api.methods.response.views.ViewsOpenResponse;
import com.slack.api.model.Conversation;
import com.slack.api.model.block.Blocks;
import com.slack.api.model.block.InputBlock;
import com.slack.api.model.block.LayoutBlock;
import com.slack.api.model.block.SectionBlock;
import com.slack.api.model.block.composition.MarkdownTextObject;
import com.slack.api.model.block.composition.OptionObject;
import com.slack.api.model.block.composition.PlainTextObject;
import com.slack.api.model.block.element.BlockElements;
import com.slack.api.model.block.element.ButtonElement;
import com.slack.api.model.event.AppHomeOpenedEvent;
import com.slack.api.model.event.AppMentionEvent;
import com.slack.api.model.event.MemberJoinedChannelEvent;
import com.slack.api.model.view.View;
import com.slack.api.model.view.ViewState;
import com.slack.api.model.view.ViewTitle;
import com.slack.api.model.view.Views;
import com.slack.api.socket_mode.SocketModeClient;
import io.respondnow.dto.incident.CreateRequest;
import io.respondnow.model.incident.*;
import io.respondnow.model.user.UserDetails;
import io.respondnow.service.incident.IncidentService;
import java.io.IOException;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.stream.Collectors;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class SlackServiceImpl implements SlackService {
  @Autowired private IncidentService incidentService;

  private final Slack slackClient;
  private final SocketModeClient socketModeSlackClient;
  private final App slackApp;
  private final SocketModeApp socketModeApp;
  private final ExecutorService executorService;
  private String botUserId;
  private static final Logger logger = LoggerFactory.getLogger(SlackServiceImpl.class);

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

  private void handleCreateIncidentViewSubmission() throws RuntimeException {
    try {
      slackApp.viewSubmission(
          "create_incident_modal",
          (payload, ctx) -> {
            logger.info("Received create incident view submission: {}", payload);
            createIncident(payload);
            return ctx.ack();
          });
    } catch (RuntimeException e) {
      throw new RuntimeException(e);
    }
  }

  private void handleIncidentSummaryViewSubmission() {
    slackApp.viewSubmission(
        "incident_summary_modal",
        (payload, ctx) -> {
          logger.debug("Update summary received: {}", payload);
          handleIncidentSummaryViewSubmission(payload);
          return ctx.ack();
        });
  }

  private void handleIncidentCommentViewSubmission() {
    slackApp.viewSubmission(
        "incident_comment_modal",
        (payload, ctx) -> {
          logger.debug("Update comment received: {}", payload);
          handleIncidentCommentViewSubmission(payload);
          return ctx.ack();
        });
  }

  private void handleIncidentRolesViewSubmission() {
    slackApp.viewSubmission(
        "incident_roles_modal",
        (payload, ctx) -> {
          logger.debug("Incident roles modal received: {}", payload);
          handleIncidentRolesViewSubmission(payload);
          return ctx.ack();
        });
  }

  private void handleIncidentStatusViewSubmission() {
    slackApp.viewSubmission(
        "incident_status_modal",
        (payload, ctx) -> {
          logger.debug("Incident status modal received: {}", payload);
          handleIncidentStatusViewSubmission(payload);
          return ctx.ack();
        });
  }

  private void handleIncidentSeverityViewSubmission() {
    slackApp.viewSubmission(
        "incident_severity_modal",
        (payload, ctx) -> {
          logger.debug("Incident severity modal received: {}", payload);
          handleIncidentSeverityViewSubmission(payload);
          return ctx.ack();
        });
  }

  /** Register all block action handlers for the Slack app. */
  private void registerBlockActionHandlers() throws RuntimeException {
    try {
      registerBlockActionChannelJoinButton();
    } catch (RuntimeException e) {
      throw new RuntimeException(e);
    }
  }

  private void registerBlockActionChannelJoinButton() throws RuntimeException {
    // Handle "create_incident_channel_join_channel_button" action
    slackApp.blockAction(
        "create_incident_channel_join_channel_button",
        (req, ctx) -> {
          // Log that the block action is being handled
          System.out.println("Inside block action handler");

          // Retrieve the value from the block action
          String value = req.getPayload().getActions().get(0).getValue(); // The button's value
          System.out.println("Button value: " + value);

          // Check if the response URL is available
          if (req.getPayload().getResponseUrl() != null) {
            // Respond back to the user
            ctx.respond(r -> r.text("You've sent \"" + value + "\" by clicking the button!"));
          }

          // Acknowledge the action (important to avoid timeouts)
          return ctx.ack();
        });
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
        SlackBlockFactory.createActionsBlock(
            "app_home_resp_create_incident_button",
            ButtonElement.builder()
                .text(PlainTextObject.builder().text("Start New Incident").emoji(true).build())
                .actionId("create_incident_modal")
                .value("show_incident_modal")
                .style("danger")
                .build()),
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

                              // Get user details (assume you have a method to fetch UserDetails by
                              // userId)
                              UserDetails userDetails =
                                  fetchUserDetails(
                                      payload.getPayload().getUser().getId(),
                                      payload.getPayload().getUser().getName(),
                                      "",
                                      ChannelSource.Slack);

                              // Create Role object
                              return new Role(roleType, userDetails);
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
        System.out.println("Successfully created an incident channel: " + channelId);

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
          //          createRequest.setAccountIdentifier(defaultAccountId);
          //          createRequest.setOrgIdentifier(defaultOrgId);
          //          createRequest.setProjectIdentifier(defaultProjectId);
          createRequest.setIdentifier(incidentId);
          createRequest.setName(name);
          createRequest.setType(Type.valueOf(incidentType));
          createRequest.setStatus(Status.Started);
          createRequest.setRoles(roles);
          createRequest.setSeverity(Severity.valueOf(severity));
          createRequest.setSummary(summary);
          createRequest.setIncidentChannel(incidentChannel);
          createRequest.setChannels(channels);

          UserDetails userDetails = new UserDetails();
          userDetails.setName(payload.getPayload().getUser().getName());
          userDetails.setUserName(payload.getPayload().getUser().getUsername());
          userDetails.setUserId(payload.getPayload().getUser().getId());
          //          userDetails.setEmail(payload.getPayload().getUser().get);
          userDetails.setSource(ChannelSource.Slack);
          Incident incident = incidentService.createIncident(createRequest, userDetails);

          // Post messages in Slack
          postIncidentCreationResponse(responseChannel, channelId, incident);

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
  private UserDetails fetchUserDetails(
      String userId, String userName, String userEmail, ChannelSource source) {
    UserDetails userDetails = new UserDetails();
    if (!userId.isEmpty()) {
      userDetails.setUserId(userId);
    }
    if (!userName.isEmpty()) {
      userDetails.setUserName(userName);
    }
    if (!userEmail.isEmpty()) {
      userDetails.setEmail(userEmail);
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

  private void postIncidentCreationResponse(
      String selectedChannelForResponse, String channelId, Incident newIncident)
      throws SlackApiException, IOException {
    slackApp
        .client()
        .chatPostMessage(
            r -> r.channel(channelId).blocks(createIncidentMessageBlocks(newIncident)));

    // Post another response in the incident channel
    slackApp
        .client()
        .chatPostMessage(
            r ->
                r.channel(selectedChannelForResponse)
                    .blocks(createIncidentMessageBlocks(newIncident)));
  }

  private List<LayoutBlock> createIncidentMessageBlocks(Incident newIncident) {
    // Creating Slack message blocks based on incident details
    return Arrays.asList(
        SlackBlockFactory.createHeaderBlock(":fire: New Incident", "new_incident_header"),
        SlackBlockFactory.createSectionBlock(
            ":writing_hand: *Name:* " + newIncident.getName(), "new_incident_name_section"),
        SlackBlockFactory.createSectionBlock(
            ":vertical_traffic_light: *Severity:* " + newIncident.getSeverity(),
            "new_incident_severity_section"),
        SlackBlockFactory.createSectionBlock(
            ":open_book: *Summary:* " + newIncident.getSummary(), "new_incident_summary_section"));
  }

  public void handleIncidentSummaryViewSubmission(ViewSubmissionRequest payload) {
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
  }

  public void handleIncidentCommentViewSubmission(ViewSubmissionRequest payload) {
    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    String updatedComment =
        payload
            .getPayload()
            .getView()
            .getState()
            .getValues()
            .get("create_incident_modal_comment")
            .get("create_incident_modal_set_comment")
            .getValue();
    logger.info("Incident Identifier: {}, Updated Comment: {}", incidentIdentifier, updatedComment);
  }

  public void handleIncidentRolesViewSubmission(ViewSubmissionRequest payload) {
    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
    ViewState state = payload.getPayload().getView().getState();
    Map<String, Map<String, ViewState.Value>> rolesData = state.getValues();

    // For example, iterate over predefined roles and assign users
    for (String role : rolesData.keySet()) {
      String roleKey = "create_incident_modal_set_" + role;
      if (rolesData.containsKey(roleKey)) {
        String userID = String.valueOf(rolesData.get(roleKey).get("selected_user"));
        logger.info(
            "Incident Identifier: {}, Role: {}, Assigned User: {}",
            incidentIdentifier,
            role,
            userID);
      }
    }
  }

  public void handleIncidentStatusViewSubmission(ViewSubmissionRequest payload) {
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
    logger.info("Incident Identifier: {}, Updated Status: {}", incidentIdentifier, status);
  }

  public void handleIncidentSeverityViewSubmission(ViewSubmissionRequest payload) {
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
    logger.info("Incident Identifier: {}, Updated Severity: {}", incidentIdentifier, severity);
  }

  public void createIncident(GlobalShortcutRequest req, GlobalShortcutContext ctx) {
    try {
      // Blocks for the modal
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
      blocks.add(getSummaryBlock());
      blocks.add(getSeverityBlock());
      blocks.add(getRoleBlock());
      blocks.add(getChannelSelectBlock());

      // Build the modal view
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

      // Open the modal
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

  private InputBlock getSummaryBlock() {
    return Blocks.input(
        i ->
            i.blockId("create_incident_modal_summary")
                .label(new PlainTextObject(":memo: Summary", true))
                .element(
                    BlockElements.plainTextInput(
                        pt ->
                            pt.actionId("create_incident_modal_set_summary")
                                .multiline(true)
                                .placeholder(
                                    new PlainTextObject(
                                        "A brief description of the problem.", false)))));
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

  private static Map<String, Object> errorResponse() {
    Map<String, Object> response = new HashMap<>();
    response.put("text", "Failed to open the incident modal. Please try again later.");
    return response;
  }

  public void listIncidents(GlobalShortcutContext ctx, SlackIncidentType status) throws Exception {
    try {
      // Build and send the list of incidents as a modal or message
      View modalView =
          View.builder()
              .type("modal")
              .callbackId("incident_list_modal")
              .title(
                  ViewTitle.builder()
                      .type("plain_text")
                      .text(
                          status == SlackIncidentType.Open ? "Open Incidents" : "Closed Incidents")
                      .build())
              .blocks(List.of(/* Add your incident list blocks here */ ))
              .build();

      // Open the modal
      ViewsOpenResponse response =
          ctx.client().viewsOpen(r -> r.triggerId(ctx.getTriggerId()).view(modalView));

      logger.info("ListIncidents view opened successfully: {}", response.isOk());
      ctx.ack();
    } catch (Exception e) {
      logger.error("Failed to open list incidents view modal: {}", e.getMessage(), e);
      ctx.ackWithJson(errorResponse());
    }
  }
}
