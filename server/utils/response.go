package utils

type DefaultResponseDTO struct {
	Status        string `json:"status"`
	Message       string `json:"message,omitempty"`
	CorrelationId string `json:"correlationId,omitempty"`
}

type Status string

const (
	SUCCESS Status = "SUCCESS"
	FAILURE Status = "FAILURE"
	ERROR   Status = "ERROR"
)
