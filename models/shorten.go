package models

import "time"

type Shorten struct {
	ID          uint64    `json:"id" example:"1"`
	OriginalURL string    `json:"original_url" binding:"required" example:"https://example.com/some/long/path"`
	ShortCode   string    `json:"short_code" example:"abc123"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ClickCount  uint64    `json:"click_count" example:"0"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// ShortenRequest represents the request to create a shortened URL
type ShortenRequest struct {
	OriginalURL  string `json:"original_url" binding:"required"`
	CustomCode   string `json:"custom_code,omitempty"`
	ExpiresAfter int    `json:"expires_after,omitempty"` // In days
}

// ShortenResponse represents the response for a shortened URL
type ShortenResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// ErrorResponse represents error responses
// type ErrorResponse struct {
// 	Error string `json:"error"`
// }
