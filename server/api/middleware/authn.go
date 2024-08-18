package middleware

import "github.com/golang-jwt/jwt"

type CurrentUser struct {
	AccountID string `json:"accountId,omitempty"`
	AuthToken string `json:"authToken,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	jwt.StandardClaims
}
