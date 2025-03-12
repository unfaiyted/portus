// models/responses.go
package models

import "time"

type ErrorType string

const (
	ErrorTypeFailedCheck         ErrorType = "FAILED_CHECK"
	ErrorTypeUnauthorized        ErrorType = "UNAUTHORIZED"
	ErrorTypeNotFound            ErrorType = "NOT_FOUND"
	ErrorTypeBadRequest          ErrorType = "BAD_REQUEST"
	ErrorTypeInternalError       ErrorType = "INTERNAL_ERROR"
	ErrorTypeForbidden           ErrorType = "FORBIDDEN"
	ErrorTypeConflict            ErrorType = "CONFLICT"
	ErrorTypeValidation          ErrorType = "VALIDATION_ERROR"
	ErrorTypeRateLimited         ErrorType = "RATE_LIMITED"
	ErrorTypeTimeout             ErrorType = "TIMEOUT"
	ErrorTypeServiceUnavailable  ErrorType = "SERVICE_UNAVAILABLE"
	ErrorTypeUnprocessableEntity ErrorType = "UNPROCESSABLE_ENTITY"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     ErrorType              `json:"error" example:"FAILED_CHECK"`
	Message   string                 `json:"message" example:"This is a pretty message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
}

// APIResponse represents a generic API response
type APIResponse[T any] struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Operation successful"`
	Data    T      `json:"data,omitempty"`
}

// Type-specific response creators
func NewShortenResponse(shorten *Shorten, message string) APIResponse[ShortenData] {
	return APIResponse[ShortenData]{
		Success: true,
		Message: message,
		Data:    ShortenData{Shorten: shorten},
	}
}

// func NewPasteListResponse(pastes []Paste, count int) APIResponse[PasteListData] {
// 	return APIResponse[PasteListData]{
// 		Success: true,
// 		Message: "Pastes retrieved successfully",
// 		Data:    PasteListData{Pastes: pastes, Count: count},
// 	}
// }
