package socketmode

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	slackincident "github.com/respondnow/respond/server/clients/slack/modals/incident"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
)

func ConnectSlackInSocketMode() error {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return fmt.Errorf("SLACK_APP_TOKEN must be set")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return fmt.Errorf("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("SLACK_BOT_TOKEN must be set")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		return fmt.Errorf("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
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

	// Handle all SlashCommand
	socketmodeHandler.Handle(socketmode.EventTypeSlashCommand, middlewareSlashCommand)
	socketmodeHandler.HandleSlashCommand("/rocket", middlewareSlashCommand)

	// socketmodeHandler.HandleDefault(middlewareDefault)

	socketmodeHandler.RunEventLoopContext(context.Background())

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
		}
	default:
		client.Debugf("unsupported Events API event received")
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
		client.Debugf("button clicked!")
	case slack.InteractionTypeShortcut:
		middlewareInteractionTypeShortcut(evt, client)
	case slack.InteractionTypeViewSubmission:
		// See https://api.slack.com/apis/connections/socket-implement#modal
		middlewareInteractionTypeViewSubmission(evt, client)
	case slack.InteractionTypeDialogSubmission:
	default:

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

	switch viewSubmission.View.CallbackID {
	case "create_incident_modal":
		slackincident.NewIncidentService(client).CreateIncident(evt)
	default:
		logrus.Infof("unsupported viewSubmission callback received: %s", viewSubmission.CallbackID)
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
	default:
		logrus.Infof("unsupported shortcut callback received: %s", shortcut.CallbackID)
	}
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
		default:
			logrus.Infof("unsupported blockAction callback received: %s", blockAction.ActionID)
		}
	}
}

func middlewareSlashCommand(evt *socketmode.Event, client *socketmode.Client) {
	cmd, ok := evt.Data.(slack.SlashCommand)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}

	client.Debugf("Slash command received: %+v", cmd)

	payload := map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewSectionBlock(
				&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "foo",
				},
				nil,
				slack.NewAccessory(
					slack.NewButtonBlockElement(
						"",
						"somevalue",
						&slack.TextBlockObject{
							Type: slack.PlainTextType,
							Text: "bar",
						},
					),
				),
			),
		}}
	client.Ack(*evt.Request, payload)
}

// func middlewareDefault(evt *socketmode.Event, client *socketmode.Client) {
// 	// fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
// }
