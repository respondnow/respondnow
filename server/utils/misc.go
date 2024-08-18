package utils

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/respondnow/respond/server/pkg/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (u *utils) RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (u *utils) StrToBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "yes", "y", "t":
		return true
	default:
		return false
	}
}

func (u *utils) GenerateSearchFilter(field string, search string, options string) primitive.E {
	return primitive.E{
		Key: field,
		Value: primitive.Regex{
			Pattern: search,
			Options: options,
		},
	}
}

// GetCorrelationID returns correlation id from echo context. If correlation id is present in querry param return that
// if correlation id not present then generate and return a new correlation id
func (u *utils) GetCorrelationID(ctx *gin.Context) string {
	correlationID := ctx.Query(constant.CorrelationID)
	if correlationID == "" {
		correlationID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return correlationID
}

// GetPagination returns page index, per page limit and show all from echo context
func (u *utils) GetPagination(ctx *gin.Context) (int64, int64, bool) {
	page, err := strconv.ParseInt(ctx.Query(constant.Page), 10, 64)
	if err != nil {
		page = constant.DefaultPage
	}
	limit, err := strconv.ParseInt(ctx.Query(constant.Limit), 10, 64)
	if err != nil {
		limit = constant.DefaultLimit
	}
	if page <= 0 {
		page = constant.DefaultPage
	}
	if limit <= 0 {
		limit = constant.DefaultLimit
	}
	all, err := strconv.ParseBool(ctx.Query(constant.All))
	if err != nil {
		all = false
	}
	if all {
		limit = 0
	}
	return page, limit, all
}
