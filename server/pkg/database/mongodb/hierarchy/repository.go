package hierarchy

import (
	"context"

	"github.com/respondnow/respond/server/pkg/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HierarchyOperator interface {
	AddAccount(ctx context.Context, account Account, opts ...*options.InsertOneOptions) (Account, error)
	GetAccountByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (Account, error)
	UpdateAccount(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	GetAllAccountsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]Account, error)

	AddOrganization(ctx context.Context, org Organization, opts ...*options.InsertOneOptions) (Organization, error)
	GetOrganizationByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (Organization, error)
	UpdateOrganization(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	GetAllOrganizationsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]Organization, error)

	AddProject(ctx context.Context, project Project, opts ...*options.InsertOneOptions) (Project, error)
	GetProjectByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (Project, error)
	UpdateProject(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	GetAllProjectsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]Project, error)

	AddUserMapping(ctx context.Context, mapping UserMapping, opts ...*options.InsertOneOptions) (UserMapping, error)
	GetUserMappingByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (UserMapping, error)
	UpdateUserMapping(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	GetAllUserMappingsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]UserMapping, error)
}

type hierarchyOperator struct {
	operator mongodb.MongoOperator
}

func NewHierarchyOperator(mongodbOperator mongodb.MongoOperator) HierarchyOperator {
	return &hierarchyOperator{
		operator: mongodbOperator,
	}
}
