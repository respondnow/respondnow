package slackclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	slackincident "github.com/respondnow/respond/server/clients/slack/modals/incident"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	incidentdb "github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
)

func (s slackService) ConnectSlackInSocketMode(ctx context.Context) error {
	client := socketmode.New(
		s.client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)

	//\\ EventTypeEventsAPI //\\
	// Handle all EventsAPI
	socketmodeHandler.Handle(socketmode.EventTypeEventsAPI, middlewareEventsAPI)

	// Handle a specific event from EventsAPI
	socketmodeHandler.HandleEvents(slackevents.AppMention, middlewareAppMentionEvent)

	//\\ EventTypeInteractive //\\
	// Handle all Interactive Events
	socketmodeHandler.Handle(socketmode.EventTypeInteractive, middlewareInteractive)

	// Handle a specific Interaction
	socketmodeHandler.HandleInteraction(slack.InteractionTypeBlockActions, middlewareInteractionTypeBlockActions)

	socketmodeHandler.RunEventLoopContext(ctx)

	return nil
}

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	logrus.Info("Connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	logrus.Error("Slack Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	logrus.Info("Connected to Slack with Socket Mode.")
}

func middlewareEventsAPI(evt *socketmode.Event, client *socketmode.Client) {
	logrus.Infof("middlewareEventsAPI")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}

	logrus.Infof("Event received: %+v\n", eventsAPIEvent)

	client.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			logrus.Infof("We have been mentionned in %v", ev.Channel)
			_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				logrus.Errorf("failed posting message: %v", err)
			}
		case *slackevents.MemberJoinedChannelEvent:
			logrus.Infof("user %q joined to channel %q", ev.User, ev.Channel)
		case *slackevents.AppHomeOpenedEvent:
			logrus.Infof("App home has been opened by %s", ev.User)
			svc, err := New()
			if err != nil {
				logrus.Errorf("failed to create slack service: %v", err)
				return
			}
			svc.HandleAppHome(ev)
		}
	default:
		logrus.Errorf("unsupported Events API event received: %+v", eventsAPIEvent.Type)
	}
}

func middlewareAppMentionEvent(evt *socketmode.Event, client *socketmode.Client) {

	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		logrus.Infof("Ignored %+v\n", ev)
		return
	}

	fmt.Printf("We have been mentioned in %v\n", ev.Channel)
	_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
	if err != nil {
		logrus.Errorf("failed posting message: %v", err)
	}
}

func middlewareInteractive(evt *socketmode.Event, client *socketmode.Client) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}

	logrus.Infof("Interaction received: %+v\n", callback)

	var payload interface{}

	switch callback.Type {
	case slack.InteractionTypeBlockActions:
		// See https://api.slack.com/apis/connections/socket-implement#button
		middlewareInteractionTypeBlockActions(evt, client)
	case slack.InteractionTypeShortcut:
		middlewareInteractionTypeShortcut(evt, client)
	case slack.InteractionTypeViewSubmission:
		// See https://api.slack.com/apis/connections/socket-implement#modal
		middlewareInteractionTypeViewSubmission(evt, client)
	default:
		logrus.Errorf("unsupported interactive event received: %+v", callback.Type)
	}

	client.Ack(*evt.Request, payload)
}

func middlewareInteractionTypeViewSubmission(evt *socketmode.Event, client *socketmode.Client) {
	viewSubmission, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	client.Debugf("ViewSubmission received: %+v", viewSubmission)

	incidentIdentifier := viewSubmission.View.PrivateMetadata // Retrieve the incident identifier

	switch viewSubmission.View.CallbackID {
	case "create_incident_modal":
		slackincident.NewIncidentService(client).CreateIncident(evt)
	case "incident_summary_modal":
		logrus.Infof("update summary : %+v\n", viewSubmission.View.State.Values)
		summary := viewSubmission.View.State.Values["create_incident_modal_summary"]["create_incident_modal_set_summary"].Value
		logrus.Infof("Incident Identifier: %s, Updated Summary: %s", incidentIdentifier, summary)
		slackincident.NewIncidentService(client).UpdateIncidentSummary(evt)
	case "incident_comment_modal":
		logrus.Infof("update comment : %+v\n", viewSubmission.View.State.Values)
		summary := viewSubmission.View.State.Values["create_incident_modal_comment"]["create_incident_modal_set_comment"].Value
		logrus.Infof("Incident Identifier: %s, Updated Summary: %s", incidentIdentifier, summary)
		slackincident.NewIncidentService(client).UpdateIncidentComment(evt)
	case "incident_roles_modal":
		rolesData := viewSubmission.View.State.Values
		supportedIncidentRoles := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentRoles()

		for _, role := range supportedIncidentRoles {
			roleKey := "create_incident_modal_set_" + string(role)
			for _, roleData := range rolesData {
				if roleInfo, exists := roleData[roleKey]; exists {
					userID := roleInfo.SelectedUser
					logrus.Infof("Incident Identifier: %s, Role: %s, Assigned User: %s", incidentIdentifier, role, userID)
				}
			}
		}
		slackincident.NewIncidentService(client).UpdateIncidentRole(evt)
	case "incident_status_modal":
		status := viewSubmission.View.State.Values["incident_status"]["create_incident_modal_set_incident_status"].SelectedOption.Value
		logrus.Infof("Incident Identifier: %s, Updated Status: %s", incidentIdentifier, status)
		slackincident.NewIncidentService(client).UpdateIncidentStatus(evt)
	case "incident_severity_modal":
		slackincident.NewIncidentService(client).UpdateIncidentSeverity(evt)
		severity := viewSubmission.View.State.Values["incident_severity"]["create_incident_modal_set_incident_severity"].SelectedOption.Value
		logrus.Infof("Incident Identifier: %s, Updated Severity: %s", incidentIdentifier, severity)
	default:
		logrus.Infof("unsupported viewSubmission callback received: %s", viewSubmission.View.CallbackID)
	}
}

