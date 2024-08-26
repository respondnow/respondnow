package bot

import (
	"context"
	"fmt"

	"github.com/respondnow/respond/server/config"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

var BotUserID string

func SetBotUserID(ctx context.Context, client *socketmode.Client) error {
	botDetails, err := GetBotDetails(ctx, client)
	if err != nil {
		return err
	}
	if len(botDetails.UserID) == 0 {
		return fmt.Errorf("failed to set bot user ID: %s", "bot user ID is emty")
	}
	BotUserID = botDetails.UserID

	return nil
}

func GetBotDetails(ctx context.Context, client *socketmode.Client) (*slack.AuthTestResponse, error) {
	authTestDetails, err := client.AuthTestContext(ctx)
	if err != nil {
		return authTestDetails, err
	}

	return authTestDetails, nil
}

func AddBotUserToIncidentChannel(ctx context.Context, client *socketmode.Client, channelID string) error {
	userList, err := ListUsers(ctx, client, channelID)
	if err != nil {
		return err
	}

}

func ListUsers(ctx context.Context, client *socketmode.Client, channelID string) ([]string, error) {
	userList := make([]string, 0)
	users, nextCur, err := client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
		ChannelID: channelID,
		Limit:     1000,
	})
	if err != nil {
		return userList, err
	}
	userList = append(userList, users...)

	for len(nextCur) > 0 {
		users, nextCur, err = client.GetUsersInConversationContext(ctx, &slack.GetUsersInConversationParameters{
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

func ListChannels(ctx context.Context, client *socketmode.Client) ([]slack.Channel, error) {
	channelList := make([]slack.Channel, 0)
	channels, nextCur, err := client.GetConversationsContext(ctx, &slack.GetConversationsParameters{
		ExcludeArchived: true,
		Limit:           1000,
	})
	if err != nil {
		return channelList, err
	}
	channelList = append(channelList, channels...)

	for len(nextCur) > 0 {
		channels, nextCur, err = client.GetConversationsContext(ctx, &slack.GetConversationsParameters{
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
