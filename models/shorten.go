package models

import "time"

type Shorten struct {
	ID          uint64    `json:"id" example:"1"`
	OriginalURL string    `json:"originalUrl" binding:"required" example:"https://example.com/some/long/path"`
	ShortCode   string    `json:"shortCode" example:"abc123"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ClickCount  uint64    `json:"clickCount" example:"0"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

// ShortenRequest represents the request to create a shortened URL
type ShortenRequest struct {
	OriginalURL  string `json:"originalUrl" binding:"required"`
	CustomCode   string `json:"customCode,omitempty"`
	ExpiresAfter int    `json:"expiresAfter,omitempty"` // In days
}

// ShortenResponse represents the response for a shortened URL
type ShortenResponse struct {
	ShortCode   string    `json:"shortCode"`
	OriginalURL string    `json:"originalUrl"`
	ShortURL    string    `json:"shortUrl"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type ShortenData struct {
	Shorten *Shorten `json:"shorten"`
}
