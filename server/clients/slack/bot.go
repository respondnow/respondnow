package slackclient

import (
	"context"
	"fmt"

	"github.com/respondnow/respondnow/server/config"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

var BotUserID string
var BotUserName string

func (s slackService) SetBotUserIDAndName(ctx context.Context) error {
	botDetails, err := s.GetBotDetails(ctx)
	if err != nil {
		return err
	}
	if len(botDetails.UserID) == 0 {
		return fmt.Errorf("failed to set bot user ID: %s", "bot user ID is emty")
	}
	BotUserID = botDetails.UserID
	BotUserName = botDetails.User

	return nil
}

func (s slackService) GetBotDetails(ctx context.Context) (*slack.AuthTestResponse, error) {
	authTestDetails, err := s.client.AuthTestContext(ctx)
	if err != nil {
		return authTestDetails, err
	}

	return authTestDetails, nil
}

func (s slackService) AddBotUserToIncidentChannel(ctx context.Context, botUserID, channelID string) error {
	userList, err := s.ListUsers(ctx, channelID)
	if err != nil {
		return err
	}

	for _, user := range userList {
		if user == botUserID {
			logrus.Infof("Bot user %s already exists in incident channel: %s", botUserID, channelID)
			return nil
		}
	}

	_, _, _, err = s.client.JoinConversationContext(ctx, channelID)
	if err != nil {
		return err
	}

	logrus.Infof("Added bot user %s to incident channel: %s", botUserID, channelID)
	return nil
}

func (s slackService) ListUsers(ctx context.Context, channelID string) ([]string, error) {
	userList := make([]string, 0)
	users, nextCur, err := s.client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     1000,
	})
	if err != nil {
		return userList, err
	}
	userList = append(userList, users...)

	for len(nextCur) > 0 {
		users, nextCur, err = s.client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
			ChannelID: config.EnvConfig.SlackConfig.IncidentChannelID,
			Limit:     1000,
		})
		if err != nil {
			return userList, err
		}
		userList = append(userList, users...)
	}

	return userList, err
}

func (s slackService) ListChannels(ctx context.Context) ([]slack.Channel, error) {
	channelList := make([]slack.Channel, 0)
	channels, nextCur, err := s.client.GetConversationsContext(ctx, &slack.GetConversationsParameters{
		ExcludeArchived: true,
		Limit:           1000,
	})
	if err != nil {
		return channelList, err
	}
	channelList = append(channelList, channels...)

	for len(nextCur) > 0 {
		channels, nextCur, err = s.client.GetConversationsContext(ctx, &slack.GetConversationsParameters{
			ExcludeArchived: true,
			Limit:           1000,
			Cursor:          nextCur,
		})
		if err != nil {
			return channelList, err
		}
		channelList = append(channelList, channels...)
	}

	return channelList, err
}
