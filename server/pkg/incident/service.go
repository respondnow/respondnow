package incident

import (
	"context"

	"github.com/respondnow/respond/server/pkg/api"
	"github.com/respondnow/respond/server/pkg/database/mongodb/incident"
)

type IncidentService interface {
	ListIncidents(context context.Context, page, limit int64, search string) (incident.IncidentList, error)
}

type incidentService struct {
	incidentOperator incident.IncidentOperator
}

func NewIncidentService(
	incidentOperator incident.IncidentOperator) IncidentService {
	return &incidentService{
		incidentOperator: incidentOperator,
	}
}

func (is incidentService) ListIncidents(context context.Context, page, limit int64,
	search string) (incident.IncidentList, error) {
	incidentList, count, err := is.incidentOperator.List(context, page, limit, search)
	if err != nil {
		return incident.IncidentList{}, err
	}
	p := api.Pagination{
		Index: page,
		Limit: limit,
		TotalPages: func() int64 {
			if count%limit == 0 {
				return count / limit
			}
			return (count / limit) + 1
		}(),
		TotalItems: count,
	}
	return incident.IncidentList{Data: incidentList, Pagination: p}, nil
}
