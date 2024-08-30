package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/respondnow/respondnow/server/config"
	"github.com/respondnow/respondnow/server/pkg/auth"
	"github.com/respondnow/respondnow/server/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationToken := c.GetHeader(config.Authorization)
		if authorizationToken == "" {
			handleError(c, http.StatusUnauthorized, "missing auth token")
			return
		}

		tokenParts, err := extractToken(authorizationToken)
		if err != nil {
			handleError(c, http.StatusUnauthorized, err.Error())
			return
		}

		cl, err := validateToken(tokenParts[0], tokenParts[1])
		if err != nil {
			handleError(c, http.StatusInternalServerError, "error processing jwt token: "+err.Error())
			return
		}

		c.Set(config.AuthToken, tokenParts[1])
		c.Set(config.AccountUUID, cl.Name)
		c.Set(config.PrincipalType, cl.Type)

		logrus.Debugf("jwt token is valid for userId: %s", cl.Username)
		c.Next()
	}
}

func extractToken(authorizationToken string) ([]string, error) {
	tokenParts := strings.Split(authorizationToken, " ")
	if len(tokenParts) != 2 {
		return nil, fmt.Errorf("invalid auth token")
	}
	return tokenParts, nil
}

func validateToken(prefix, token string) (*auth.CurrentUser, error) {
	if prefix != config.BearerType {
		return nil, fmt.Errorf("invalid token prefix")
	}
	return auth.ValidateJWT(token, config.EnvConfig.Auth.JWTSecret)
}

func handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, utils.DefaultResponseDTO{
		Status:  http.StatusText(statusCode),
		Message: message,
	})
	c.Abort()
}
