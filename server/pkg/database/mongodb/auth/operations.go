package auth

import (
	"context"

	"github.com/respondnow/respond/server/pkg/constant"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *authOperator) AddUser(ctx context.Context, user User, opts ...*options.InsertOneOptions) (User, error) {
	res, err := a.operator.Create(ctx, mongodb.UserCollection, user, opts...)
	if err != nil {
		return User{}, err
	}

	query := bson.M{constant.ID: res.InsertedID, "removed": false}
	return a.GetUserByQuery(ctx, query)
}

func (a *authOperator) GetUserByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (User, error) {
	out := User{}

	// Execute the query
	res, err := a.operator.Get(ctx, mongodb.UserCollection, query, opts...)
	if err != nil {
		return out, err
	}

	// Decode the result into the User struct
	if err := res.Decode(&out); err != nil {
		return User{}, err
	}

	return out, nil
}

func (a *authOperator) UpdateUserByField(ctx context.Context, field string, value interface{}, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: field, Value: value}}

	update := bson.M{"$set": updates}

	result, err := a.operator.Update(ctx, mongodb.UserCollection, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