func middlewareInteractionTypeShortcut(evt *socketmode.Event, client *socketmode.Client) {
	shortcut, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	client.Debugf("Shortcut received: %+v", shortcut)

	switch shortcut.CallbackID {
	case "open_incident_modal":
		slackincident.NewIncidentService(client).CreateIncidentView(evt)
	case "list_open_incidents_modal":
		slackincident.NewIncidentService(client).ListIncidents(evt, incident.Open)
	case "list_closed_incidents_modal":
		slackincident.NewIncidentService(client).ListIncidents(evt, incident.Closed)
	default:
		logrus.Infof("Unsupported action callback received: %s", shortcut.CallbackID)

	}
}

func getSeverityBlock() *slack.InputBlock {
	supportedIncidentSeverities := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentSeverities()
	initialOptionForIncidentSeverity := string(supportedIncidentSeverities[len(supportedIncidentSeverities)-1])
	incidentSevOptions := []*slack.OptionBlockObject{}
	for _, incidentSev := range supportedIncidentSeverities {
		incidentSevOptions = append(incidentSevOptions, slack.NewOptionBlockObject(
			string(incidentSev),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(incidentSev), true, false),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(incidentSev), true, false),
		))
	}

	return slack.NewInputBlock(
		"incident_severity",
		&slack.TextBlockObject{
			Type:  slack.PlainTextType,
			Text:  ":vertical_traffic_light: Severity",
			Emoji: false,
		},
		nil,
		&slack.SelectBlockElement{
			Type:     slack.OptTypeStatic,
			ActionID: "create_incident_modal_set_incident_severity",
			Placeholder: slack.NewTextBlockObject(slack.PlainTextType,
				"Select severity of the incident...", false, false),
			InitialOption: slack.NewOptionBlockObject(
				initialOptionForIncidentSeverity,
				slack.NewTextBlockObject(slack.PlainTextType,
					initialOptionForIncidentSeverity, false, false),
				slack.NewTextBlockObject(slack.PlainTextType,
					initialOptionForIncidentSeverity, false, false),
			),
			Options: incidentSevOptions,
		},
	)
}

func getStatusBlock() *slack.InputBlock {
	// Fetch the list of statuses from the database or other source
	supportedIncidentStatuses := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentStageStatuses()
	initialOptionForIncidentStatus := string(supportedIncidentStatuses[len(supportedIncidentStatuses)-1])
	incidentStatusOptions := []*slack.OptionBlockObject{}

	for _, incidentStatus := range supportedIncidentStatuses {
		incidentStatusOptions = append(incidentStatusOptions, slack.NewOptionBlockObject(
			string(incidentStatus),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(incidentStatus), true, false),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(incidentStatus), true, false),
		))
	}

	return slack.NewInputBlock(
		"incident_status",
		&slack.TextBlockObject{
			Type:  slack.PlainTextType,
			Text:  ":arrows_counterclockwise: Status",
			Emoji: false,
		},
		nil,
		&slack.SelectBlockElement{
			Type:     slack.OptTypeStatic,
			ActionID: "create_incident_modal_set_incident_status",
			Placeholder: slack.NewTextBlockObject(slack.PlainTextType,
				"Select status of the incident...", false, false),
			InitialOption: slack.NewOptionBlockObject(
				initialOptionForIncidentStatus,
				slack.NewTextBlockObject(slack.PlainTextType,
					initialOptionForIncidentStatus, false, false),
				slack.NewTextBlockObject(slack.PlainTextType,
					initialOptionForIncidentStatus, false, false),
			),
			Options: incidentStatusOptions,
		},
	)
}

func getSummaryBlock() *slack.InputBlock {
	return slack.NewInputBlock("create_incident_modal_summary", slack.NewTextBlockObject(
		slack.PlainTextType, ":memo: Summary", false, false,
	), nil, slack.PlainTextInputBlockElement{
		Type:      slack.METPlainTextInput,
		Multiline: true,
		ActionID:  "create_incident_modal_set_summary",
		Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "A brief description of the problem.",
			false, false),
	})
}

