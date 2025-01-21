package io.respondnow.service.slack;

import com.slack.api.Slack;
import com.slack.api.bolt.App;
import com.slack.api.bolt.context.builtin.GlobalShortcutContext;
import com.slack.api.bolt.request.builtin.GlobalShortcutRequest;
import com.slack.api.bolt.request.builtin.ViewSubmissionRequest;
import com.slack.api.methods.SlackApiException;
import com.slack.api.model.Conversation;
import com.slack.api.model.event.AppHomeOpenedEvent;
import io.respondnow.model.incident.SlackIncidentType;

import java.io.IOException;
import java.util.List;

public interface SlackService {
  void setBotUserIDAndName() throws Exception;

  String getBotUserId() throws Exception;

  void addBotUserToIncidentChannel(String botUserID, String channelID) throws Exception;

  List<String> listAllMembersOfChannel(String channelId) throws Exception;

  List<String> listUsers(String channelID) throws Exception;

  List<Conversation> listChannels() throws Exception;

  void handleAppHome(AppHomeOpenedEvent event) throws InterruptedException;

  String getIncidentChannelID();

  String getBotToken();

  String getAppToken();

  boolean isBotInChannel(String botUserID, String channelID) throws Exception;

  void startApp();

  Slack getSlackClient();

  App getSlackApp();

  void createIncident(GlobalShortcutRequest req, GlobalShortcutContext ctx);

  void listIncidents(GlobalShortcutContext ctx, SlackIncidentType status) throws Exception;

  void createIncident(ViewSubmissionRequest viewSubmission);

  void handleIncidentSummaryViewSubmission(ViewSubmissionRequest viewSubmission) throws SlackApiException, IOException;

  void handleIncidentCommentViewSubmission(ViewSubmissionRequest viewSubmission);

  void handleIncidentRolesViewSubmission(ViewSubmissionRequest viewSubmission);

  void handleIncidentStatusViewSubmission(ViewSubmissionRequest viewSubmission);

  void handleIncidentSeverityViewSubmission(ViewSubmissionRequest viewSubmission);
}
