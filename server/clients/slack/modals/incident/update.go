package incident

import (
	"context"
	"fmt"
	"strings"

	"github.com/respondnow/respond/server/config"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	incidentdb "github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respond/server/pkg/incident"
	"github.com/respondnow/respond/server/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func (is incidentService) UpdateIncidentSummary(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	logrus.Infof("callback data for update: %+v\n", callback)

	is.client.Ack(*evt.Request)

	responseSubmitted := callback.View.State.Values
	slackUser := callback.User

	if responseSubmitted == nil {
		logrus.Errorf("failed to process update incident summary request from slack: %s", "empty response submitted")
		return
	}

	summary := responseSubmitted["create_incident_modal_summary"]["create_incident_modal_set_summary"]

	incidentIdentifier := callback.View.PrivateMetadata

	incidentService := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
		config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		config.EnvConfig.DefaultHierarchy.DefaultProjectId,
	)

	updatedIncident, err := incidentService.UpdateSummary(context.TODO(), incidentIdentifier, summary.Value, utils.UserDetails{
		Email:    slackUser.Profile.Email,
		Name:     slackUser.Profile.DisplayName,
		UserName: slackUser.Name,
		UserId:   slackUser.ID,
		Source:   utils.Slack,
	})
	if err != nil {
		logrus.Errorf("failed to update incident summary: %+v", err)
		return
	}
	logrus.Infof("Incident summary has been updated: %+v", updatedIncident)
	logrus.Infof("Incident channels for updated summary: %+v\n", updatedIncident.Channels[0].ID)

	err = is.sendUpdateSummaryResponseMsg(updatedIncident.Channels[0].ID, updatedIncident, summary.Value)
	if err != nil {
		logrus.Errorf("failed to post summary update to the channel: %s, error: %+v",
			updatedIncident.IncidentChannel.Slack.ChannelID, err)
		return
	} else {
		logrus.Infof("Summary update confirmation successfully posted to channel: %s", updatedIncident.IncidentChannel.Slack.ChannelID)
	}
}

func (is incidentService) AddIncidentComment(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	logrus.Infof("callback data for update: %+v\n", callback)

	is.client.Ack(*evt.Request)

	responseSubmitted := callback.View.State.Values
	slackUser := callback.User

	if responseSubmitted == nil {
		logrus.Errorf("failed to process update incident comment update request from slack: %s",
			"empty response submitted")
		return
	}

	comment := responseSubmitted["create_incident_modal_comment"]["create_incident_modal_set_comment"]

	incidentIdentifier := callback.View.PrivateMetadata

	incidentService := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
		config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		config.EnvConfig.DefaultHierarchy.DefaultProjectId,
	)

	updatedIncident, err := incidentService.UpdateComment(context.TODO(), incidentIdentifier, comment.Value, utils.UserDetails{
		Email:    slackUser.Profile.Email,
		Name:     slackUser.Profile.DisplayName,
		UserName: slackUser.Name,
		UserId:   slackUser.ID,
		Source:   utils.Slack,
	})
	if err != nil {
		logrus.Errorf("failed to update incident comment: %+v", err)
		return
	}
	logrus.Infof("Incident comment has been updated: %+v", updatedIncident)
	logrus.Infof("Incident channels for updated comment: %+v\n", updatedIncident.Channels[0].ID)

	err = is.sendUpdateCommentResponseMsg(updatedIncident.Channels[0].ID, updatedIncident, comment.Value)
	if err != nil {
		logrus.Errorf("failed to post summary update to the channel: %s, error: %+v",
			updatedIncident.IncidentChannel.Slack.ChannelID, err)
		return
	} else {
		logrus.Infof("Comment update confirmation successfully posted to channel: %s", updatedIncident.IncidentChannel.Slack.ChannelID)
	}
}

func (is incidentService) UpdateIncidentStatus(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	logrus.Infof("callback data for update: %+v\n", callback)

	is.client.Ack(*evt.Request)

	responseSubmitted := callback.View.State.Values
	slackUser := callback.User

	if responseSubmitted == nil {
		logrus.Errorf("failed to process update incident status request from slack: %s",
			"empty response submitted")
		return
	}

	status := responseSubmitted["incident_status"]["create_incident_modal_set_incident_status"]
	logrus.Infof("new status is: %v\n", incidentdb.Status(status.SelectedOption.Value))

	incidentIdentifier := callback.View.PrivateMetadata

	incidentService := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
		config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		config.EnvConfig.DefaultHierarchy.DefaultProjectId,
	)

	updatedIncident, err := incidentService.UpdateStatus(context.TODO(), incidentIdentifier, status.SelectedOption.Value,
		utils.UserDetails{
			Email:    slackUser.Profile.Email,
			Name:     slackUser.Profile.DisplayName,
			UserName: slackUser.Name,
			UserId:   slackUser.ID,
			Source:   utils.Slack,
		})
	if err != nil {
		logrus.Errorf("failed to update incident status: %+v", err)
		return
	}
	logrus.Infof("Incident status has been updated: %+v channel: %v", updatedIncident, updatedIncident.Channels[0].ID)

	err = is.sendUpdateStatusResponseMsg(updatedIncident.Channels[0].ID, updatedIncident, status.SelectedOption.Value)
	if err != nil {
		logrus.Errorf("failed to post status update confirmation to the channel: %s, error: %+v",
			updatedIncident.IncidentChannel.Slack.ChannelID, err)
		return
	} else {
		logrus.Infof("Status update confirmation successfully posted to channel: %s", updatedIncident.IncidentChannel.Slack.ChannelID)
	}
}

