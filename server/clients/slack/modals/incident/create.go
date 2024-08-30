package incident

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/respondnow/respondnow/server/config"

	"github.com/respondnow/respondnow/server/pkg/database/mongodb"
	incidentdb "github.com/respondnow/respondnow/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respondnow/server/pkg/incident"
	"github.com/respondnow/respondnow/server/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

const (
	// How many total characters are allowed in a Slack channel name?
	// Limit the channel name to 76 to take this into account
	slackChannelNameLengthCap = 80
	// How many characters does the incident prefix take up?
	channelNamePrefixLength = len("rn-2006-01-02-15-04-05-")
	// How long can the provided name can be?
	incidentNameMaxLength = slackChannelNameLengthCap - channelNamePrefixLength
)

func getIncidentCommanderRoleDescription() string {
	return "The Incident Commander is the decision maker during a major incident, delegating tasks and listening to input from subject matter experts in order to bring the incident to resolution. They become the highest ranking individual on any major incident call, regardless of their day-to-day rank. Their decisions made as commander are final.\n\nYour job as an Incident Commander is to listen to the call and to watch the incident Slack room in order to provide clear coordination, recruiting others to gather context and details. You should not be performing any actions or remediations, checking graphs, or investigating logs. Those tasks should be delegated.\n\nAn IC should also be considering next steps and backup plans at every opportunity, in an effort to avoid getting stuck without any clear options to proceed and to keep things moving towards resolution.\n\nMore information: https://response.pagerduty.com/training/incident_commander/"
}

func getCommunicationsLeadRoleDescription() string {
	return "The purpose of the Communications Liaison is to be the primary individual in charge of notifying our customers of the current conditions, and informing the Incident Commander of any relevant feedback from customers as the incident progresses.\n\nIt's important for the rest of the command staff to be able to focus on the problem at hand, rather than worrying about crafting messages to customers.\n\nYour job as Communications Liaison is to listen to the call, watch the incident Slack room, and track incoming customer support requests, keeping track of what's going on and how far the incident is progressing (still investigating vs close to resolution).\n\nThe Incident Commander will instruct you to notify customers of the incident and keep them updated at various points throughout the call. You will be required to craft the message, gain approval from the IC, and then disseminate that message to customers.\n\nMore information: https://response.pagerduty.com/training/customer_liaison/"
}

func getUserRoleNotificationBlocks(role incidentdb.RoleType, channel string) []slack.Block {
	slackBlocks := []slack.Block{
		slack.NewHeaderBlock(&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: fmt.Sprintf(":wave: You have been elected as the %s for an incident.", role),
		}, slack.HeaderBlockOptionBlockID("user_role_notification_header")),
	}

	if role == incidentdb.IncidentCommander {
		slackBlocks = append(slackBlocks, slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: getIncidentCommanderRoleDescription(),
		}, nil, nil, slack.SectionBlockOptionBlockID("user_role_notification_description")))
	} else if role == incidentdb.CommunicationsLead {
		slackBlocks = append(slackBlocks, slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: getCommunicationsLeadRoleDescription(),
		}, nil, nil, slack.SectionBlockOptionBlockID("user_role_notification_description")))
	}

	slackBlocks = append(slackBlocks,

		slack.NewDividerBlock(),

		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf("Please join the channel here: <#%s>", channel),
		}, nil, nil, slack.SectionBlockOptionBlockID("user_role_notification_channel")))

	return slackBlocks
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

func getNameBlock() *slack.InputBlock {
	return slack.NewInputBlock("create_incident_modal_name", slack.NewTextBlockObject(
		slack.PlainTextType, ":writing_hand: Incident Name", false, false,
	), nil, slack.PlainTextInputBlockElement{
		Type:      slack.METPlainTextInput,
		MaxLength: incidentNameMaxLength,
		ActionID:  "create_incident_modal_set_name",
		Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "IAM service is down",
			false, false),
	})
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
		Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "Add a comment.",
			false, false),
	})
}

