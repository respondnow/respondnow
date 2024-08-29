package incident

import (
	"github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type IncidentService interface {
	CreateIncidentView(evt *socketmode.Event)
	CreateIncident(evt *socketmode.Event)
	HandleJoinChannelAction(evt *socketmode.Event, blockAction *slack.BlockAction)
	ListIncidents(evt *socketmode.Event, slackIncidentType incident.SlackIncidentType)
	ShowIncident(evt *socketmode.Event, incidentID string)
}

type incidentService struct {
	client *socketmode.Client
}

func NewIncidentService(client *socketmode.Client) IncidentService {
	return &incidentService{
		client: client,
	}
}
