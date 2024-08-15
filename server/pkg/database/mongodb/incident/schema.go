package incident

import (
	"time"

	"github.com/respondnow/respond/server/pkg/api"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"github.com/respondnow/respond/server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IncidentList struct {
	Data       []Incident     `bson:"data" json:"data"`
	Pagination api.Pagination `json:"page"`
}

type Status string

const (
	Started             Status = "started"
	Detected            Status = "detected"
	Acknowledged        Status = "acknowledged"
	Identified          Status = "identified"
	Mitigated           Status = "mitigated"
	Investigating       Status = "investigating"
	PostmortemStarted   Status = "postmortemStarted"
	PostmortemCompleted Status = "PostmortemCompleted"
	Resolved            Status = "resolved"
)

type Severity string

const (
	Severity1 Severity = "SEV1"
	Severity2 Severity = "SEV2"
	Severity3 Severity = "SEV3"
	Severity4 Severity = "SEV4"
)

type Incident struct {
	ID                      primitive.ObjectID `bson:"_id" json:"id"`
	mongodb.ResourceDetails `bson:",inline" json:",inline"`
	Severity                Severity        `bson:"severity" json:"severity" binding:"required"`
	Status                  Status          `bson:"status" json:"status" binding:"required"`
	Summary                 string          `bson:"summary" json:"summary" binding:"required"`
	Active                  bool            `bson:"active" json:"active" binding:"required"`
	Services                []Service       `bson:"services" json:"services"`
	Environments            []Environment   `bson:"environments" json:"environments"`
	Functionalities         []Functionality `bson:"functionalities" json:"functionalities"`
	Roles                   []Role          `bson:"roles" json:"roles"`
	Stages                  []Stage         `bson:"stages" json:"stages"`
	Channels                []Channel       `bson:"channel" json:"channel"`
	Slack                   *Slack          `bson:"slack,omitempty" json:"slack,omitempty"`
	ConferenceDetails       []Conference    `bson:"conferenceDetails" json:"conferenceDetails"`
	mongodb.AuditDetails    `bson:",inline" json:",inline"`
}

type ChannelSource string

const (
	SlackSource ChannelSource = "slack"
)

type Channel struct {
	ID     string        `bson:"id" json:"id"`
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
	IncidentCommander     RoleType = "incidentCommander"
	CommunicationsLiaison RoleType = "communicationsLiaison"
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
	ID        string     `bson:"stageId" json:"stageId"`
	Type      Status     `bson:"type" json:"type"`
	Duration  int64      `bson:"duration" json:"duration"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt *time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type ChannelStatus string

const (
	Operational ChannelStatus = "operational"
)

type Slack struct {
	SlackChannel `bson:",inline" json:",inline"`
}

type SlackChannel struct {
	ChannelID          string        `bson:"channelId" json:"channelId"`
	ChannelName        string        `bson:"channelName" json:"channelName"`
	ChannelReference   string        `bson:"channelReference" json:"channelReference"`
	ChannelDescription string        `bson:"channelDescription" json:"channelDescription"`
	ChannelStatus      ChannelStatus `bson:"channelStatus" json:"channelStatus"`
}