func getCommentBlock() *slack.InputBlock {
	return slack.NewInputBlock("create_incident_modal_comment", slack.NewTextBlockObject(
		slack.PlainTextType, ":speech_balloon: Comment", false, false,
	), nil, slack.PlainTextInputBlockElement{
		Type:      slack.METPlainTextInput,
		Multiline: true,
		ActionID:  "create_incident_modal_set_comment",
		Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "Add a comment",
			false, false),
	})
}

func getRolesBlock() []slack.Block {
	supportedIncidentRoles := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentRoles()
	var blocks []slack.Block
	for _, role := range supportedIncidentRoles {
		roleText := slack.NewTextBlockObject(slack.PlainTextType, string(role), false, false)
		userSelect := slack.NewOptionsSelectBlockElement(
			slack.OptTypeUser,
			slack.NewTextBlockObject(slack.PlainTextType, "Select a user", false, false),
			"create_incident_modal_set_"+string(role),
		)
		section := slack.NewSectionBlock(roleText, nil, slack.NewAccessory(userSelect))
		blocks = append(blocks, section)
	}
	return blocks
}

func middlewareInteractionTypeBlockActions(evt *socketmode.Event, client *socketmode.Client) {
	blockActions, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	client.Debugf("BlockActions received: %+v", blockActions)

	for _, blockAction := range blockActions.ActionCallback.BlockActions {
		switch blockAction.ActionID {
		case "create_incident_channel_join_channel_button":
			slackincident.NewIncidentService(client).HandleJoinChannelAction(evt, blockAction)
		case "create_incident_modal":
			slackincident.NewIncidentService(client).CreateIncidentView(evt)
		case "update_incident_summary_button":
			client.Debugf("Displaying modal for incident summary")
			modalRequest := slack.ModalViewRequest{
				Type:            slack.VTModal,
				PrivateMetadata: blockAction.Value,
				CallbackID:      "incident_summary_modal",
				Title: slack.NewTextBlockObject("plain_text",
					"Update Incident Summary", false, false),
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						getSummaryBlock(),
					},
				},
				Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
			}

			_, err := client.OpenView(blockActions.TriggerID, modalRequest)
			if err != nil {
				logrus.Errorf("Error opening modal: %v", err)
			}
		case "update_incident_comment_button":
			client.Debugf("Displaying modal for incident comment")
			modalRequest := slack.ModalViewRequest{
				Type:            slack.VTModal,
				PrivateMetadata: blockAction.Value,
				CallbackID:      "incident_comment_modal",
				Title: slack.NewTextBlockObject("plain_text",
					"Add comment", false, false),
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						getCommentBlock(),
					},
				},
				Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
			}

			_, err := client.OpenView(blockActions.TriggerID, modalRequest)
			if err != nil {
				logrus.Errorf("Error opening modal: %v", err)
			}
		case "update_incident_assign_roles_button":
			client.Debugf("Displaying modal for assigning roles")
			modalRequest := slack.ModalViewRequest{
				Type:            slack.VTModal,
				PrivateMetadata: blockAction.Value,
				CallbackID:      "incident_roles_modal",
				Title: slack.NewTextBlockObject("plain_text",
					"Assign Incident Roles", false, false),
				Blocks: slack.Blocks{
					BlockSet: getRolesBlock(),
				},
				Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
			}

			_, err := client.OpenView(blockActions.TriggerID, modalRequest)
			if err != nil {
				logrus.Errorf("Error opening modal: %v", err)
			}
		case "update_incident_status_button":
			client.Debugf("Displaying modal for incident severity selection")
			modalRequest := slack.ModalViewRequest{
				Type:            slack.VTModal,
				CallbackID:      "incident_status_modal",
				PrivateMetadata: blockAction.Value,
				Title: slack.NewTextBlockObject("plain_text",
					"Update Incident Status", false, false),
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						getStatusBlock(),
					},
				},
				Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
			}

			_, err := client.OpenView(blockActions.TriggerID, modalRequest)
			if err != nil {
				logrus.Errorf("Error opening modal: %v", err)
			}
		case "update_incident_severity_button":
			client.Debugf("Displaying modal for incident severity selection")
			modalRequest := slack.ModalViewRequest{
				Type:            slack.VTModal,
				CallbackID:      "incident_severity_modal",
				PrivateMetadata: blockAction.Value,
				Title: slack.NewTextBlockObject("plain_text",
					"Update Incident Severity", false, false),
				Blocks: slack.Blocks{
					BlockSet: []slack.Block{
						getSeverityBlock(),
					},
				},
				Submit: slack.NewTextBlockObject("plain_text", "Submit", false, false),
			}

			_, err := client.OpenView(blockActions.TriggerID, modalRequest)
			if err != nil {
				logrus.Errorf("Error opening modal: %v", err)
			}
		default:
			incidentID := blockAction.ActionID
			if strings.HasPrefix(incidentID, "view_incident_") {
				slackincident.NewIncidentService(client).ShowIncident(evt, blockAction.Value)
				return
			} else {
				logrus.Infof("Unsupported action callback received: %s", blockActions.CallbackID)
			}
		}
	}
}
