package user

import (
	"context"

	"github.com/respondnow/respondnow/server/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthOperator interface {
	AddUser(ctx context.Context, user User, opts ...*options.InsertOneOptions) (User, error)
	GetUserByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (User, error)
	UpdateUser(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteUser(ctx context.Context, filter bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

type authOperator struct {
	operator mongodb.MongoOperator
}

func NewAuthOperator(mongodbOperator mongodb.MongoOperator) AuthOperator {
	return &authOperator{
		operator: mongodbOperator,
	}
}
