package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/respondnow/respond/server/pkg/database/mongodb/user"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt/v4"
	"github.com/respondnow/respond/server/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (a *authService) Signup(ctx context.Context, input AddUserInput) (user.User, error) {
	query := bson.M{user.Email: input.Email, "removed": false}
	_, err := a.authOperator.GetUserByQuery(ctx, query)
	if err == nil {
		// Email already exists, return an error indicating duplicate user
		return user.User{}, fmt.Errorf("email: %v already exists", input.Email)
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), config.EnvConfig.Auth.PasswordHashCost)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	currentTime := time.Now().Unix()

	u := user.User{
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

	createdUser, err := a.authOperator.AddUser(ctx, u)
	if err != nil {
		return user.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil

}

func (a *authService) Login(ctx context.Context, input LoginUserInput) (user.User, error) {
	query := bson.M{user.Email: input.Email, "removed": false}
	u, err := a.authOperator.GetUserByQuery(ctx, query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, fmt.Errorf("email %v not found", input.Email)
		}
		return user.User{}, fmt.Errorf("error fetching user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return user.User{}, errors.New("invalid user ID or password")
		}
		return user.User{}, fmt.Errorf("error comparing passwords: %w", err)
	}

	return u, nil
}

func (a *authService) UpdateLastLogin(ctx context.Context, input LoginUserInput) error {
	query := bson.M{user.Email: input.Email, "removed": false}
	updates := bson.M{"lastLoginAt": time.Now().Unix(), "updatedAt": time.Now().Unix()}

	_, err := a.authOperator.UpdateUser(ctx, query, updates)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) UpdateUser(ctx context.Context, query bson.M, updates bson.M) error {
	_, err := a.authOperator.UpdateUser(ctx, query, updates)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	query := bson.M{
		"_id": id,
	}
	_, err := a.authOperator.DeleteUser(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) ChangePassword(ctx context.Context, input ChangeUserPasswordInput) error {
	query := bson.M{user.Email: input.Email, "removed": false}
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
		Type:     UserPrincipalType,
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
