package incident

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/respondnow/respond/server/config"
	"github.com/respondnow/respond/server/pkg/api"
	"github.com/respondnow/respond/server/pkg/constant"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"github.com/respondnow/respond/server/pkg/database/mongodb/incident"
	"github.com/respondnow/respond/server/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type IncidentService interface {
	Get(ctx context.Context, id string) (incident.Incident, error)
	List(ctx context.Context, token string, listFilters ListFilters, correlationID, search string,
		limit, page int64, all bool) (ListResponse, error)
	Create(ctx context.Context, request CreateRequest,
		currentUser utils.UserDetails, correlationID string) (CreateResponse, error)
	AddConferenceDetailsForIncident(conferenceType incident.ConferenceType) (incident.Conference, error)
	GenerateIncidentIdentifier(createdAt int64) string
}

type incidentService struct {
	incidentOperator  incident.IncidentOperator
	accountIdentifier string
	orgIdentifier     string
	projectIdentifier string
}

func NewIncidentService(
	incidentOperator incident.IncidentOperator, accountIdentifier, orgIdentifier,
	projectIdentifier string) IncidentService {
	return &incidentService{
		incidentOperator:  incidentOperator,
		accountIdentifier: accountIdentifier,
		orgIdentifier:     orgIdentifier,
		projectIdentifier: projectIdentifier,
	}
}

func (is incidentService) GenerateIncidentIdentifier(createdAt int64) string {

	return strconv.Itoa(int(createdAt)) + "-" + uuid.New().String()
}

func (is incidentService) Create(ctx context.Context, request CreateRequest,
	currentUser utils.UserDetails, correlationID string) (CreateResponse, error) {
	resp := CreateResponse{
		Incident:      incident.Incident{},
		CorrelationID: correlationID,
	}

	err := is.ValidateCreateRequest(request)
	if err != nil {
		return resp, err
	}

	ts := time.Now().Unix()
	// Set default values
	if request.Status == "" {
		request.Status = incident.Started
	}
	if request.Identifier == "" {
		request.Identifier = is.GenerateIncidentIdentifier(ts)
	}
	channelCreatedDetails := incident.Slack{}
	if len(request.Channels) > 0 {
		for _, channel := range request.Channels {
			channelCreatedDetails = incident.Slack{
				ChannelID:     channel.ID,
				ChannelName:   channel.Name,
				ChannelStatus: channel.Status,
				SlackTeam: incident.SlackTeam{
					TeamID: channel.TeamID,
				},
			}
		}
	}

	newIncident := incident.Incident{
		ResourceDetails: request.ResourceDetails,
		IdentifierDetails: mongodb.IdentifierDetails{
			AccountIdentifier: is.accountIdentifier,
			OrgIdentifier:     is.orgIdentifier,
			ProjectIdentifier: is.projectIdentifier,
		},
		Type:            request.Type,
		Severity:        request.Severity,
		Status:          incident.Started,
		Summary:         request.Summary,
		Active:          true,
		IncidentChannel: request.IncidentChannel,
		Channels:        request.Channels,
		Services:        request.Services,
		Functionalities: request.Functionalities,
		Environments:    request.Environments,
		Attachments:     request.Attachments,
		Timelines: []incident.Timeline{
			{
				ID:        strconv.Itoa(int(time.Now().Unix())),
				Type:      incident.ChangeTypeIncidentCreated,
				CreatedAt: ts,
				UpdatedAt: &ts,
				User:      currentUser,
				Slack:     request.IncidentChannel.Slack,
			},
			{
				ID:        strconv.Itoa(int(time.Now().Unix())),
				Type:      incident.ChangeTypeSlackChannelCreated,
				CreatedAt: ts,
				UpdatedAt: &ts,
				User:      currentUser,
				Slack:     &channelCreatedDetails,
			},
		},
		Stages: []incident.Stage{
			{
				ID:        strconv.Itoa(int(time.Now().Unix())),
				Type:      request.Status,
				CreatedAt: ts,
				UpdatedAt: &ts,
				User:      currentUser,
			},
		},
		Roles: request.Roles,
		AuditDetails: mongodb.AuditDetails{
			CreatedBy: currentUser,
			CreatedAt: ts,
			UpdatedBy: currentUser,
			UpdatedAt: &ts,
		},
	}

	// add conference details
	if request.AddConference != nil {
		confDetails, err := is.AddConferenceDetailsForIncident(request.AddConference.Type)
		if err != nil {
			return resp, err
		}
		newIncident.ConferenceDetails = append(newIncident.ConferenceDetails, confDetails)
	}

	incident, err := is.incidentOperator.Create(ctx, newIncident)
	if err != nil {
		return resp, err
	}
	resp.Incident = incident

	return resp, nil
}

