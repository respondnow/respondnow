package slackclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

type SlackService interface {
	ConnectSlackInSocketMode(ctx context.Context) error
	SetBotUserID(ctx context.Context) error
	GetBotDetails(ctx context.Context) (*slack.AuthTestResponse, error)
	AddBotUserToIncidentChannel(ctx context.Context, botUserID, channelID string) error
	ListUsers(ctx context.Context, channelID string) ([]string, error)
	ListChannels(ctx context.Context) ([]slack.Channel, error)
}

type slackService struct {
	client *slack.Client
}

func New() (SlackService, error) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return nil, fmt.Errorf("SLACK_APP_TOKEN must be set")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, fmt.Errorf("SLACK_APP_TOKEN must have the prefix \"xapp-\"")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN must be set")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN must have the prefix \"xoxb-\"")
	}

	return &slackService{
		client: slack.New(
			botToken,
			slack.OptionDebug(true),
			slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
			slack.OptionAppLevelToken(appToken),
		),
	}, nil
}