func getChannelSelectBlock() *slack.InputBlock {
	return slack.NewInputBlock("create_incident_modal_conversation_select", slack.NewTextBlockObject(
		slack.PlainTextType, "Select a channel to post the incident details", false, false,
	), nil, slack.SelectBlockElement{
		Type:               slack.OptTypeChannels,
		InitialChannel:     config.EnvConfig.SlackConfig.IncidentChannelID,
		ActionID:           "create_incident_modal_select_conversation",
		ResponseURLEnabled: true,
	})
}

func getRoleBlock() *slack.InputBlock {
	supportedRoles := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentRoles()
	initialOptionForIncidentRole := string(incidentdb.IncidentCommander)
	incidentRoleOptions := []*slack.OptionBlockObject{}
	for _, role := range supportedRoles {
		incidentRoleOptions = append(incidentRoleOptions, slack.NewOptionBlockObject(
			string(role),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(role), false, false),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(role), false, false),
		))
	}

	return slack.NewInputBlock(
		"incident_role",
		&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":firefighter: Assign role to yourself",
		},
		nil,
		&slack.MultiSelectBlockElement{
			Type:     slack.MultiOptTypeStatic,
			ActionID: "create_incident_modal_set_incident_role",
			Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "Assign role to yourself...",
				false, false),
			InitialOptions: []*slack.OptionBlockObject{
				slack.NewOptionBlockObject(
					initialOptionForIncidentRole,
					slack.NewTextBlockObject(slack.PlainTextType,
						initialOptionForIncidentRole, false, false),
					slack.NewTextBlockObject(slack.PlainTextType,
						initialOptionForIncidentRole, false, false),
				),
			},
			Options: incidentRoleOptions,
		},
	)
}

func getTypeBlock() *slack.InputBlock {
	supportedIncidentTypes := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentTypes()
	initialOptionForIncidentType := string(supportedIncidentTypes[0])
	incidentTypeOptions := []*slack.OptionBlockObject{}
	for _, incidentType := range supportedIncidentTypes {
		incidentTypeOptions = append(incidentTypeOptions, slack.NewOptionBlockObject(
			string(incidentType),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(incidentType), false, false),
			slack.NewTextBlockObject(slack.PlainTextType,
				string(incidentType), false, false),
		))
	}

	return slack.NewInputBlock(
		"incident_type",
		&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":fire: Incident Type",
		},
		nil,
		&slack.SelectBlockElement{
			Type:        slack.OptTypeStatic,
			ActionID:    "create_incident_modal_set_incident_type",
			Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "Select incident type...", false, false),
			InitialOption: slack.NewOptionBlockObject(
				initialOptionForIncidentType,
				slack.NewTextBlockObject(slack.PlainTextType,
					initialOptionForIncidentType, false, false),
				slack.NewTextBlockObject(slack.PlainTextType,
					initialOptionForIncidentType, false, false),
			),
			Options: incidentTypeOptions,
		},
	)
}

// func getStageBlock() *slack.SectionBlock {
// 	supportedIncidentTypes := incidentdb.NewIncidentOperator(mongodb.Operator).GetIncidentTypes()
// 	initialOptionForIncidentType := string(supportedIncidentTypes[0])
// 	incidentTypeOptions := []*slack.OptionBlockObject{}
// 	for _, incidentType := range supportedIncidentTypes {
// 		incidentTypeOptions = append(incidentTypeOptions, slack.NewOptionBlockObject(
// 			string(incidentType),
// 			slack.NewTextBlockObject(slack.PlainTextType,
// 				string(incidentType), false, false),
// 			slack.NewTextBlockObject(slack.PlainTextType,
// 				string(incidentType), false, false),
// 		))
// 	}

// 	return slack.NewSectionBlock(
// 		&slack.TextBlockObject{
// 			Type: slack.MarkdownType,
// 			Text: ":trackball: *Type of the incident*",
// 		},
// 		nil, slack.NewAccessory(
// 			&slack.SelectBlockElement{
// 				Type:        slack.OptTypeStatic,
// 				ActionID:    "create_incident_modal_set_incident_type",
// 				Placeholder: slack.NewTextBlockObject(slack.PlainTextType, "Select incident type...", false, false),
// 				InitialOption: slack.NewOptionBlockObject(
// 					initialOptionForIncidentType,
// 					slack.NewTextBlockObject(slack.PlainTextType,
// 						initialOptionForIncidentType, false, false),
// 					slack.NewTextBlockObject(slack.PlainTextType,
// 						initialOptionForIncidentType, false, false),
// 				),
// 				Options: incidentTypeOptions,
// 			},
// 		), slack.SectionBlockOptionBlockID("incident_type"),
// 	)
// }