func (is incidentService) ValidateCreateRequest(request CreateRequest) error {
	if request.IncidentChannel == nil {
		return fmt.Errorf("incident channel must not be nil")
	} else if request.IncidentChannel.Type != "" {
		if request.IncidentChannel.Type == incident.ChannelSlack {
			if request.IncidentChannel.Slack.ChannelID == "" {
				return fmt.Errorf("incident slack channel id must not be nil")
			}
		}
	}

	return nil
}

func (is incidentService) AddConferenceDetailsForIncident(conferenceType incident.ConferenceType) (incident.Conference, error) {
	switch conferenceType {
	// TODO: Add zoom integration for generating new zoom links
	case incident.Zoom:
		return incident.Conference{
			ID:   strconv.Itoa(int(time.Now().Unix())),
			Type: conferenceType,
			URL:  config.EnvConfig.Conferences.ZoomLink,
		}, nil
	default:
		return incident.Conference{},
			fmt.Errorf("unsupported conference type provided: %s, supported: %s", conferenceType, incident.Zoom)
	}
}

func (is incidentService) List(ctx context.Context, token string, listFilters ListFilters, correlationID,
	search string, limit, page int64, all bool) (ListResponse, error) {
	resp := ListResponse{}
	filter := bson.D{}
	filter = append(filter, bson.E{
		Key:   constant.AccountIdentifier,
		Value: is.accountIdentifier,
	})
	if is.orgIdentifier != "" {
		filter = append(filter, bson.E{
			Key:   constant.OrgIdentifier,
			Value: is.orgIdentifier,
		})
	}
	if is.projectIdentifier != "" {
		filter = append(filter, bson.E{
			Key:   constant.ProjectIdentifier,
			Value: is.projectIdentifier,
		})
	}

	filter = append(filter, bson.E{
		Key:   "removed",
		Value: false,
	})
	if len(listFilters.Active) != 0 {
		isActive, err := strconv.ParseBool(listFilters.Active)
		if err != nil {
			return resp, err
		}
		filter = append(filter, bson.E{
			Key:   constant.Active,
			Value: isActive,
		})
	}
	if len(listFilters.Severity) > 0 {
		filter = append(filter, bson.E{
			Key:   constant.Severity,
			Value: listFilters.Severity,
		})
	}
	if len(listFilters.Type) > 0 {
		filter = append(filter, bson.E{
			Key:   constant.Type,
			Value: listFilters.Type,
		})
	}
	if len(listFilters.Status) > 0 {
		filter = append(filter, bson.E{
			Key:   constant.Status,
			Value: listFilters.Status,
		})
	}
	if len(listFilters.IncidentChannelType) > 0 {
		filter = append(filter, bson.E{
			Key:   constant.IncidentChannelType,
			Value: listFilters.IncidentChannelType,
		})
	}
	if search != "" {
		filter = append(filter, utils.NewUtils().GenerateSearchFilter(constant.Name, search, "i"))
	}

	incidentList, err := is.incidentOperator.CustomList(ctx, filter, limit, page*limit)
	if err != nil {
		return resp, err
	}
	count, err := is.incidentOperator.CountDocuments(ctx, filter)
	if err != nil {
		return resp, err
	}

	return ListResponse{
		Content:       incidentList,
		Pagination:    api.GetPagination(page, limit, count, all),
		CorrelationID: correlationID,
	}, nil
}

func (is incidentService) Get(ctx context.Context, id string) (incident.Incident, error) {
	filter := bson.D{}
	filter = append(filter, bson.E{
		Key:   constant.AccountIdentifier,
		Value: is.accountIdentifier,
	})
	if is.orgIdentifier != "" {
		filter = append(filter, bson.E{
			Key:   constant.OrgIdentifier,
			Value: is.orgIdentifier,
		})
	}
	if is.projectIdentifier != "" {
		filter = append(filter, bson.E{
			Key:   constant.ProjectIdentifier,
			Value: is.projectIdentifier,
		})
	}

	filter = append(filter, bson.E{
		Key:   "identifier",
		Value: id,
	})
	filter = append(filter, bson.E{
		Key:   "removed",
		Value: false,
	})

	resp, err := is.incidentOperator.Get(ctx, filter)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
