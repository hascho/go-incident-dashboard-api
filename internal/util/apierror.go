package util

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
	Internal   error // holds the underlying go error for internal logging/debugging
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error (%d): %s", e.StatusCode, e.Message)
}

func NewNotFoundError(message string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}
