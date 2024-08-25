package mongodb

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type GetCollectionInterface interface {
	getCollection(collectionType int) (*mongo.Collection, error)
}

type GetCollectionStruct struct{}

var GetCollectionClient GetCollectionInterface = &GetCollectionStruct{}

// getCollection function returns the appropriate DB collection based on the collection value passed
func (g *GetCollectionStruct) getCollection(collectionType int) (*mongo.Collection, error) {
	mongoClient := MClient
	switch collectionType {
	case IncidentCollection:
		return mongoClient.(*MongoClient).IncidentCollection, nil
	case UsersCollection:
		return mongoClient.(*MongoClient).UsersCollection, nil
	case AccountsCollection:
		return mongoClient.(*MongoClient).AccountsCollection, nil
	case OrganizationsCollection:
		return mongoClient.(*MongoClient).OrganizationsCollection, nil
	case ProjectsCollection:
		return mongoClient.(*MongoClient).ProjectsCollection, nil
	case UserMappingsCollection:
		return mongoClient.(*MongoClient).UserMappingsCollection, nil
	default:
		return nil, errors.New("unknown collection name")
	}
}
