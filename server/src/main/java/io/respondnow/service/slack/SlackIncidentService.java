// package io.respondnow.service.slack;
//
// import com.slack.api.bolt.context.builtin.GlobalShortcutContext;
// import com.slack.api.bolt.request.builtin.GlobalShortcutRequest;
// import com.slack.api.bolt.request.builtin.ViewSubmissionRequest;
// import io.respondnow.model.incident.SlackIncidentType;
//
// public interface SlackIncidentService {
//  void createIncident(GlobalShortcutRequest req, GlobalShortcutContext ctx);
//
//  void listIncidents(GlobalShortcutContext ctx, SlackIncidentType status) throws Exception;
//
//  void createIncident(ViewSubmissionRequest viewSubmission);
//
//  void handleIncidentSummaryViewSubmission(ViewSubmissionRequest viewSubmission);
//
//  void handleIncidentCommentViewSubmission(ViewSubmissionRequest viewSubmission);
//
//  void handleIncidentRolesViewSubmission(ViewSubmissionRequest viewSubmission);
//
//  void handleIncidentStatusViewSubmission(ViewSubmissionRequest viewSubmission);
//
//  void handleIncidentSeverityViewSubmission(ViewSubmissionRequest viewSubmission);
// }
