package index

import (
	"github.com/respondnow/respond/server/pkg/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// don't remove any item from the list mark Available as false if you don't want the index.
// don't change Name of the index also, Name should be unique.
func GetIncidentIndexList() *IndexList {
	return &IndexList{
		Items: []Index{
			{
				Name:      "unique_incidentIdentifier_1",
				Available: true,
				Model: mongo.IndexModel{
					Keys: bson.D{
						{
							Key:   "identifier",
							Value: constant.One,
						},
					},
					Options: options.Index().
						SetName("unique_incidentIdentifier_1").
						SetUnique(true).
						SetPartialFilterExpression(bson.D{
							{
								Key:   "removed",
								Value: false,
							},
						}),
				},
			},
		},
	}
}
