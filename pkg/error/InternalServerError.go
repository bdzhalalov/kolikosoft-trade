package error

import "net/http"

type InternalServerError struct {
	Message string
	Code    int
}

func (e *InternalServerError) New() *BaseError {
	return &BaseError{
		Message: "Internal Server Error",
		Code:    http.StatusInternalServerError,
	}
}
