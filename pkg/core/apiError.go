package core

import (
	"runtime"
)

// ApiError struct to hold error details.
type ApiError struct {
	StatusCode int      `json:"statusCode"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors"`
	Stack      string   `json:"stack"`
}

// NewApiError creates a new instance of ApiError
func NewApiError(statusCode int, message string, errors []string) *ApiError {
	stack := captureStackTrace()
	return &ApiError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
		Stack:      stack,
	}
}

// captureStackTrace captures the stack trace
func captureStackTrace() string {
	stack := make([]byte, 1024)
	n := runtime.Stack(stack, false)
	return string(stack[:n])
}

// ToMap converts ApiError to a map for easier JSON encoding
func (e *ApiError) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": e.StatusCode,
		"message":    e.Message,
		"errors":     e.Errors,
		"stack":      e.Stack,
	}
}