func (is incidentService) UpdateIncidentRole(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	logrus.Infof("callback data for update: %+v\n", callback)

	is.client.Ack(*evt.Request)

	responseSubmitted := callback.View.State.Values
	slackUser := callback.User

	incidentIdentifier := callback.View.PrivateMetadata

	if responseSubmitted == nil {
		logrus.Errorf("failed to process update incident status request from slack: %s",
			"empty response submitted")
		return
	}
	logrus.Infof("responseSubmitted: %+v\n", responseSubmitted)

	rolesData := callback.View.State.Values
	supportedIncidentRoles := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentRoles()

	roleAssignments := make(map[string]utils.UserDetails)

	for _, role := range supportedIncidentRoles {
		roleKey := "create_incident_modal_set_" + string(role)
		for _, roleData := range rolesData {
			if roleInfo, exists := roleData[roleKey]; exists {
				userID := roleInfo.SelectedUser
				user, err := is.client.GetUserInfo(userID)
				if err != nil {
					logrus.Errorf("Failed to fetch user info for userID: %s, error: %v", userID, err)
					continue
				}
				roleAssignments[string(role)] = utils.UserDetails{
					Source:   utils.Slack,
					UserId:   userID,
					UserName: user.Name,
				}
				logrus.Infof("Incident Identifier: %s, Role: %s, Assigned User: %s", incidentIdentifier, role, userID)
			}
		}
	}

	logrus.Infof("Updated incident role assignments: %+v\n", roleAssignments)

	incidentService := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
		config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		config.EnvConfig.DefaultHierarchy.DefaultProjectId,
	)

	updatedIncident, err := incidentService.UpdateRoles(context.TODO(), incidentIdentifier, roleAssignments,
		utils.UserDetails{
			Email:    slackUser.Profile.Email,
			Name:     slackUser.Profile.DisplayName,
			UserName: slackUser.Name,
			UserId:   slackUser.ID,
			Source:   utils.Slack,
		})
	if err != nil {
		logrus.Errorf("failed to update incident status: %+v", err)
		return
	}
	logrus.Infof("Incident status has been updated: %+v channel: %v", updatedIncident, updatedIncident.Channels[0].ID)

	err = is.sendUpdateRolesResponseMsg(updatedIncident.Channels[0].ID, updatedIncident, roleAssignments)
	if err != nil {
		logrus.Errorf("failed to post roles update confirmation to the channel: %s, error: %+v",
			updatedIncident.IncidentChannel.Slack.ChannelID, err)
		return
	} else {
		logrus.Infof("Roles update confirmation successfully posted to channel: %s", updatedIncident.IncidentChannel.Slack.ChannelID)
	}
}

func (is incidentService) UpdateIncidentSeverity(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	logrus.Infof("callback data for update: %+v\n", callback)

	is.client.Ack(*evt.Request)

	responseSubmitted := callback.View.State.Values
	slackUser := callback.User

	if responseSubmitted == nil {
		logrus.Errorf("failed to process update incident severity request from slack: %s", "empty response submitted")
		return
	}

	severity := responseSubmitted["incident_severity"]["create_incident_modal_set_incident_severity"]
	logrus.Infof("new severity is: %v\n", incidentdb.Severity(severity.SelectedOption.Value))

	incidentIdentifier := callback.View.PrivateMetadata

	incidentService := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
		config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		config.EnvConfig.DefaultHierarchy.DefaultProjectId,
	)

	updatedIncident, err := incidentService.UpdateSeverity(context.TODO(), incidentIdentifier, severity.SelectedOption.Value,
		utils.UserDetails{
			Email:    slackUser.Profile.Email,
			Name:     slackUser.Profile.DisplayName,
			UserName: slackUser.Name,
			UserId:   slackUser.ID,
			Source:   utils.Slack,
		})
	if err != nil {
		logrus.Errorf("failed to update incident summary: %+v", err)
		return
	}
	logrus.Infof("Incident summary has been updated: %+v channel: %v", updatedIncident, updatedIncident.Channels[0].ID)

	err = is.sendUpdateSeverityResponseMsg(updatedIncident.Channels[0].ID, updatedIncident, severity.SelectedOption.Value)
	if err != nil {
		logrus.Errorf("failed to post severity update to the channel: %s, error: %+v",
			updatedIncident.IncidentChannel.Slack.ChannelID, err)
		return
	} else {
		logrus.Infof("Severity update done: %s", updatedIncident.IncidentChannel.Slack.ChannelID)
	}
}

