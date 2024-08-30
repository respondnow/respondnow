package user

import (
	"context"

	"github.com/respondnow/respondnow/server/pkg/database/mongodb/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	Signup(ctx context.Context, input AddUserInput) (user.User, error)
	Login(ctx context.Context, input LoginUserInput) (user.User, error)
	ChangePassword(ctx context.Context, input ChangeUserPasswordInput) error
	CreateJWTToken(email, userID, name string) (string, error)
	UpdateLastLogin(ctx context.Context, input LoginUserInput) error
	UpdateUser(ctx context.Context, query bson.M, updates bson.M) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}

type authService struct {
	authOperator user.AuthOperator
}

// NewAuthService creates and returns a new instance of AuthService
func NewAuthService(authOperator user.AuthOperator) AuthService {
	return &authService{
		authOperator: authOperator,
	}
}