func (is incidentService) CreateIncidentView(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	// Send Ack
	is.client.Ack(*evt.Request)

	slackBlocks := []slack.Block{
		slack.NewSectionBlock(
			&slack.TextBlockObject{
				Type: slack.MarkdownType,
				Text: "This will start a new incident channel and you will " +
					"be invited to it. From there, please use our incident " +
					"management process to run the incident or coordinate " +
					"with others to do so.",
			},
			nil, nil,
		),
		getNameBlock(),
		getTypeBlock(),
		getSummaryBlock(),
		getSeverityBlock(),
		getRoleBlock(),
		getChannelSelectBlock(),
	}

	viewResp, err := is.client.OpenView(callback.TriggerID, slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: "create_incident_modal",
		Title:      slack.NewTextBlockObject(slack.PlainTextType, "Start a new incident", false, false),
		Blocks: slack.Blocks{
			BlockSet: slackBlocks,
		},
		Submit: &slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: "Start",
		},
	})
	logrus.Infof("resp: %+v, err: %+v", viewResp, err)
}

func (is incidentService) CreateIncident(evt *socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	// Send Ack
	is.client.Ack(*evt.Request)

	responseSubmitted := callback.View.State.Values
	slackUser := callback.User

	if responseSubmitted == nil {
		logrus.Errorf("failed to process create incident request from slack: %s", "empty response submitted")
		return
	}

	name := responseSubmitted["create_incident_modal_name"]["create_incident_modal_set_name"]
	incidentType := responseSubmitted["incident_type"]["create_incident_modal_set_incident_type"]
	summary := responseSubmitted["create_incident_modal_summary"]["create_incident_modal_set_summary"]
	severity := responseSubmitted["incident_severity"]["create_incident_modal_set_incident_severity"]

	var incidentRoles []incidentdb.Role
	roles := responseSubmitted["incident_role"]["create_incident_modal_set_incident_role"]
	for _, role := range roles.SelectedOptions {
		incidentRoles = append(incidentRoles, incidentdb.Role{
			Type: incidentdb.RoleType(role.Value),
			User: utils.UserDetails{
				Source:   utils.Slack,
				UserId:   slackUser.ID,
				UserName: slackUser.Name,
			},
		})
	}

	selectedChannelForResponse := responseSubmitted["create_incident_modal_conversation_select"]["create_incident_modal_select_conversation"]

	createdAt := time.Now()
	incidentService := incident.NewIncidentService(incidentdb.NewIncidentOperator(mongodb.Operator),
		config.EnvConfig.DefaultHierarchy.DefaultAccountId,
		config.EnvConfig.DefaultHierarchy.DefaultOrgId,
		config.EnvConfig.DefaultHierarchy.DefaultProjectId,
	)
	incidentIdentifier := incidentService.GenerateIncidentIdentifier(createdAt.Unix())
	channel, err := is.client.Client.CreateConversation(slack.CreateConversationParams{
		ChannelName: generateSlackChannelName(name.Value, &createdAt),
		IsPrivate:   false,
		TeamID:      callback.Team.ID,
	})
	if err != nil {
		logrus.Errorf("failed to create incident channel: %+v", err)
		return
	} else {
		logrus.Infof("Successfully created an incident channel: %+v", channel)
	}
	_, err = is.client.Client.InviteUsersToConversation(channel.ID, slackUser.ID)
	if err != nil {
		logrus.Errorf("failed to invite user %s to incident channel %s: %+v", slackUser.ID, channel.ID, err)
		return
	} else {
		logrus.Infof("Successfully invited user %s to incident channel %s", slackUser.ID, channel.ID)
	}

	createIncidentReq := incident.CreateRequest{
		ResourceDetails: mongodb.ResourceDetails{
			Name:        name.Value,
			Identifier:  incidentIdentifier,
			Description: summary.Value,
		},
		Type:     incidentdb.Type(incidentType.SelectedOption.Value),
		Summary:  summary.Value,
		Severity: incidentdb.Severity(severity.SelectedOption.Value),
		Roles:    incidentRoles,
		Status:   incidentdb.Started,
		IncidentChannel: &incidentdb.IncidentChannel{
			Type: incidentdb.ChannelSlack,
			Slack: &incidentdb.Slack{
				SlackTeam: incidentdb.SlackTeam{
					TeamID:     callback.Team.ID,
					TeamName:   callback.Team.Name,
					TeamDomain: callback.Team.Domain,
				},
				ChannelID: config.EnvConfig.SlackConfig.IncidentChannelID,
			},
		},
		Channels: []incidentdb.Channel{
			{
				ID:     channel.ID,
				TeamID: callback.Team.ID,
				Name:   channel.Name,
				Source: incidentdb.SlackSource,
				Status: incidentdb.Operational,
				URL:    fmt.Sprintf("https://%s.slack.com/archives/%s", callback.Team.ID, channel.ID),
			},
		},
	}
	user := utils.UserDetails{
		Email:    slackUser.Profile.Email,
		Name:     slackUser.Profile.DisplayName,
		UserName: slackUser.Name,
		UserId:   slackUser.ID,
		Source:   utils.Slack,
	}
	newIncident, err := incidentService.Create(context.TODO(), createIncidentReq, user, "")
	if err != nil {
		logrus.Errorf("failed to create incident: %+v", err)
		return
	}
	logrus.Infof("A new incident has been created: %+v", newIncident)

	err = is.sendCreateIncidentResponseMsg(slackUser.TeamID, selectedChannelForResponse.SelectedChannel,
		channel.ID, channel.Name, newIncident.Incident)
	if err != nil {
		logrus.Errorf("failed to post message to the channel: %s, error: %+v",
			selectedChannelForResponse.SelectedChannel, err)
		return
	} else {
		logrus.Infof("A new incident creation response successfully posted to channel:%s",
			selectedChannelForResponse.SelectedChannel)
	}

	err = is.sendCreateIncidentResponseMsg(slackUser.TeamID, channel.ID, "", "",
		newIncident.Incident)
	if err != nil {
		logrus.Errorf("failed to post message to the channel: %s, error: %+v", channel.ID, err)
		return
	} else {
		logrus.Infof("A new incident creation response successfully posted to channel:%s", channel.ID)
	}

	var rolesAssigned []incidentdb.RoleType
	for _, role := range newIncident.Roles {
		rolesAssigned = append(rolesAssigned, role.Type)
	}
	for _, role := range rolesAssigned {
		_, _, err := is.client.Client.PostMessageContext(context.TODO(), callback.User.ID,
			slack.MsgOptionBlocks(getUserRoleNotificationBlocks(role, channel.ID)...))
		if err != nil {
			logrus.Errorf("failed to notify user for the role assigned: %+v", err)
		}
	}
}

