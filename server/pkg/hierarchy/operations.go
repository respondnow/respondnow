package hierarchy

import (
	"context"
	"errors"
	"time"

	"github.com/respondnow/respond/server/pkg/database/mongodb/hierarchy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *hierarchyManager) CreateAccount(ctx context.Context, account hierarchy.Account) error {
	filter := bson.M{hierarchy.AccountIDField: account.AccountID, hierarchy.FieldRemoved: false}
	_, err := h.operator.GetAccountByQuery(ctx, filter)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	} else {
		return errors.New("account with the given account_id already exists")
	}

	_, err = h.operator.AddAccount(ctx, account)
	return err
}

func (h *hierarchyManager) DeleteAccount(ctx context.Context, id string) error {
	filter := bson.M{hierarchy.AccountIDField: id, hierarchy.FieldRemoved: false}
	update := bson.M{
		"$set": bson.M{
			hierarchy.FieldRemoved:   true,
			hierarchy.FieldUpdatedAt: time.Now().Unix(),
		},
	}

	result, err := h.operator.UpdateAccount(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("account not found or already removed")
	}
	return nil
}

func (h *hierarchyManager) ReadAccount(ctx context.Context, id string) (hierarchy.Account, error) {
	filter := bson.M{hierarchy.AccountIDField: id, hierarchy.FieldRemoved: false}
	var account hierarchy.Account
	res, err := h.operator.GetAccountByQuery(ctx, filter)
	if err != nil {
		return account, err
	}

	return res, err
}

func (h *hierarchyManager) GetAllAccounts(ctx context.Context) ([]hierarchy.Account, error) {
	filter := bson.M{hierarchy.FieldRemoved: false}
	accounts, err := h.operator.GetAllAccountsByQuery(ctx, filter)
	if err != nil {
		return nil, err
	}

	return accounts, err
}

func (h *hierarchyManager) CreateOrganization(ctx context.Context, org hierarchy.Organization) error {
	_, err := h.operator.AddOrganization(ctx, org)
	return err
}

func (h *hierarchyManager) DeleteOrganization(ctx context.Context, id string) error {
	filter := bson.M{hierarchy.OrganizationIDField: id, hierarchy.FieldRemoved: false}
	update := bson.M{
		"$set": bson.M{
			hierarchy.FieldRemoved:   true,
			hierarchy.FieldUpdatedAt: time.Now().Unix(),
		},
	}
	result, err := h.operator.UpdateOrganization(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("organization not found or already removed")
	}
	return nil
}

func (h *hierarchyManager) ReadOrganization(ctx context.Context, id string) (hierarchy.Organization, error) {
	filter := bson.M{hierarchy.OrganizationIDField: id, hierarchy.FieldRemoved: false}
	res, err := h.operator.GetOrganizationByQuery(ctx, filter)
	if err != nil {
		return hierarchy.Organization{}, err
	}
	return res, err
}

func (h *hierarchyManager) GetAllOrganizations(ctx context.Context) ([]hierarchy.Organization, error) {
	filter := bson.M{hierarchy.FieldRemoved: false}
	organizations, err := h.operator.GetAllOrganizationsByQuery(ctx, filter)
	if err != nil {
		return nil, err
	}

	return organizations, err
}

func (h *hierarchyManager) CreateProject(ctx context.Context, proj hierarchy.Project) error {
	_, err := h.operator.AddProject(ctx, proj)
	return err
}

func (h *hierarchyManager) DeleteProject(ctx context.Context, id string) error {
	filter := bson.M{hierarchy.ProjectIDField: id, hierarchy.FieldRemoved: false}
	update := bson.M{
		"$set": bson.M{
			hierarchy.FieldRemoved:   true,
			hierarchy.FieldUpdatedAt: time.Now().Unix(),
		},
	}
	result, err := h.operator.UpdateProject(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("project not found or already removed")
	}
	return nil
}

func (h *hierarchyManager) ReadProject(ctx context.Context, id string) (hierarchy.Project, error) {
	filter := bson.M{hierarchy.ProjectIDField: id, hierarchy.FieldRemoved: false}
	res, err := h.operator.GetProjectByQuery(ctx, filter)
	if err != nil {
		return hierarchy.Project{}, err
	}

	return res, err
}

func (h *hierarchyManager) GetAllProjects(ctx context.Context) ([]hierarchy.Project, error) {
	filter := bson.M{hierarchy.FieldRemoved: false}
	projects, err := h.operator.GetAllProjectsByQuery(ctx, filter)
	if err != nil {
		return nil, err
	}

	return projects, err
}

func (h *hierarchyManager) CreateUserMapping(ctx context.Context, userID, accountID, orgID, projectID string, isDefault bool) (primitive.ObjectID, error) {
	mapping := hierarchy.UserMapping{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		AccountID: accountID,
		IsDefault: isDefault,
		CreatedAt: time.Now().Unix(),
	}

	if orgID != "" {
		mapping.OrgID = orgID
	}
	if projectID != "" {
		mapping.ProjectID = projectID
	}

	_, err := h.operator.AddUserMapping(ctx, mapping)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return mapping.ID, nil
}
