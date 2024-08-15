package incident

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/respondnow/respond/server/pkg/database/mongodb"
)

func (i *incidentOperator) List(ctx context.Context, page, limit int64,
	search string) ([]Incident, int64, error) {
	filter := bson.D{
		{
			Key:   "removed",
			Value: false,
		},
	}

	sort := bson.D{{
		Key:   "created_at",
		Value: -1,
	}}

	results, err := i.operator.List(ctx, mongodb.IncidentCollection, filter,
		options.MergeFindOptions().SetSort(sort).SetSkip(page*limit).SetLimit(limit))
	if err != nil {
		return nil, 0, err
	}

	var incidents []Incident
	err = results.All(ctx, &incidents)
	if err != nil {
		return nil, 0, err
	}
	count, err := i.operator.CountDocuments(ctx, mongodb.IncidentCollection, filter)
	if err != nil {
		return nil, 0, err
	}

	// Pass an empty array instead of nil
	if incidents == nil {
		incidents = make([]Incident, 0)
	}

	return incidents, count, nil
}
