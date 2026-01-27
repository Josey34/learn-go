package domain

type ExternalErrors struct {
	Message       string
	Code          string
	Status        int
	CorrelationID string
}

func NewExternalError(message string, code string, status int, correlationId string) *ExternalErrors {
	return &ExternalErrors{
		Message:       message,
		Code:          code,
		Status:        status,
		CorrelationID: correlationId,
	}
}
