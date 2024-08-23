package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
	Active                 bool               `bson:"active" json:"active" binding:"required"`
	Name                   string             `bson:"name" json:"name"`
	UserID                 string             `bson:"userId" json:"userId"`
	Email                  string             `bson:"email" json:"email"`
	Password               string             `bson:"password" json:"password"`
	ChangePasswordRequired bool               `bson:"changePasswordRequired" json:"changePasswordRequired"`
	CreatedAt              int64              `bson:"createdAt" json:"createdAt"`
	UpdatedAt              *int64             `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy              string             `bson:"createdBy" json:"createdBy"`
	UpdatedBy              string             `bson:"updatedBy" json:"updatedBy"`
	RemovedAt              *int64             `bson:"removedAt,omitempty" json:"removedAt,omitempty"`
	Removed                bool               `bson:"removed" json:"removed"`
	LastLoginAt            int64              `bson:"lastLoginAt" json:"lastLoginAt"`
}

const (
	UserID string = "userId"
	Email  string = "email"
)
