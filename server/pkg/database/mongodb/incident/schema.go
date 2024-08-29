package incident

import (
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"github.com/respondnow/respond/server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	Started       Status = "Started"
	Acknowledged  Status = "Acknowledged"
	Investigating Status = "Investigating"
	Identified    Status = "Identified"
	Mitigated     Status = "Mitigated"
	Resolved      Status = "Resolved"

	// Detected            Status = "detected"
	// PostmortemStarted   Status = "postmortemStarted"
	// PostmortemCompleted Status = "PostmortemCompleted"
)

type Severity string

const (
	Severity0 Severity = "SEV0 - Critical, High Impact"
	Severity1 Severity = "SEV1 - Major, Significant Impact"
	Severity2 Severity = "SEV2 - Minor, Low Impact"
)

type Type string

const (
	Availability Type = "Availability"
	Latency      Type = "Latency"
	Security     Type = "Security"
	Other        Type = "Other"
)

type Incident struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
	mongodb.ResourceDetails   `bson:",inline" json:",inline"`
	mongodb.IdentifierDetails `bson:",inline" json:",inline"`
	Type                      Type             `bson:"type" json:"type"`
	Severity                  Severity         `bson:"severity" json:"severity" binding:"required"`
	Status                    Status           `bson:"status" json:"status" binding:"required"`
	Summary                   string           `bson:"summary" json:"summary" binding:"required"`
	Active                    bool             `bson:"active" json:"active" binding:"required"`
	Services                  []Service        `bson:"services,omitempty" json:"services,omitempty"`
	Environments              []Environment    `bson:"environments,omitempty" json:"environments,omitempty"`
	Functionalities           []Functionality  `bson:"functionalities,omitempty" json:"functionalities,omitempty"`
	Roles                     []Role           `bson:"roles,omitempty" json:"roles,omitempty"`
	Timelines                 []Timeline       `bson:"timelines,omitempty" json:"timelines,omitempty"`
	Stages                    []Stage          `bson:"stages,omitempty" json:"stages,omitempty"`
	Channels                  []Channel        `bson:"channels,omitempty" json:"channels,omitempty"`
	IncidentChannel           *IncidentChannel `bson:"incidentChannel,omitempty" json:"incidentChannel,omitempty"`
	ConferenceDetails         []Conference     `bson:"conferenceDetails,omitempty" json:"conferenceDetails,omitempty"`
	Attachments               []Attachment     `bson:"attachments,omitempty" json:"attachments,omitempty"`
	mongodb.AuditDetails      `bson:",inline" json:",inline"`
}

type ChangeType string

const (
	ChangeTypeSeverity            ChangeType = "severity"
	ChangeTypeStatus              ChangeType = "status"
	ChangeTypeComment             ChangeType = "comment"
	ChangeTypeSlackChannelCreated ChangeType = "slackChannelCreated"
	ChangeTypeIncidentCreated     ChangeType = "incidentCreated"
)

type Timeline struct {
	ID            string            `bson:"id" json:"id"`
	Type          ChangeType        `bson:"type" json:"type"`
	CreatedAt     int64             `bson:"createdAt" json:"createdAt"`
	UpdatedAt     *int64            `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	PreviousState *string           `bson:"previousState,omitempty" json:"previousState,omitempty"`
	CurrentState  *string           `bson:"currentState,omitempty" json:"currentState,omitempty"`
	Slack         *Slack            `bson:"slack,omitempty" json:"slack,omitempty"`
	User          utils.UserDetails `bson:"userDetails" json:"userDetails"`
	Message       *string           `bson:"message,omitempty" json:"message,omitempty"`
}

type AttachmentType string

const (
	Link AttachmentType = "link"
)

type Attachment struct {
	Type        AttachmentType `bson:"type" json:"type"`
	Description string         `bson:"description" json:"description"`
	URL         string         `bson:"url" json:"url"`
}

type ChannelSource string

const (
	SlackSource ChannelSource = "slack"
)

type Channel struct {
	ID     string        `bson:"id" json:"id"`
	TeamID string        `bson:"teamId" json:"teamId"`
	Name   string        `bson:"name" json:"name"`
	Source ChannelSource `bson:"source" json:"source"`
	URL    string        `bson:"url" json:"url"`
	Status ChannelStatus `bson:"status" json:"status"`
}

type ConferenceType string

const (
	Zoom ConferenceType = "zoom"
)

type Conference struct {
	ID   string         `bson:"conferenceId" json:"conferenceId"`
	Type ConferenceType `bson:"type" json:"type"`
	URL  string         `bson:"url" json:"url"`
}

type RoleType string

const (
	IncidentCommander  RoleType = "Incident Commander"
	CommunicationsLead RoleType = "Communications Lead"
)

type Role struct {
	Type RoleType          `bson:"roleType" json:"roleType"`
	User utils.UserDetails `bson:"userDetails" json:"userDetails"`
}

type Functionality struct {
	ID   string `bson:"functionalityId" json:"functionalityId"`
	Name string `bson:"functionalityName" json:"functionalityName"`
}

type Service struct {
	ID   string `bson:"serviceId" json:"serviceId"`
	Name string `bson:"serviceName" json:"serviceName"`
}

type Environment struct {
	ID   string `bson:"environmentId" json:"environmentId"`
	Name string `bson:"environmentName" json:"environmentName"`
}

type Stage struct {
	ID        string            `bson:"stageId" json:"stageId"`
	Type      Status            `bson:"type" json:"type"`
	Duration  int64             `bson:"duration" json:"duration"`
	CreatedAt int64             `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64            `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	User      utils.UserDetails `bson:"userDetails" json:"userDetails"`
}

type ChannelStatus string

const (
	Operational ChannelStatus = "operational"
)

type IncidentChannelType string

const (
	ChannelSlack IncidentChannelType = "slack"
)

type IncidentChannel struct {
	Type  IncidentChannelType `bson:"type" json:"type"`
	Slack *Slack              `bson:"slack" json:"slack"`
}

type Slack struct {
	SlackTeam          `bson:",inline" json:",inline"`
	ChannelID          string        `bson:"channelId" json:"channelId"`
	ChannelName        string        `bson:"channelName" json:"channelName"`
	ChannelReference   string        `bson:"channelReference" json:"channelReference"`
	ChannelDescription string        `bson:"channelDescription" json:"channelDescription"`
	ChannelStatus      ChannelStatus `bson:"channelStatus" json:"channelStatus"`
}

type SlackTeam struct {
	TeamID     string `bson:"teamId" json:"teamId"`
	TeamName   string `bson:"teamName" json:"teamName"`
	TeamDomain string `bson:"teamDomain" json:"teamDomain"`
}

type SlackIncidentType string

const (
	Closed SlackIncidentType = "CLOSED"
	Open   SlackIncidentType = "OPEN"
)
