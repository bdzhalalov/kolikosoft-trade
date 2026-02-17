package error

import "net/http"

type BadRequestError struct {
	Message string
	Code    int
}

func (e *BadRequestError) New(message string) *BaseError {
	return &BaseError{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}
