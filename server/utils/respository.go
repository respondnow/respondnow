package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Utils interface {
	RestCall(method string, url string, payload []byte, opt ...func(option *RestCallOption)) (int, []byte, error)
	GenerateJWTToken(identity string, subject string, signingKey []byte) (string, error)
	GenerateSearchFilter(field string, search string, options string) primitive.E
	GetPagination(ctx *gin.Context) (int64, int64, bool)
	GetCorrelationID(ctx *gin.Context) string
	RandStringBytes(n int) string
	StrToBool(s string) bool
}

type utils struct{}

func NewUtils() Utils {
	return &utils{}
}

type MockUtils struct {
	mock.Mock
}