func (is incidentService) sendCreateIncidentResponseMsg(teamID, channelID, joinChannelID, joinChannelName string,
	newIncident incidentdb.Incident) error {
	var commander string
	for _, role := range newIncident.Roles {
		if role.Type == incidentdb.IncidentCommander {
			commander = role.User.UserId
			break
		}
	}

	blocks := []slack.Block{
		slack.NewHeaderBlock(&slack.TextBlockObject{
			Type: slack.PlainTextType,
			Text: ":fire: :mega: New Incident",
		}, slack.HeaderBlockOptionBlockID("create_incident_channel_resp_header")),

		slack.NewSectionBlock(nil, []*slack.TextBlockObject{
			{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf(":writing_hand: *Name*\n _%s_", newIncident.Name),
			},
			{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf(":vertical_traffic_light: *Severity*\n _%s_", newIncident.Severity),
			},
		}, nil, slack.SectionBlockOptionBlockID("create_incident_channel_resp_name_severity")),

		slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: fmt.Sprintf(":open_book: *Summary*\n _%s_", newIncident.Summary),
		}, nil, nil, slack.SectionBlockOptionBlockID("create_incident_channel_resp_summary")),

		slack.NewDividerBlock(),

		slack.NewSectionBlock(nil, []*slack.TextBlockObject{
			{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf(":eyes: *Current Status*\n _%s_", newIncident.Status),
			},
			{
				Type: slack.MarkdownType,
				Text: fmt.Sprintf(":firefighter: *Commander*\n _<@%s>_", commander),
			},
		}, nil, slack.SectionBlockOptionBlockID("create_incident_channel_resp_status")),
	}

	shouldBePinned := true
	if joinChannelID != "" && joinChannelName != "" {
		shouldBePinned = false
		blocks = append(blocks,
			slack.NewDividerBlock(),

			slack.NewActionBlock("create_incident_channel_join_channel", slack.ButtonBlockElement{
				Type: slack.METButton,
				Text: &slack.TextBlockObject{
					Type: slack.PlainTextType,
					Text: fmt.Sprintf(":slack: Join %s", joinChannelName),
				},
				Style:    slack.StylePrimary,
				URL:      fmt.Sprintf("https://%s.slack.com/archives/%s", teamID, joinChannelID),
				ActionID: "create_incident_channel_join_channel_button",
				Value:    joinChannelID,
			}),
		)
	}

	blocks = append(blocks,
		slack.NewDividerBlock(),

		slack.NewActionBlock("incident_action_buttons",
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Update Summary"},
				ActionID: "update_incident_summary_button",
				Value:    newIncident.Identifier,
			},
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Add a comment"},
				ActionID: "update_incident_comment_button",
				Value:    newIncident.Identifier,
			},
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Assign Roles"},
				ActionID: "update_incident_assign_roles_button",
				Value:    newIncident.Identifier,
			},
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Update Status"},
				ActionID: "update_incident_status_button",
				Value:    newIncident.Identifier,
			},
			slack.ButtonBlockElement{
				Type:     slack.METButton,
				Text:     &slack.TextBlockObject{Type: slack.PlainTextType, Text: "Update Severity"},
				ActionID: "update_incident_severity_button",
				Value:    newIncident.Identifier,
			},
		),
	)

	blocks = append(blocks,
		slack.NewContextBlock("create_incident_channel_resp_createdAt",
			slack.NewTextBlockObject(
				slack.PlainTextType, fmt.Sprintf(":clock1: Started At: %v",
					time.Unix(newIncident.CreatedAt, 0).UTC()), false, false,
			),
			slack.NewTextBlockObject(
				slack.PlainTextType, fmt.Sprintf(":man: Started By: %s", newIncident.CreatedBy.UserName),
				false, false,
			)),
	)

	_, respTS, err := is.client.Client.PostMessageContext(context.TODO(), channelID, slack.MsgOptionBlocks(blocks...))
	if err != nil {
		return err
	}

	if shouldBePinned {
		err = is.client.Client.AddPinContext(context.TODO(), channelID, slack.NewRefToMessage(
			channelID, respTS,
		))
		if err != nil {
			logrus.Errorf("failed to pin message sent at %v in channel: %s", respTS, channelID)
		}
	}

	return nil
}

func generateSlackChannelName(incidentName string, createdAt *time.Time) string {
	fmtDateTime := createdAt.Format("2006-01-02-15-04-05")
	incidentName = strings.ToLower(strings.TrimSpace(incidentName))
	incidentName = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(incidentName, "-")
	incidentName = strings.ReplaceAll(strings.Trim(incidentName, "-"), "--", "-")

	return "rn-" + fmtDateTime + "-" + incidentName
}

func (is incidentService) HandleJoinChannelAction(evt *socketmode.Event, blockAction *slack.BlockAction) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		logrus.Infof("Ignored %+v\n", evt)
		return
	}
	// Send Ack
	is.client.Ack(*evt.Request)

	_, err := is.client.Client.InviteUsersToConversation(blockAction.Value, callback.User.ID)
	if err != nil {
		logrus.Errorf("failed to invite user %s to incident channel %s: %+v", callback.User.ID,
			blockAction.Value, err)
		return
	} else {
		logrus.Infof("Successfully invited user %s to incident channel %s", callback.User.ID, blockAction.Value)
	}
}
