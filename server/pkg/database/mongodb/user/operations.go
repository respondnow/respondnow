package user

import (
	"context"

	"github.com/respondnow/respond/server/pkg/constant"
	"github.com/respondnow/respond/server/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (a *authOperator) AddUser(ctx context.Context, user User, opts ...*options.InsertOneOptions) (User, error) {
	res, err := a.operator.Create(ctx, mongodb.UsersCollection, user, opts...)
	if err != nil {
		return User{}, err
	}

	query := bson.M{constant.ID: res.InsertedID, "removed": false}
	return a.GetUserByQuery(ctx, query)
}

func (a *authOperator) GetUserByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (User, error) {
	out := User{}

	// Execute the query
	res, err := a.operator.Get(ctx, mongodb.UsersCollection, query, opts...)
	if err != nil {
		return out, err
	}

	// Decode the result into the User struct
	if err := res.Decode(&out); err != nil {
		return User{}, err
	}

	return out, nil
}

func (a *authOperator) UpdateUser(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updates}

	result, err := a.operator.Update(ctx, mongodb.UsersCollection, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authOperator) DeleteUser(ctx context.Context, filter bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := a.operator.Delete(ctx, mongodb.UsersCollection, filter, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
