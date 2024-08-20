package auth

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
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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
	ChangeUserPassword bool   `json:"changeUserPassword"`
}

type ChangePasswordResponseDTO struct {
	utils.DefaultResponseDTO `json:",inline"`
	Data                     LoginResponse `json:"data"`
}
