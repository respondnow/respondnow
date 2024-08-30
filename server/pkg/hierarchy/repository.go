package hierarchy

import (
	"context"

	"github.com/respondnow/respondnow/server/pkg/database/mongodb/hierarchy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HierarchyManager interface {
	CreateAccount(ctx context.Context, account hierarchy.Account) error
	DeleteAccount(ctx context.Context, id string) error
	ReadAccount(ctx context.Context, id string) (hierarchy.Account, error)
	GetAllAccounts(ctx context.Context) ([]hierarchy.Account, error)

	CreateOrganization(ctx context.Context, org hierarchy.Organization) error
	DeleteOrganization(ctx context.Context, id string) error
	ReadOrganization(ctx context.Context, id string) (hierarchy.Organization, error)
	GetAllOrganizations(ctx context.Context) ([]hierarchy.Organization, error)

	CreateProject(ctx context.Context, proj hierarchy.Project) error
	DeleteProject(ctx context.Context, id string) error
	ReadProject(ctx context.Context, id string) (hierarchy.Project, error)
	GetAllProjects(ctx context.Context) ([]hierarchy.Project, error)
	CreateUserMapping(ctx context.Context, userID, accountID, orgID, projectID string, isDefault bool) (primitive.ObjectID, error)
}

type hierarchyManager struct {
	operator hierarchy.HierarchyOperator
}

func NewHierarchyManager(mongoOperator hierarchy.HierarchyOperator) HierarchyManager {
	return &hierarchyManager{
		operator: mongoOperator,
	}
}
