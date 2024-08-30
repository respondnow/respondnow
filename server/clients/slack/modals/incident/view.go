package incident

import (
	"context"
	"fmt"
	"time"

	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	incident2 "github.com/respondnow/respond/server/pkg/incident"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func (is incidentService) ListIncidents(evt *socketmode.Event, slackIncidentType incident.SlackIncidentType) {
	is.client.Ack(*evt.Request)

	inf := incident2.NewIncidentService(incident.NewIncidentOperator(mongodb.Operator), "", "", "")

	listIncidents, err := inf.ListIncidentsForSlackView(context.Background(), slackIncidentType)
	if err != nil {
		logrus.Errorf("failed to list incidents: %v", err)
		return
	}

	createMarkdownTextBlock := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject(slack.MarkdownType, text, false, false)
	}

	createPlainTextBlock := func(text string) *slack.TextBlockObject {
		return slack.NewTextBlockObject(slack.PlainTextType, text, false, false)
	}

	var blocks []slack.Block

	if len(listIncidents) == 0 {
		text := ":information_source: No incidents found."
		blocks = append(blocks, slack.NewSectionBlock(createMarkdownTextBlock(text), nil, nil))
	} else {
		for _, i := range listIncidents {
			var commander string
			for _, role := range i.Roles {
				if role.Type == incident.IncidentCommander {
					commander = role.User.UserId
				}
			}

			if commander != "" {
				commander = fmt.Sprintf("<@%s>", commander)
			} else {
				commander = "N/A"
			}

			text := fmt.Sprintf(
				":writing_hand: *Name:* %s\n:vertical_traffic_light: *Severity:* %s\n:firefighter: *Commander:* %s\n:eyes: *Current Status:* %s\n\n",
				i.Name, i.Severity, commander, i.Status)

			sectionBlock := slack.NewSectionBlock(
				createMarkdownTextBlock(text),
				nil, slack.NewAccessory(
					slack.NewButtonBlockElement("view_incident_"+i.Identifier, i.Identifier, createPlainTextBlock("🔍 View Details"))))

			blocks = append(blocks, sectionBlock, slack.NewDividerBlock())
		}
	}

	modalRequest := slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: "view_more_modal",
		Title:      createPlainTextBlock("📋 Incident List"),
		Close:      createPlainTextBlock("Close"),
		Blocks:     slack.Blocks{BlockSet: blocks},
	}

	triggerID := evt.Data.(slack.InteractionCallback).TriggerID
	if _, err := is.client.Client.OpenView(triggerID, modalRequest); err != nil {
		logrus.Errorf("failed to list incidents: %v", err)
	}
}

func (is incidentService) ShowIncident(evt *socketmode.Event, incidentID string) {
	is.client.Ack(*evt.Request)

	inf := incident2.NewIncidentService(incident.NewIncidentOperator(mongodb.Operator), "", "", "")

	getIncident, err := inf.GetIncidentForSlackView(context.Background(), incidentID)
	if err != nil {
		logrus.Errorf("failed to get incident: %v", err)
		return
	}

	if getIncident.Identifier == "" {
		logrus.Errorf("Incident with ID %s not found", incidentID)
		return
	}

	formatUserMention := func(userID string) string {
		if userID != "" {
			return fmt.Sprintf("<@%s>", userID)
		}
		return "N/A"
	}

	var commander, communicationsLead string
	for _, role := range getIncident.Roles {
		switch role.Type {
		case incident.IncidentCommander:
			commander = role.User.UserId
		case incident.CommunicationsLead:
			communicationsLead = role.User.UserId
		}
	}

	commander = formatUserMention(commander)
	communicationsLead = formatUserMention(communicationsLead)
	startedAt := time.Unix(getIncident.CreatedAt, 0).UTC().Format("Monday, Jan 2, 2006 at 3:04 PM")

	blocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, "*Incident Details*", false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":writing_hand: *Name:* %s", getIncident.Name), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":vertical_traffic_light: *Severity:* %s", string(getIncident.Severity)), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":eyes: *Current Status:* %s", string(getIncident.Status)), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":firefighter: *Commander:* %s", commander), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":phone: *Communications Lead:* %s", communicationsLead), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":open_book: *Summary:* %s", getIncident.Summary), false, false),
			nil,
			nil,
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":clock1: *Started At:* %s", startedAt), false, false),
			nil,
			nil,
		),
	}

	if getIncident.Status == incident.Resolved && getIncident.UpdatedAt != nil {
		completedAt := time.Unix(*getIncident.UpdatedAt, 0).UTC().Format("Monday, Jan 2, 2006 at 3:04 PM")
		blocks = append(blocks, slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(":checkered_flag: *Completed At:* %s", completedAt), false, false),
			nil,
			nil,
		))
	}

	blocks = append(blocks, slack.NewDividerBlock())

	modal := slack.ModalViewRequest{
		Type:       slack.ViewType("modal"),
		CallbackID: "incident_details_modal",
		Title:      slack.NewTextBlockObject("plain_text", "🚨 Incident Details", false, false),
		Blocks:     slack.Blocks{BlockSet: blocks},
		Close:      slack.NewTextBlockObject("plain_text", "Close", false, false),
	}

	// Open the modal view
	triggerID := evt.Data.(slack.InteractionCallback).TriggerID
	if _, err = is.client.Client.PushViewContext(context.Background(), triggerID, modal); err != nil {
		logrus.Errorf("failed to open incident details view: %v", err)
	}
}
