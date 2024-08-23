package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/respondnow/respond/server/config"
	"github.com/respondnow/respond/server/pkg/database/mongodb/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (a *authService) Signup(ctx context.Context, input AddUserInput) (auth.User, error) {
	query := bson.M{auth.Email: input.Email, "removed": false}
	_, err := a.authOperator.GetUserByQuery(ctx, query)
	if err == nil {
		// Email already exists, return an error indicating duplicate user
		return auth.User{}, fmt.Errorf("email: %v already exists", input.Email)
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return auth.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), config.EnvConfig.Auth.PasswordHashCost)
	if err != nil {
		return auth.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	currentTime := time.Now().Unix()

	user := auth.User{
		Active:                 false,
		Name:                   input.Name,
		UserID:                 input.UserID,
		Email:                  input.Email,
		Password:               string(hashedPassword),
		ChangePasswordRequired: true,
		CreatedAt:              currentTime,
		UpdatedAt:              &currentTime,
		CreatedBy:              "",
		UpdatedBy:              "",
		RemovedAt:              nil,
		Removed:                false,
		LastLoginAt:            0,
	}

	createdUser, err := a.authOperator.AddUser(ctx, user)
	if err != nil {
		return auth.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil

}

func (a *authService) Login(ctx context.Context, input LoginUserInput) (auth.User, error) {
	query := bson.M{auth.Email: input.Email, "removed": false}
	user, err := a.authOperator.GetUserByQuery(ctx, query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return auth.User{}, fmt.Errorf("email %v not found", auth.Email)
		}
		return auth.User{}, fmt.Errorf("error fetching user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return auth.User{}, errors.New("invalid user ID or password")
		}
		return auth.User{}, fmt.Errorf("error comparing passwords: %w", err)
	}

	return user, nil
}

func (a *authService) UpdateLastLogin(ctx context.Context, input LoginUserInput) error {
	query := bson.M{auth.Email: input.Email, "removed": false}
	updates := bson.M{"lastLoginAt": time.Now().Unix(), "updatedAt": time.Now().Unix()}

	_, err := a.authOperator.UpdateUser(ctx, query, updates)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) ChangePassword(ctx context.Context, input ChangeUserPasswordInput) error {
	query := bson.M{auth.Email: input.Email, "removed": false}
	user, err := a.authOperator.GetUserByQuery(ctx, query)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updates := bson.M{"password": string(hashedPassword), "updatedAt": time.Now().Unix(), "changePasswordRequired": false, "active": true}
	_, err = a.authOperator.UpdateUser(ctx, query, updates)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) CreateJWTToken(email, userID, name string) (string, error) {
	timeOut := time.Hour * 24
	expirationTime := time.Now().Add(timeOut)

	claims := &CustomClaims{
		Email:    email,
		Name:     name,
		Username: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    JWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.EnvConfig.Auth.JWTSecret))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", tokenString), nil
}
