package auth

import (
	"context"

	"github.com/respondnow/respond/server/pkg/database/mongodb/auth"
)

type AuthService interface {
	Signup(ctx context.Context, input AddUserInput) (auth.User, error)
	Login(ctx context.Context, input LoginUserInput) (auth.User, error)
	ChangePassword(ctx context.Context, input ChangeUserPasswordInput) error
	CreateJWTToken(email, userID, name string) (string, error)
	UpdateLastLogin(ctx context.Context, input LoginUserInput) error
}

type authService struct {
	authOperator auth.AuthOperator
}

// NewAuthService creates and returns a new instance of AuthService
func NewAuthService(authOperator auth.AuthOperator) AuthService {
	return &authService{
		authOperator: authOperator,
	}
}
