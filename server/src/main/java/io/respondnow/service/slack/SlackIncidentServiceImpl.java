// package io.respondnow.service.slack;
//
// import com.slack.api.Slack;
// import com.slack.api.app_backend.views.payload.ViewSubmissionPayload;
// import com.slack.api.bolt.context.builtin.GlobalShortcutContext;
// import com.slack.api.bolt.request.builtin.GlobalShortcutRequest;
// import com.slack.api.bolt.request.builtin.ViewSubmissionRequest;
// import com.slack.api.methods.SlackApiException;
// import com.slack.api.methods.request.conversations.ConversationsCreateRequest;
// import com.slack.api.methods.request.conversations.ConversationsInviteRequest;
// import com.slack.api.methods.response.conversations.ConversationsCreateResponse;
// import com.slack.api.methods.response.conversations.ConversationsInviteResponse;
// import com.slack.api.methods.response.views.ViewsOpenResponse;
// import com.slack.api.model.block.*;
// import com.slack.api.model.block.composition.MarkdownTextObject;
// import com.slack.api.model.block.composition.OptionObject;
// import com.slack.api.model.block.composition.PlainTextObject;
// import com.slack.api.model.block.element.BlockElements;
// import com.slack.api.model.view.View;
// import com.slack.api.model.view.ViewState;
// import com.slack.api.model.view.ViewTitle;
// import com.slack.api.model.view.Views;
// import io.respondnow.model.incident.*;
// import io.respondnow.service.incident.IncidentService;
// import java.io.IOException;
// import java.time.Instant;
// import java.util.*;
// import java.util.stream.Collectors;
// import org.slf4j.Logger;
// import org.slf4j.LoggerFactory;
// import org.springframework.beans.factory.annotation.Autowired;
// import org.springframework.beans.factory.annotation.Value;
// import org.springframework.stereotype.Service;
//
// @Service
// public class SlackIncidentServiceImpl implements SlackIncidentService {
//  @Autowired private IncidentService incidentService;
//  @Autowired private SlackService slackService;
//
//  @Value("${hierarchy.defaultAccount.id:default_account_id}")
//  private String defaultAccountId;
//
//  @Value("${hierarchy.defaultOrg.id:default_org_id}")
//  private String defaultOrgId;
//
//  @Value("${hierarchy.defaultProject.id:default_project_id}")
//  private String defaultProjectId;
//
//  private static final Logger logger = LoggerFactory.getLogger(SlackIncidentServiceImpl.class);
//
//  public void createIncident(ViewSubmissionRequest payload) {
//    Map<String, Map<String, ViewState.Value>> stateValues =
//        Optional.ofNullable(payload.getPayload())
//            .map(ViewSubmissionPayload::getView)
//            .map(View::getState)
//            .map(ViewState::getValues)
//            .orElse(Collections.emptyMap());
//
//    // Safely extract values from the payload
//    String incidentIdentifier =
//        Optional.ofNullable(payload.getPayload().getView().getPrivateMetadata())
//            .orElse("Unknown Incident Identifier");
//
//    String name =
//        Optional.ofNullable(stateValues.get("create_incident_modal_name"))
//            .map(inner -> inner.get("create_incident_modal_set_name"))
//            .map(ViewState.Value::getValue)
//            .orElse(null);
//
//    if (name == null) {
//      logger.error("Name field is missing in view submission payload.");
//    }
//
//    String incidentType =
//        Optional.ofNullable(stateValues.get("incident_type"))
//            .map(inner -> inner.get("create_incident_modal_set_incident_type"))
//            .map(ViewState.Value::getSelectedOption)
//            .map(ViewState.SelectedOption::getValue)
//            .orElse(null);
//
//    if (incidentType == null) {
//      logger.error("Incident type field is missing in view submission payload.");
//    }
//
//    String summary =
//        Optional.ofNullable(stateValues.get("create_incident_modal_summary"))
//            .map(inner -> inner.get("create_incident_modal_set_summary"))
//            .map(ViewState.Value::getValue)
//            .orElse(null);
//
//    if (summary == null) {
//      logger.error("Summary field is missing in view submission payload.");
//    }
//
//    String severity =
//        Optional.ofNullable(stateValues.get("incident_severity"))
//            .map(inner -> inner.get("create_incident_modal_set_incident_severity"))
//            .map(ViewState.Value::getSelectedOption)
//            .map(ViewState.SelectedOption::getValue)
//            .orElse(null);
//
//    if (severity == null) {
//      logger.error("Severity field is missing in view submission payload.");
//    }
//
//    String responseChannel =
//        Optional.ofNullable(stateValues.get("create_incident_modal_conversation_select"))
//            .map(inner -> inner.get("create_incident_modal_select_conversation"))
//            .map(ViewState.Value::getSelectedConversation)
//            .orElse(null);
//
//    if (responseChannel == null) {
//      logger.error("Response channel field is missing in view submission payload.");
//    }
//
//    List<String> roles =
//        Optional.ofNullable(stateValues.get("incident_role"))
//            .map(inner -> inner.get("create_incident_modal_set_incident_role"))
//            .map(ViewState.Value::getSelectedOptions) // Get the list of selected options
//            .map(
//                selectedOptions ->
//                    selectedOptions.stream()
//                        .map(
//                            ViewState.SelectedOption
//                                ::getValue) // Extract the value of each selected option
//                        .collect(Collectors.toList())) // Collect the values into a List
//            .orElse(Collections.emptyList()); // Return an empty list if no roles are found
//
//    if (roles.isEmpty()) {
//      logger.error("No roles are selected in view submission payload.");
//    } else {
//      logger.info("Selected Roles: {}", roles);
//    }
//
//    // Logging extracted values
//    logger.info(
//        "Incident Identifier: {}, Name: {}, Incident Type: {}, Summary: {}, Severity: {}, Response
// Channel: {}",
//        incidentIdentifier,
//        name,
//        incidentType,
//        summary,
//        severity,
//        responseChannel);
//
//    // Proceed with further business logic if all required fields are present
//    if (name != null
//        && incidentType != null
//        && summary != null
//        && severity != null
//        && !roles.isEmpty()
//        && responseChannel != null) {
//      logger.info("Creating incident with the provided details.");
//
//      long createdAt = Instant.now().getEpochSecond();
//      String incidentId = incidentService.generateIncidentIdentifier(createdAt);
//
//      Slack slackClient = slackService.getSlackClient();
//      // Create Slack channel using conversations.create
//      try {
//        ConversationsCreateRequest createChannelRequest =
//            ConversationsCreateRequest.builder()
//                .name("incident-" + name + "-" + createdAt) // Channel name
//                .isPrivate(false) // Change to true if you want a private channel
//                .teamId(payload.getPayload().getTeam().getId()) // Team ID from the callback
//                .build();
//
//        ConversationsCreateResponse createChannelResponse =
//            slackClient
//                .methods(slackService.getBotToken())
//                .conversationsCreate(createChannelRequest);
//
//        if (!createChannelResponse.isOk()) {
//          throw new IOException(
//              "Failed to create Slack channel: " + createChannelResponse.getError());
//        }
//
//        String channelId = createChannelResponse.getChannel().getId();
//        System.out.println("Successfully created an incident channel: " + channelId);
//
//        // Invite users to the channel
//        try {
//          ConversationsInviteRequest conversationsInviteRequest =
//              ConversationsInviteRequest.builder()
//                  .channel(channelId)
//                  .users(
//                      Collections.singletonList(
//                          payload.getPayload().getUser().getId())) // User ID from the callback
//                  .build();
//
//          ConversationsInviteResponse conversationsInviteResponse =
//              slackClient
//                  .methods(slackService.getBotToken())
//                  .conversationsInvite(conversationsInviteRequest);
//
//          if (!conversationsInviteResponse.isOk()) {
//            throw new IOException(
//                "Failed to invite user to Slack channel: "
//                    + conversationsInviteResponse.getError());
//          }
//
//          // Create incident record in the database
//          Incident newIncident = new Incident();
//          newIncident.setAccountIdentifier(defaultAccountId);
//          newIncident.setOrgIdentifier(defaultOrgId);
//          newIncident.setProjectIdentifier(defaultProjectId);
//          newIncident.setIdentifier(incidentIdentifier);
//          newIncident.setName(name);
//          newIncident.setSeverity(Severity.valueOf(severity));
//          newIncident.setIncidentChannel(new IncidentChannel());
//
//          Incident incident = incidentService.createIncident(newIncident);
//
//          // Post messages in Slack
//          postIncidentCreationResponse(responseChannel, channelId, newIncident);
//
//        } catch (SlackApiException | IOException e) {
//          throw new RuntimeException(e);
//        }
//
//      } catch (SlackApiException | IOException e) {
//        throw new RuntimeException(e);
//      }
//    } else {
//      logger.error("Failed to create incident due to missing required fields.");
//    }
//  }
//
//  private void postIncidentCreationResponse(
//      String selectedChannelForResponse, String channelId, Incident newIncident)
//      throws SlackApiException, IOException {
//    slackService
//        .getSlackApp()
//        .client()
//        .chatPostMessage(
//            r -> r.channel(channelId).blocks(createIncidentMessageBlocks(newIncident)));
//
//    // Post another response in the incident channel
//    slackService
//        .getSlackApp()
//        .client()
//        .chatPostMessage(
//            r ->
//                r.channel(selectedChannelForResponse)
//                    .blocks(createIncidentMessageBlocks(newIncident)));
//  }
//
//  private List<LayoutBlock> createIncidentMessageBlocks(Incident newIncident) {
//    // Creating Slack message blocks based on incident details
//    return Arrays.asList(
//        SlackBlockFactory.createHeaderBlock(":fire: New Incident", "new_incident_header"),
//        SlackBlockFactory.createSectionBlock(
//            ":writing_hand: *Name:* " + newIncident.getName(), "new_incident_name_section"),
//        SlackBlockFactory.createSectionBlock(
//            ":vertical_traffic_light: *Severity:* " + newIncident.getSeverity(),
//            "new_incident_severity_section"),
//        SlackBlockFactory.createSectionBlock(
//            ":open_book: *Summary:* " + newIncident.getSummary(),
// "new_incident_summary_section"));
//  }
//
//  public void handleIncidentSummaryViewSubmission(ViewSubmissionRequest payload) {
//    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
//    String updatedSummary =
//        payload
//            .getPayload()
//            .getView()
//            .getState()
//            .getValues()
//            .get("create_incident_modal_summary")
//            .get("create_incident_modal_set_summary")
//            .getValue();
//    logger.info("Incident Identifier: {}, Updated Summary: {}", incidentIdentifier,
// updatedSummary);
//  }
//
//  public void handleIncidentCommentViewSubmission(ViewSubmissionRequest payload) {
//    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
//    String updatedComment =
//        payload
//            .getPayload()
//            .getView()
//            .getState()
//            .getValues()
//            .get("create_incident_modal_comment")
//            .get("create_incident_modal_set_comment")
//            .getValue();
//    logger.info("Incident Identifier: {}, Updated Comment: {}", incidentIdentifier,
// updatedComment);
//  }
//
//  public void handleIncidentRolesViewSubmission(ViewSubmissionRequest payload) {
//    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
//    ViewState state = payload.getPayload().getView().getState();
//    Map<String, Map<String, ViewState.Value>> rolesData = state.getValues();
//
//    // For example, iterate over predefined roles and assign users
//    for (String role : rolesData.keySet()) {
//      String roleKey = "create_incident_modal_set_" + role;
//      if (rolesData.containsKey(roleKey)) {
//        String userID = String.valueOf(rolesData.get(roleKey).get("selected_user"));
//        logger.info(
//            "Incident Identifier: {}, Role: {}, Assigned User: {}",
//            incidentIdentifier,
//            role,
//            userID);
//      }
//    }
//  }
//
//  public void handleIncidentStatusViewSubmission(ViewSubmissionRequest payload) {
//    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
//    String status =
//        payload
//            .getPayload()
//            .getView()
//            .getState()
//            .getValues()
//            .get("incident_status")
//            .get("create_incident_modal_set_incident_status")
//            .getSelectedOption()
//            .getValue();
//    logger.info("Incident Identifier: {}, Updated Status: {}", incidentIdentifier, status);
//  }
//
//  public void handleIncidentSeverityViewSubmission(ViewSubmissionRequest payload) {
//    String incidentIdentifier = payload.getPayload().getView().getPrivateMetadata();
//    String severity =
//        payload
//            .getPayload()
//            .getView()
//            .getState()
//            .getValues()
//            .get("incident_severity")
//            .get("create_incident_modal_set_incident_severity")
//            .getSelectedOption()
//            .getValue();
//    logger.info("Incident Identifier: {}, Updated Severity: {}", incidentIdentifier, severity);
//  }
//
//  public void createIncident(GlobalShortcutRequest req, GlobalShortcutContext ctx) {
//    try {
//      // Blocks for the modal
//      List<LayoutBlock> blocks = new ArrayList<>();
//      blocks.add(
//          SectionBlock.builder()
//              .text(
//                  MarkdownTextObject.builder()
//                      .text(
//                          "This will start a new incident channel, and you will "
//                              + "be invited to it. From there, please use our incident "
//                              + "management process to run the incident or coordinate "
//                              + "with others to do so.")
//                      .build())
//              .build());
//      blocks.add(getNameBlock());
//      blocks.add(getTypeBlock());
//      blocks.add(getSummaryBlock());
//      blocks.add(getSeverityBlock());
//      blocks.add(getRoleBlock());
//      blocks.add(getChannelSelectBlock());
//
//      // Build the modal view
//      View modalView =
//          Views.view(
//              v ->
//                  v.type("modal")
//                      .callbackId("create_incident_modal")
//                      .title(
//                          Views.viewTitle(
//                              title -> title.type("plain_text").text("Start a new incident")))
//                      .blocks(blocks)
//                      .submit(Views.viewSubmit(submit ->
// submit.type("plain_text").text("Start"))));
//
//      // Open the modal
//      ViewsOpenResponse response =
//          ctx.client().viewsOpen(r ->
// r.triggerId(req.getPayload().getTriggerId()).view(modalView));
//
//      logger.info("CreateIncident view opened successfully: {}", response.isOk());
//      ctx.ack();
//    } catch (Exception e) {
//      logger.error("Failed to open create incident view modal: {}", e.getMessage(), e);
//      ctx.ackWithJson(errorResponse());
//    }
//  }
//
//  private InputBlock getNameBlock() {
//    return Blocks.input(
//        i ->
//            i.blockId("create_incident_modal_name")
//                .label(new PlainTextObject(":writing_hand: Incident Name", true))
//                .element(
//                    BlockElements.plainTextInput(
//                        pt ->
//                            pt.actionId("create_incident_modal_set_name")
//                                .placeholder(new PlainTextObject("IAM service is down", false))
//                                .maxLength(76))));
//  }
//
//  private InputBlock getTypeBlock() {
//    List<OptionObject> options = getIncidentTypes(); // Fetch supported incident types
//    return Blocks.input(
//        i ->
//            i.blockId("incident_type")
//                .label(new PlainTextObject(":fire: Incident Type", true))
//                .element(
//                    BlockElements.staticSelect(
//                        s ->
//                            s.actionId("create_incident_modal_set_incident_type")
//                                .placeholder(new PlainTextObject("Select incident type...",
// false))
//                                .options(options))));
//  }
//
//  private InputBlock getSummaryBlock() {
//    return Blocks.input(
//        i ->
//            i.blockId("create_incident_modal_summary")
//                .label(new PlainTextObject(":memo: Summary", true))
//                .element(
//                    BlockElements.plainTextInput(
//                        pt ->
//                            pt.actionId("create_incident_modal_set_summary")
//                                .multiline(true)
//                                .placeholder(
//                                    new PlainTextObject(
//                                        "A brief description of the problem.", false)))));
//  }
//
//  private InputBlock getSeverityBlock() {
//    List<OptionObject> options = getIncidentSeverities(); // Fetch severities
//    return Blocks.input(
//        i ->
//            i.blockId("incident_severity")
//                .label(new PlainTextObject(":vertical_traffic_light: Severity", true))
//                .element(
//                    BlockElements.staticSelect(
//                        s ->
//                            s.actionId("create_incident_modal_set_incident_severity")
//                                .placeholder(
//                                    new PlainTextObject(
//                                        "Select severity of the incident...", false))
//                                .options(options))));
//  }
//
//  private InputBlock getRoleBlock() {
//    List<OptionObject> options = getIncidentRoles(); // Fetch roles
//    return Blocks.input(
//        i ->
//            i.blockId("incident_role")
//                .label(new PlainTextObject(":firefighter: Assign role to yourself", true))
//                .element(
//                    BlockElements.multiStaticSelect(
//                        ms ->
//                            ms.actionId("create_incident_modal_set_incident_role")
//                                .placeholder(
//                                    new PlainTextObject("Assign role to yourself...", false))
//                                .options(options))));
//  }
//
//  private InputBlock getChannelSelectBlock() {
//    return Blocks.input(
//        i ->
//            i.blockId("create_incident_modal_conversation_select")
//                .label(new PlainTextObject("Select a channel to post the incident details",
// false))
//                .element(
//                    BlockElements.conversationsSelect(
//                        cs ->
//                            cs.actionId("create_incident_modal_select_conversation")
//                                .responseUrlEnabled(true)
//                                .placeholder(new PlainTextObject("Select a channel...",
// false)))));
//  }
//
//  /**
//   * Fetches the list of incident types as Slack option objects.
//   *
//   * @return List of Slack OptionObjects representing incident types.
//   */
//  private List<OptionObject> getIncidentTypes() {
//    List<OptionObject> options = new ArrayList<>();
//
//    // Dynamically add options based on the Type enum
//    for (Type type : Type.values()) {
//      options.add(
//          OptionObject.builder()
//              .text(new PlainTextObject(type.name(), false))
//              .value(type.name().toLowerCase())
//              .build());
//    }
//
//    return options;
//  }
//
//  /**
//   * Fetches the list of incident severities as Slack option objects.
//   *
//   * @return List of Slack OptionObjects representing incident severities.
//   */
//  private List<OptionObject> getIncidentSeverities() {
//    List<OptionObject> options = new ArrayList<>();
//
//    // Dynamically add options based on the Severity enum
//    for (Severity severity : Severity.values()) {
//      options.add(
//          OptionObject.builder()
//              .text(
//                  new PlainTextObject(
//                      severity.getDescription(), false)) // Using description from enum
//              .value(severity.name().toLowerCase()) // Using name of the enum as value
//              .build());
//    }
//
//    return options;
//  }
//
//  /**
//   * Fetches the list of incident roles as Slack option objects.
//   *
//   * @return List of Slack OptionObjects representing incident roles.
//   */
//  private List<OptionObject> getIncidentRoles() {
//    List<OptionObject> options = new ArrayList<>();
//
//    // Dynamically add options based on the RoleType enum
//    for (RoleType role : RoleType.values()) {
//      options.add(
//          OptionObject.builder()
//              .text(
//                  new PlainTextObject(
//                      role.name().replace("_", " "),
//                      false)) // Converting enum name to human-readable text
//              .value(role.name().toLowerCase()) // Using name of the enum as value
//              .build());
//    }
//
//    return options;
//  }
//
//  private static Map<String, Object> errorResponse() {
//    Map<String, Object> response = new HashMap<>();
//    response.put("text", "Failed to open the incident modal. Please try again later.");
//    return response;
//  }
//
//  public void listIncidents(GlobalShortcutContext ctx, SlackIncidentType status) throws Exception
// {
//    try {
//      // Build and send the list of incidents as a modal or message
//      View modalView =
//          View.builder()
//              .type("modal")
//              .callbackId("incident_list_modal")
//              .title(
//                  ViewTitle.builder()
//                      .type("plain_text")
//                      .text(
//                          status == SlackIncidentType.Open ? "Open Incidents" : "Closed
// Incidents")
//                      .build())
//              .blocks(List.of(/* Add your incident list blocks here */ ))
//              .build();
//
//      // Open the modal
//      ViewsOpenResponse response =
//          ctx.client().viewsOpen(r -> r.triggerId(ctx.getTriggerId()).view(modalView));
//
//      logger.info("ListIncidents view opened successfully: {}", response.isOk());
//      ctx.ack();
//    } catch (Exception e) {
//      logger.error("Failed to open list incidents view modal: {}", e.getMessage(), e);
//      ctx.ackWithJson(errorResponse());
//    }
//  }
// }
