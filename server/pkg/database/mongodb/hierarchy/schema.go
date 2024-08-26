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
	AccountIDField   = "accountIdentifier"
	AccountNameField = "name"
)

const (
	OrganizationIDField        = "orgIdentifier"
	OrganizationNameField      = "name"
	OrganizationAccountIDField = "accountIdentifier"
)

const (
	ProjectIDField        = "projectIdentifier"
	ProjectNameField      = "name"
	ProjectOrgIDField     = "orgIdentifier"
	ProjectAccountIDField = "accountIdentifier"
)

const (
	UserMappingUserIDField    = "userId"
	UserMappingAccountIDField = "accountIdentifier"
	UserMappingOrgIDField     = "orgIdentifier"
	UserMappingProjectIDField = "projectIdentifier"
	UserMappingIsDefaultField = "isDefault"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccountID string             `bson:"accountIdentifier" json:"accountIdentifier"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
	Removed   bool               `bson:"removed" json:"removed"`
}

type Organization struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrgID     string             `bson:"orgIdentifier" json:"orgIdentifier"`
	Name      string             `bson:"name" json:"name"`
	AccountID string             `bson:"accountIdentifier" json:"accountIdentifier"` // Link to Account
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
	Removed   bool               `bson:"removed" json:"removed"`
}

type Project struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"projectIdentifier" json:"projectIdentifier"`
	Name      string             `bson:"name" json:"name"`
	OrgID     string             `bson:"orgIdentifier" json:"orgIdentifier"`         // Link to Organization
	AccountID string             `bson:"accountIdentifier" json:"accountIdentifier"` // Link to Account
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy string             `bson:"updatedBy" json:"updatedBy"`
	Removed   bool               `bson:"removed" json:"removed"`
}

type UserMapping struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId" json:"userId"`
	AccountID string             `bson:"accountIdentifier" json:"accountIdentifier"`
	OrgID     string             `bson:"orgIdentifier,omitempty" json:"orgIdentifier,omitempty"`
	ProjectID string             `bson:"projectIdentifier,omitempty" json:"projectIdentifier,omitempty"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	Removed   bool               `bson:"removed" json:"removed"`
	IsDefault bool               `bson:"isDefault" json:"isDefault"`
}
