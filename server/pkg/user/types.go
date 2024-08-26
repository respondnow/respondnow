package user

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/respondnow/respond/server/utils"
)

type AddUserInput struct {
	Name     string `json:"name" validate:"required"`
	UserID   string `json:"userId" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangeUserPasswordInput struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type CustomClaims struct {
	Type     PrincipalType `json:"type"`
	Email    string        `json:"email"`
	Name     string        `json:"name"`
	Username string        `json:"username"`
	jwt.RegisteredClaims
}

type PrincipalType string

const (
	UserPrincipalType    PrincipalType = "USER"
	ServicePrincipalType PrincipalType = "SERVICE"
)

const (
	JWTIssuer = "respondNow"
)

type SignupResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     interface{} `json:"data"`
}

type LoginResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     LoginResponse `json:"data"`
}

type LoginResponse struct {
	Token              string `json:"token,omitempty"`
	LastLoginAt        int64  `json:"lastLoginAt"`
	ChangeUserPassword bool   `json:"changeUserPassword"`
}

type ChangePasswordResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     ChangePasswordResponse `json:"data"`
}

type ChangePasswordResponse struct {
	Token       string `json:"token,omitempty"`
	LastLoginAt int64  `json:"lastLoginAt"`
}

type GetUserMappingResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     UserMapping `json:"data"`
}

type UserMapping struct {
	DefaultMapping Identifiers   `json:"defaultMapping"`
	Mappings       []Identifiers `json:"mappings"`
}

type Identifiers struct {
	AccountID   string `json:"accountId"`
	AccountName string `json:"accountName"`
	OrgID       string `json:"orgId,omitempty"`
	OrgName     string `json:"orgName"`
	ProjectID   string `json:"projectId,omitempty"`
	ProjectName string `json:"projectName"`
}
