package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type CurrentUser struct {
	AccountID string `json:"accountId,omitempty"`
	AuthToken string `json:"authToken,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	jwt.StandardClaims
}

func ValidateJWT(token string, secret string) (*CurrentUser, error) {
	tkn, err := jwt.ParseWithClaims(token, &CurrentUser{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := tkn.Claims.(*CurrentUser)
	if ok {
		if claims.Name == "" || claims.Type == "" {
			return nil, fmt.Errorf("missing principal id or type")
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
