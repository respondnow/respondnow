package incident

import (
	"context"
	"fmt"

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

	currentIncident, err := incidentService.Get(context.TODO(), incidentIdentifier)
	if err != nil {
		logrus.Errorf("failed to retrieve current incident details: %+v", err)
		return
	}

	oldSummary := currentIncident.Summary

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

	err = is.sendUpdateSummaryResponseMsg(updatedIncident.IncidentChannel.Slack.ChannelID, updatedIncident, oldSummary, summary.Value)
	if err != nil {
		logrus.Errorf("failed to post summary update confirmation to the channel: %s, error: %+v",
			updatedIncident.IncidentChannel.Slack.ChannelID, err)
		return
	} else {
		logrus.Infof("Summary update confirmation successfully posted to channel: %s", updatedIncident.IncidentChannel.Slack.ChannelID)
	}
}

func (is incidentService) sendUpdateSummaryResponseMsg(channelID string,
	updatedIncident incidentdb.Incident, oldSummary,
	newSummary string) error {
	messageText := fmt.Sprintf(
		"*User:* %s updated the incident summary.\n\n*Old Summary:*\n_%s_\n\n*New Summary:*\n_%s_",
		updatedIncident.AuditDetails.UpdatedBy.UserName,
		oldSummary,
		newSummary,
	)

	logrus.Info("send update summary response message")
	_, _, err := is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionText(messageText, false))
	if err != nil {
		logrus.Info("there is some error while posting the update message back to the slack")
		return err
	}

	return nil
}
