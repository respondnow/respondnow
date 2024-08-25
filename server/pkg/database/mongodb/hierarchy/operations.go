package hierarchy

import (
	"context"

	"github.com/respondnow/respond/server/pkg/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *hierarchyOperator) AddAccount(ctx context.Context, account Account, opts ...*options.InsertOneOptions) (Account, error) {
	res, err := h.operator.Create(ctx, mongodb.AccountsCollection, account, opts...)
	if err != nil {
		return Account{}, err
	}

	query := bson.M{FieldID: res.InsertedID, FieldRemoved: false}
	return h.GetAccountByQuery(ctx, query)
}

func (h *hierarchyOperator) GetAccountByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (Account, error) {
	var out Account

	res, err := h.operator.Get(ctx, mongodb.AccountsCollection, query, opts...)
	if err != nil {
		return out, err
	}

	if err := res.Decode(&out); err != nil {
		return Account{}, err
	}

	return out, nil
}

func (h *hierarchyOperator) UpdateAccount(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updates}

	result, err := h.operator.Update(ctx, mongodb.AccountsCollection, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *hierarchyOperator) GetAllAccountsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]Account, error) {
	var accounts []Account
	cursor, err := h.operator.List(ctx, mongodb.AccountsCollection, query, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &accounts)
	return accounts, err
}

func (h *hierarchyOperator) AddOrganization(ctx context.Context, org Organization, opts ...*options.InsertOneOptions) (Organization, error) {
	res, err := h.operator.Create(ctx, mongodb.OrganizationsCollection, org, opts...)
	if err != nil {
		return Organization{}, err
	}

	query := bson.M{FieldID: res.InsertedID, FieldRemoved: false}
	return h.GetOrganizationByQuery(ctx, query)
}

func (h *hierarchyOperator) GetOrganizationByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (Organization, error) {
	var out Organization

	res, err := h.operator.Get(ctx, mongodb.OrganizationsCollection, query, opts...)
	if err != nil {
		return out, err
	}

	if err := res.Decode(&out); err != nil {
		return Organization{}, err
	}

	return out, nil
}

func (h *hierarchyOperator) UpdateOrganization(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updates}

	result, err := h.operator.Update(ctx, mongodb.OrganizationsCollection, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *hierarchyOperator) GetAllOrganizationsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]Organization, error) {
	var organizations []Organization
	cursor, err := h.operator.List(ctx, mongodb.OrganizationsCollection, query, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &organizations)
	return organizations, err
}

func (h *hierarchyOperator) AddProject(ctx context.Context, project Project, opts ...*options.InsertOneOptions) (Project, error) {
	res, err := h.operator.Create(ctx, mongodb.ProjectsCollection, project, opts...)
	if err != nil {
		return Project{}, err
	}

	query := bson.M{FieldID: res.InsertedID, FieldRemoved: false}
	return h.GetProjectByQuery(ctx, query)
}

func (h *hierarchyOperator) GetProjectByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (Project, error) {
	var out Project

	res, err := h.operator.Get(ctx, mongodb.ProjectsCollection, query, opts...)
	if err != nil {
		return out, err
	}

	if err := res.Decode(&out); err != nil {
		return Project{}, err
	}

	return out, nil
}

func (h *hierarchyOperator) UpdateProject(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updates}

	result, err := h.operator.Update(ctx, mongodb.ProjectsCollection, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *hierarchyOperator) GetAllProjectsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]Project, error) {
	var projects []Project
	cursor, err := h.operator.List(ctx, mongodb.ProjectsCollection, query, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &projects)
	return projects, err
}

func (h *hierarchyOperator) AddUserMapping(ctx context.Context, mapping UserMapping, opts ...*options.InsertOneOptions) (UserMapping, error) {
	_, err := h.operator.Create(ctx, mongodb.UserMappingsCollection, mapping, opts...)
	if err != nil {
		return UserMapping{}, err
	}
	return mapping, nil
}

func (h *hierarchyOperator) GetUserMappingByQuery(ctx context.Context, query bson.M, opts ...*options.FindOneOptions) (UserMapping, error) {
	var mapping UserMapping
	res, err := h.operator.Get(ctx, mongodb.UserMappingsCollection, query, opts...)
	if err != nil {
		return mapping, err
	}
	err = res.Decode(&mapping)
	return mapping, err
}

func (h *hierarchyOperator) UpdateUserMapping(ctx context.Context, filter bson.M, updates bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update := bson.M{"$set": updates}
	result, err := h.operator.Update(ctx, mongodb.UserMappingsCollection, filter, update, opts...)
	return result, err
}

func (h *hierarchyOperator) GetAllUserMappingsByQuery(ctx context.Context, query bson.M, opts ...*options.FindOptions) ([]UserMapping, error) {
	var mappings []UserMapping
	cursor, err := h.operator.List(ctx, mongodb.UserMappingsCollection, query, opts...)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &mappings)
	return mappings, err
}
