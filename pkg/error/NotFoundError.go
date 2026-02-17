package error

import "net/http"

type NotFoundError struct {
	Message string
	Code    int
}

func (e *NotFoundError) New(message string) *BaseError {
	return &BaseError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}
