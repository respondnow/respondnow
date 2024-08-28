package slackclient

import (
	"context"
	"fmt"
	"log"
	"os"

	slackincident "github.com/respondnow/respond/server/clients/slack/modals/incident"
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
