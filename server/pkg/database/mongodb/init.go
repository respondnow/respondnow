package mongodb

import (
	"context"
	"fmt"

	"github.com/respondnow/respond/server/pkg/constant"
	"github.com/respondnow/respond/server/pkg/database/mongodb/index"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Enum for Database collections
const (
	IncidentCollection = iota
	UserCollection
)

var (
	MClient     MongoInterface = &MongoClient{}
	MgoClient   *mongo.Client
	Collections = map[int]string{
		IncidentCollection: constant.IncidentCollection,
		UserCollection:     constant.UserCollection,
	}
)

var MongoSession = MgoClient

// MongoInterface requires a MongoClient that implements the Initialize method to create the Mongo DB client
// and a initAllCollection method to initialize all DB Collections
type MongoInterface interface {
	Initialize(client *mongo.Client, db string) (*MongoClient, error)
	initAllCollection() error
}

// MongoClient structure contains all the Database collections and the instance of the Database
type MongoClient struct {
	Database           *mongo.Database
	IncidentCollection *mongo.Collection
	UsersCollection    *mongo.Collection
}

// Initialize initializes database connection
func (m *MongoClient) Initialize(client *mongo.Client, db string) (*MongoClient, error) {
	m.Database = client.Database(db)
	err := m.initAllCollection()

	return m, err
}

// initAllCollection initializes all the database collections
func (m *MongoClient) initAllCollection() error {
	ctx := context.TODO()

	for _, coll := range Collections {
		err := m.initCollection(ctx, coll)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MongoClient) initCollection(ctx context.Context, collectionName string) error {
	logrus.Infof("setting up `%s` collection", collectionName)
	var (
		indexes = &index.IndexList{}
	)

	cs, err := m.Database.ListCollectionSpecifications(ctx,
		bson.D{primitive.E{Key: constant.Name, Value: collectionName}})
	if err != nil {
		return err
	}
	if len(cs) == 0 {
		logrus.Infof("collection `%s` not available, creating...", collectionName)
		if err := m.Database.CreateCollection(ctx, collectionName, nil); err != nil {
			return err
		}
	} else {
		logrus.Infof("collection `%s` already available", collectionName)
	}
	coll := m.Database.Collection(collectionName)
	switch collectionName {
	case constant.IncidentCollection:
		m.IncidentCollection = coll
		indexes = index.GetIncidentIndexList()
	case constant.UserCollection:
		m.UsersCollection = coll
	default:
		return fmt.Errorf("unknown collection given to initialize: %s", collectionName)
	}

	sc, err := coll.Indexes().ListSpecifications(ctx)
	if err != nil {
		return err
	}

	indexavailable := make(map[string]string)
	for _, item := range sc {
		indexavailable[item.Name] = item.Name
	}
	for name := range indexes.GetInactiveIndexes() {
		if _, ok := indexavailable[name]; ok {
			_, err := coll.Indexes().DropOne(ctx, name)
			if err != nil {
				return err
			}
			logrus.Infof("index `%s` dropped", name)
			continue
		}
		logrus.Infof("index `%s` not available", name)
	}
	for name, indexModel := range indexes.GetActiveIndexes() {
		if _, ok := indexavailable[name]; !ok {
			res, err := coll.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				return err
			}
			logrus.Infof("index `%s` created", res)
			continue
		}
		logrus.Infof("index `%s` available", name)
	}

	return nil
}
