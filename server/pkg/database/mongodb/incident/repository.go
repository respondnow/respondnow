package incident

import (
	"context"

	"github.com/respondnow/respondnow/server/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IncidentOperator interface {
	Create(ctx context.Context, in Incident, opts ...*options.InsertOneOptions) (Incident, error)
	GetByID(ctx context.Context, id interface{}, opts ...*options.FindOneOptions) (Incident, error)
	Get(ctx context.Context, id interface{}, opts ...*options.FindOneOptions) (Incident, error)
	List(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]Incident, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	UpdateByID(ctx context.Context, in Incident, opts ...*options.UpdateOptions) (Incident, error)
	CustomList(ctx context.Context, filter interface{}, limit, skip int64) ([]Incident, error)
	BulkProcessWithSessionContext(ctx mongo.SessionContext, createList, updateList []Incident,
		opts ...*options.BulkWriteOptions) error
	WithDefaults(in *Incident)
	Validate(in *Incident) error
	GetIncidentTypes() []Type
	GetIncidentSeverities() []Severity
	GetIncidentAttachmentType() []AttachmentType
	GetIncidentStageStatuses() []Status
	GetIncidentRoles() []RoleType
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
