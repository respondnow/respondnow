package incident

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type IncidentService interface {
	CreateIncidentView(evt *socketmode.Event)
	CreateIncident(evt *socketmode.Event)
	HandleJoinChannelAction(evt *socketmode.Event, blockAction *slack.BlockAction)
}

type incidentService struct {
	client *socketmode.Client
}

func NewIncidentService(client *socketmode.Client) IncidentService {
	return &incidentService{
		client: client,
	}
}
