package incident

import (
	"github.com/respondnow/respondnow/server/pkg/api"
	"github.com/respondnow/respondnow/server/pkg/database/mongodb"
	"github.com/respondnow/respondnow/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respondnow/server/utils"
)

type ListResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     ListResponse `json:"data"`
}

type ListResponse struct {
	Content       []incident.Incident `json:"content"`
	Pagination    api.Pagination      `json:"pagination"`
	CorrelationID string              `json:"correlationID"`
}

type GetResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     incident.Incident `json:"data"`
}

type GetFilters struct {
	mongodb.IdentifierDetails `json:",inline"`
	IncidentId                string `json:"incidentId"`
}

type ListFilters struct {
	Type                incident.Type
	Severity            incident.Severity
	IncidentChannelType incident.IncidentChannelType
	Status              incident.Status
	Active              string
}

type CreateRequest struct {
	mongodb.ResourceDetails `json:",inline"`
	Type                    incident.Type             `json:"type" binding:"required"`
	Severity                incident.Severity         `json:"severity" binding:"required"`
	Summary                 string                    `json:"summary" binding:"required"`
	IncidentChannel         *incident.IncidentChannel `json:"incidentChannel" binding:"required"`
	Status                  incident.Status           `json:"status"`
	Services                []incident.Service        `json:"services,omitempty"`
	Environments            []incident.Environment    `json:"environments,omitempty"`
	Functionalities         []incident.Functionality  `json:"functionalities,omitempty"`
	Channels                []incident.Channel        `json:"channels,omitempty"`
	Roles                   []incident.Role           `json:"roles,omitempty"`
	AddConference           *AddConference            `json:"addConference,omitempty"`
	Attachments             []incident.Attachment     `json:"attachments,omitempty"`
}

type AddConference struct {
	Type incident.ConferenceType
}

type CreateResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     CreateResponse `json:"data"`
}

type CreateResponse struct {
	incident.Incident `json:",inline" binding:"required"`
	CorrelationID     string `json:"correlationID"`
}
