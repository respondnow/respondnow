package incident

import (
	"context"

	"github.com/respondnow/respond/server/pkg/database/mongodb"
)

type IncidentOperator interface {
	List(ctx context.Context, page, limit int64, search string) ([]Incident, int64, error)
}

// Operator is the struct for incident operator
type incidentOperator struct {
	operator mongodb.MongoOperator
}

// NewIncidentOperator returns a new instance of incident operator
func NewIncidentOperator(mongodbOperator mongodb.MongoOperator) IncidentOperator {
	return &incidentOperator{
		operator: mongodbOperator,
	}
}
