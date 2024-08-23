package utils

type UserDetails struct {
	Email    string `json:"email"`
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Source   Source `json:"source"`
}

type Source string

const (
	Slack Source = "slack"
)
