package slackclient

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (s slackService) HandleAppHome(evt *slackevents.AppHomeOpenedEvent) {
	if evt == nil {
		logrus.Errorf("failed to handle app_home_opened event: %s", "nil evt received in handler")
		return
	}
	logrus.Infof("event app_home_opened === %+v", evt)

	slackBlocks := []slack.Block{
		slack.NewHeaderBlock(&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":robot_face: Respond Now",
		}, slack.HeaderBlockOptionBlockID("app_home_resp_header")),

		slack.NewActionBlock("app_home_resp_create_incident_button", slack.ButtonBlockElement{
			Type: slack.METButton,
			Text: &slack.TextBlockObject{
				Type:  slack.PlainTextType,
				Text:  "Start New Incident",
				Emoji: true,
			},
			ActionID: "create_incident_modal",
			Value:    "show_incident_modal",
			Style:    slack.StyleDanger,
		}),

		slack.NewDividerBlock(),

		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "*Hi there, <@" + evt.User +
				"> :wave:*!\n\nI'm your friendly Respond Now, and my " +
				"sole purpose is to help us manage incidents.\n",
		}, nil, nil, slack.SectionBlockOptionBlockID("app_home_resp_intro")),

		slack.NewHeaderBlock(&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":slack: Adding me to a channel",
		}, slack.HeaderBlockOptionBlockID("app_home_resp_add_to_channel_header")),

		slack.NewDividerBlock(),

		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf("To add me to a new channel, please use <@%s>", BotUserID),
		}, nil, nil, slack.SectionBlockOptionBlockID("app_home_resp_add_to_channel_steps")),

		slack.NewHeaderBlock(&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":firefighter: Creating New Incidents",
		}, slack.HeaderBlockOptionBlockID("app_home_resp_creating_new_incidents_header")),

		slack.NewDividerBlock(),

		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "To create a new incident, you can do the following:\n" +
				"- Use the 'Start New Incident' button here\n " +
				"- Search for 'start a new incident' in the Slack search bar\n" +
				"- type _/start_ in any Slack channel to find my create command and run it.",
		}, nil, nil, slack.SectionBlockOptionBlockID("app_home_resp_create_incident_steps")),

		slack.NewImageBlock("https://www.shutterstock.com/shutterstock/photos/2451527535/display_1500/stock-photo-incident-management-process-business-technology-concept-businessman-using-laptop-with-incident-2451527535.jpg",
			"how to start a new incident", "start_new_incident_image_block", &slack.TextBlockObject{
				Type:  slack.PlainTextType,
				Text:  "How to start a new incident",
				Emoji: true,
			}),

		slack.NewHeaderBlock(&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":point_right: Documentation and Learning Materials",
		}, slack.HeaderBlockOptionBlockID("app_home_resp_docs_header")),

		slack.NewDividerBlock(),

		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: "I have a lot of features. To check them all out, visit my <https://github.com/respondnow/respondnow/blob/main/README.md|docs>.",
		}, nil, nil, slack.SectionBlockOptionBlockID("app_home_resp_docs_content")),

		// slack.NewHeaderBlock(&slack.TextBlockObject{
		// 	Type: slack.PlainTextType,
		// 	Text: ":point_right: My Commands",
		// }, slack.HeaderBlockOptionBlockID("app_home_resp_cmd_header")),

		// slack.NewDividerBlock(),
	}

	viewResp, err := s.socketModeClient.PublishViewContext(context.TODO(), evt.User, slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: slackBlocks,
		},
	}, "")
	if err != nil {
		logrus.Errorf("failed to handle app_home_opened event: %+v", err)
	} else {
		logrus.Infof("App home view resp: %+v, err: %+v", viewResp, err)
	}
}
