package mongodb

import (
	"context"

	"github.com/respondnow/respond/server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -destination=mocks/mock_mongodb.go -package=mock_mongodb github.com/respondnow/respond/server/pkg/database/mongodb MongoOperator
type MongoOperator interface {
	Create(ctx context.Context, collectionType int, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	CreateMany(ctx context.Context, collectionType int, documents []interface{},
		opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	Get(ctx context.Context, collectionType int, filter interface{},
		opts ...*options.FindOneOptions) (*mongo.SingleResult, error)
	List(ctx context.Context, collectionType int, filter interface{},
		opts ...*options.FindOptions) (*mongo.Cursor, error)
	Update(ctx context.Context, collectionType int, filter, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, collectionType int, filter, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateByID(ctx context.Context, collectionType int, id, update interface{},
		opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Replace(ctx context.Context, collectionType int, filter,
		replacement interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, collectionType int, filter interface{},
		opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, collectionType int, filter interface{},
		opts ...*options.CountOptions) (int64, error)
	Aggregate(ctx context.Context, collectionType int, pipeline interface{},
		opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	BulkWrite(ctx context.Context, collectionType int, document []mongo.WriteModel,
		opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	GetCollection(collectionType int) (*mongo.Collection, error)
	ListCollection(ctx context.Context, dbName string, mclient *mongo.Client) ([]string, error)
	ListDataBase(ctx context.Context, mclient *mongo.Client) ([]string, error)
}

type MongoOperations struct{}

var (
	// Operator contains all the CRUD operations of the mongo database
	Operator MongoOperator = &MongoOperations{}
)

// Create puts a document in the database
func (m *MongoOperations) Create(ctx context.Context, collectionType int, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return nil, err
	}
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateMany puts an array of documents in the database
func (m *MongoOperations) CreateMany(ctx context.Context, collectionType int, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return nil, err
	}
	res, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// BulkWrite creates a bulk write operation in the database
func (m *MongoOperations) BulkWrite(ctx context.Context, collectionType int, document []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return nil, err
	}
	result, err := collection.BulkWrite(ctx, document, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get fetches a document from the database based on a query
func (m *MongoOperations) Get(ctx context.Context, collectionType int, filter interface{},
	opts ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return nil, err
	}

	maxTime := options.FindOne().SetMaxTime(config.EnvConfig.MaxQueryExecutionTime)
	opts = append(opts, maxTime)

	result := collection.FindOne(ctx, filter, opts...)
	return result, nil
}

// List fetches a list of documents from the database based on a query
func (m *MongoOperations) List(ctx context.Context, collectionType int, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return nil, err
	}

	f := false
	for _, opt := range opts {
		if opt.MaxTime != nil {
			f = true
			break
		}
	}

	if !f {
		maxTime := options.Find().SetMaxTime(config.EnvConfig.MaxQueryExecutionTime)
		opts = append(opts, maxTime)
	}

	result, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates a document in the database based on a query
func (m *MongoOperations) Update(ctx context.Context, collectionType int, filter, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return result, err
	}
	result, err = collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateMany updates multiple documents in the database based on a query
func (m *MongoOperations) UpdateMany(ctx context.Context, collectionType int, filter, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return result, err
	}

	result, err = collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UpdateByID updates a document in the database based on document ID
func (m *MongoOperations) UpdateByID(ctx context.Context, collectionType int, id, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return result, err
	}
	result, err = collection.UpdateByID(ctx, id, update, opts...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Replace changes a document with a new one in the database based on a query
func (m *MongoOperations) Replace(ctx context.Context, collectionType int, filter,
	replacement interface{}) (*mongo.UpdateResult, error) {
	var result *mongo.UpdateResult
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return result, err
	}
	// If the given item is not present then insert.
	opts := options.Replace().SetUpsert(true)
	result, err = collection.ReplaceOne(ctx, filter, replacement, opts)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Delete removes a document from the database based on a query
func (m *MongoOperations) Delete(ctx context.Context, collectionType int, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	var result *mongo.DeleteResult
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return result, err
	}
	result, err = collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// CountDocuments returns the number of documents in the collection that matches a query
func (m *MongoOperations) CountDocuments(ctx context.Context, collectionType int,
	filter interface{}, opts ...*options.CountOptions) (int64, error) {
	var result int64 = 0
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return result, err
	}

	f := false
	for _, opt := range opts {
		if opt.MaxTime != nil {
			f = true
			break
		}
	}

	if !f {
		maxTime := options.Count().SetMaxTime(config.EnvConfig.MaxQueryExecutionTime)
		opts = append(opts, maxTime)
	}

	result, err = collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (m *MongoOperations) Aggregate(ctx context.Context, collectionType int,
	pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	collection, err := m.GetCollection(collectionType)
	if err != nil {
		return nil, err
	}

	f := false
	for _, opt := range opts {
		if opt.MaxTime != nil {
			f = true
			break
		}
	}

	if !f {
		maxTime := options.Aggregate().SetMaxTime(config.EnvConfig.MaxQueryExecutionTime)
		opts = append(opts, maxTime)
	}

	result, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetCollection fetches the correct collection based on the collection type
func (m *MongoOperations) GetCollection(collectionType int) (*mongo.Collection, error) {
	return GetCollectionClient.getCollection(collectionType)
}

func (m *MongoOperations) ListDataBase(ctx context.Context, mclient *mongo.Client) ([]string, error) {
	dbs, err := mclient.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (m *MongoOperations) ListCollection(ctx context.Context, dbName string, mclient *mongo.Client) ([]string, error) {
	cols, err := mclient.Database(dbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	return cols, nil
}
