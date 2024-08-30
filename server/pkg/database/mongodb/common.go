package mongodb

import (
	"github.com/respondnow/respondnow/server/utils"
)

type ResourceDetails struct {
	Name        string   `bson:"name" json:"name" binding:"required"`
	Identifier  string   `bson:"identifier" json:"identifier" binding:"required"`
	Description string   `bson:"description" json:"description"`
	Tags        []string `bson:"tags" json:"tags"`
}

type IdentifierDetails struct {
	AccountIdentifier string `bson:"accountIdentifier" json:"accountIdentifier"`
	OrgIdentifier     string `bson:"orgIdentifier" json:"orgIdentifier"`
	ProjectIdentifier string `bson:"projectIdentifier" json:"projectIdentifier"`
}

type AuditDetails struct {
	CreatedAt int64             `bson:"createdAt" json:"createdAt"`
	UpdatedAt *int64            `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy utils.UserDetails `bson:"createdBy" json:"createdBy"`
	UpdatedBy utils.UserDetails `bson:"updatedBy" json:"updatedBy"`
	RemovedAt *int64            `bson:"removedAt,omitempty" json:"removedAt,omitempty"`
	Removed   bool              `bson:"removed" json:"removed"`
}
