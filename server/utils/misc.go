package utils

import (
	"math/rand"
	"strings"
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
