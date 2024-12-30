package core

import (
	"runtime"
)

// ApiError struct to hold error details.
type AppError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Stack      string `json:"stack"`
}

// NewApiError creates a new instance of ApiError
func NewAppError(statusCode int, message string) *AppError {
	stack := captureStackTrace()
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
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
func (e *AppError) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": e.StatusCode,
		"message":    e.Message,
		"stack":      e.Stack,
	}
}
