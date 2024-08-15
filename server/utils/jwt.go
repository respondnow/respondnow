package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type harnessClaims struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Claims struct {
	Type string `json:"type"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	if c.Type != "SERVICE" {
		return fmt.Errorf("token has wrong type: %s", c.Type)
	}
	return c.StandardClaims.Valid()
}

func (u *utils) GenerateJWTToken(identity string, subject string, signingKey []byte) (string, error) {
	var (
		tokenTypeService = "SERVICE"
		tokenIssuer      = "Harness Inc"
	)

	// Valid from an hour ago
	issuedTime := jwt.NewNumericDate(time.Now().Add(-time.Hour))

	// Expires in an hour from now
	expiryTime := jwt.NewNumericDate(time.Now().Add(time.Hour))

	harnessClaims := harnessClaims{
		Type: tokenTypeService,
		Name: "Chaos",
	}

	harnessClaims.Issuer = tokenIssuer
	harnessClaims.IssuedAt = issuedTime
	harnessClaims.NotBefore = issuedTime
	harnessClaims.ExpiresAt = expiryTime
	if subject != "" {
		harnessClaims.Subject = subject
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, harnessClaims)
	signedJwt, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", identity, signedJwt), nil
}
