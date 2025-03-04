// models/responses.go
package models

type ErrorType string

const (
	ErrorTypeFailedCheck   ErrorType = "FAILED_CHECK"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeBadRequest    ErrorType = "BAD_REQUEST"
	ErrorTypeInternalError ErrorType = "INTERNAL_ERROR"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   ErrorType              `json:"error" example:"FAILED_CHECK"`
	Message string                 `json:"message" example:"This is a pretty message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}
