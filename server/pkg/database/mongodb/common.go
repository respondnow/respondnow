package mongodb

import (
	"time"

	"github.com/respondnow/respond/server/utils"
)

type ResourceDetails struct {
	Name        string   `bson:"name" json:"name" binding:"required"`
	Identifier  string   `bson:"identifier" json:"identifier" binding:"required"`
	Description string   `bson:"description" json:"description"`
	Tags        []string `bson:"tags" json:"tags"`
}

type AuditDetails struct {
	CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt *time.Time        `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy utils.UserDetails `bson:"createdBy" json:"createdBy"`
	UpdatedBy utils.UserDetails `bson:"updatedBy" json:"updatedBy"`
	RemovedAt *time.Time        `bson:"removedAt,omitempty" json:"removedAt,omitempty"`
	Removed   bool              `bson:"removed" json:"removed"`
}
