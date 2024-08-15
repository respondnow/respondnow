package utils

import (
	"github.com/stretchr/testify/mock"
)

type Utils interface {
	RestCall(method string, url string, payload []byte, opt ...func(option *RestCallOption)) (int, []byte, error)
	GenerateJWTToken(identity string, subject string, signingKey []byte) (string, error)
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
