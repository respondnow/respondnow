package hierarchy

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	FieldID        = "_id"
	FieldCreatedAt = "createdAt"
	FieldUpdatedAt = "updatedAt"
	FieldCreatedBy = "createdBy"
	FieldUpdatedBy = "updatedBy"
	FieldRemoved   = "removed"
)

const (
	AccountIDField   = "accountId"
	AccountNameField = "name"
)

const (
	OrganizationIDField        = "orgId"
	OrganizationNameField      = "name"
	OrganizationAccountIDField = "accountId"
)

const (
	ProjectIDField        = "projectId"
	ProjectNameField      = "name"
	ProjectOrgIDField     = "orgId"
	ProjectAccountIDField = "accountId"
)

const (
	UserMappingUserIDField    = "userId"
	UserMappingAccountIDField = "accountId"
	UserMappingOrgIDField     = "orgId"
	UserMappingProjectIDField = "projectId"
	UserMappingIsDefaultField = "isDefault"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID string             `bson:"accountId" json:"accountId"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
	Removed   bool               `bson:"removed" json:"removed"`
}

type Organization struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrgID     string             `bson:"orgId" json:"orgId"`
	Name      string             `bson:"name" json:"name"`
	AccountID string             `bson:"accountId" json:"accountId"` // Link to Account
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
	Removed   bool               `bson:"removed" json:"removed"`
}

type Project struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"projectId" json:"projectId"`
	Name      string             `bson:"name" json:"name"`
	OrgID     string             `bson:"orgId" json:"orgId"`         // Link to Organization
	AccountID string             `bson:"accountId" json:"accountId"` // Link to Account
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
	Removed   bool               `bson:"removed" json:"removed"`
}

type UserMapping struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId" json:"userId"`
	AccountID string             `bson:"accountId" json:"accountId"`
	OrgID     string             `bson:"orgId,omitempty" json:"orgId,omitempty"`
	ProjectID string             `bson:"projectId,omitempty" json:"projectId,omitempty"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	Removed   bool               `bson:"removed" json:"removed"`
	IsDefault bool               `bson:"isDefault" json:"isDefault"`
}