func (is incidentService) sendUpdateSummaryResponseMsg(channelID string,
	updatedIncident incidentdb.Incident, newSummary string) error {
	userInfo, err := is.client.GetUserInfo(updatedIncident.AuditDetails.UpdatedBy.UserId)
	if err != nil {
		logrus.Errorf("failed to fetch user info for Slack ID %s: %+v", updatedIncident.AuditDetails.UpdatedBy.UserId, err)
		return err
	}
	slackHandle := userInfo.Name

	messageText := fmt.Sprintf(
		":memo: *Summary Updated*\n <@%s> updated the summary to: _%s_",
		slackHandle,
		newSummary,
	)

	logrus.Infof("send update summary response message to channel %v", channelID)
	_, _, err = is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		logrus.Infof("there is some error while posting the update message back to the slack")
		return err
	}

	return nil
}

func (is incidentService) sendUpdateCommentResponseMsg(channelID string,
	updatedIncident incidentdb.Incident, newComment string) error {
	userInfo, err := is.client.GetUserInfo(updatedIncident.AuditDetails.UpdatedBy.UserId)
	if err != nil {
		logrus.Errorf("failed to fetch user info for Slack ID %s: %+v", updatedIncident.AuditDetails.UpdatedBy.UserId, err)
		return err
	}
	slackHandle := userInfo.Name

	messageText := fmt.Sprintf(
		":memo: *Comment Added*\n <@%s> added the comment: _%s_",
		slackHandle,
		newComment,
	)

	logrus.Infof("send update comment response message to channel %v", channelID)
	_, _, err = is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		logrus.Infof("there is some error while posting the update message back to the slack")
		return err
	}

	return nil
}

func (is incidentService) sendUpdateSeverityResponseMsg(channelID string, updatedIncident incidentdb.Incident, newSeverity string) error {
	userInfo, err := is.client.GetUserInfo(updatedIncident.AuditDetails.UpdatedBy.UserId)
	if err != nil {
		logrus.Errorf("failed to fetch user info for Slack ID %s: %+v", updatedIncident.AuditDetails.UpdatedBy.UserId, err)
		return err
	}
	slackHandle := userInfo.Name

	messageText := fmt.Sprintf(
		":vertical_traffic_light: *Severity Updated*\n <@%s> updated the severity to: _%s_",
		slackHandle,
		newSeverity,
	)

	logrus.Infof("send update severity response message to channel %v", channelID)
	_, _, err = is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		logrus.Infof("there is some error while posting the update message back to the slack")
		return err
	}

	return nil
}

func (is incidentService) sendUpdateStatusResponseMsg(channelID string, updatedIncident incidentdb.Incident,
	newStatus string) error {
	userInfo, err := is.client.GetUserInfo(updatedIncident.AuditDetails.UpdatedBy.UserId)
	if err != nil {
		logrus.Errorf("failed to fetch user info for Slack ID %s: %+v", updatedIncident.AuditDetails.UpdatedBy.UserId, err)
		return err
	}
	slackHandle := userInfo.Name

	messageText := fmt.Sprintf(
		":eyes: *Status Updated*\n <@%s> updated the status to: _%s_",
		slackHandle,
		newStatus,
	)

	logrus.Infof("send update severity response message to channel %v", channelID)
	_, _, err = is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		logrus.Infof("there is some error while posting the update message back to the slack")
		return err
	}

	return nil
}

func (is incidentService) sendUpdateRolesResponseMsg(channelID string,
	updatedIncident incidentdb.Incident, roleAssignments map[string]utils.UserDetails) error {
	var roleUpdates []string
	for role, userDetails := range roleAssignments {
		roleUpdates = append(roleUpdates, fmt.Sprintf("*%s*: <@%s>", role, userDetails.UserId))
	}

	messageText := fmt.Sprintf(
		":firefighter: *Roles Updated*\nThe following roles have been updated:\n%s",
		strings.Join(roleUpdates, "\n"),
	)

	logrus.Infof("send update roles response message to channel %v", channelID)
	_, _, err := is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		logrus.Errorf("there is some error while posting the update message back to the slack")
		return err
	}

	return nil
}
