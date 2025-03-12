package services

import (
	"context"
	"errors"
	"fmt"
	"portus/models"
	"portus/repository"
	"portus/utils"
	"time"

	"github.com/rs/zerolog/log"
)

// ShortenService provides methods to interact with URL shortening
type ShortenService interface {
	GetOriginalURL(ctx context.Context, code string) (string, error)
	Create(ctx context.Context, req models.ShortenRequest) (*models.ShortenResponse, error)
	Update(ctx context.Context, code string, req models.ShortenRequest) (*models.ShortenResponse, error)
	Delete(ctx context.Context, code string) error
	GetById(ctx context.Context, id uint64) *models.Shorten
	ShortCodeExists(ctx context.Context, randomCode string) (bool, error)
}

type shortenService struct {
	repo    repository.ShortenRepository
	baseURL string
}

// NewShortenService creates a new shortening service
func NewShortenService(repo repository.ShortenRepository, baseURL string) ShortenService {
	return &shortenService{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (s *shortenService) GetById(ctx context.Context, id uint64) *models.Shorten {
	shorten, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil
	}
	return shorten
}

func (s *shortenService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	shorten, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return "", err
	}

	if shorten == nil {
		return "", errors.New("short URL not found")
	}

	// Check if expired
	if !shorten.ExpiresAt.IsZero() && shorten.ExpiresAt.Before(time.Now()) {
		return "", errors.New("shortened URL has expired")
	}

	// Update click count asynchronously
	go s.repo.IncrementClickCount(ctx, code)

	return shorten.OriginalURL, nil
}

func (s *shortenService) Create(ctx context.Context, req models.ShortenRequest) (*models.ShortenResponse, error) {
	var shortCode string
	var err error

	log.Debug().Str("custom_code", req.CustomCode).Msg("Code passed")

	if req.CustomCode != "" {
		shortCode = req.CustomCode
		// Check if code already exists
		existing, _ := s.repo.FindByCode(ctx, shortCode)
		if existing != nil {
			return nil, errors.New("custom code already in use")
		}
	} else {
		// Generate random code
		shortCode = utils.GenerateShortCode()
	}

	// Set expiration time if specified
	var expiresAt time.Time
	if req.ExpiresAfter > 0 {
		expiresAt = time.Now().AddDate(0, 0, req.ExpiresAfter)
	}

	shorten := &models.Shorten{
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ClickCount:  0,
		ExpiresAt:   expiresAt,
	}

	newShorten, err := s.repo.Create(ctx, shorten)
	if err != nil {
		return nil, err
	}

	return &models.ShortenResponse{
		ShortCode:   newShorten.ShortCode,
		OriginalURL: newShorten.OriginalURL,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, shortCode),
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *shortenService) Update(ctx context.Context, code string, req models.ShortenRequest) (*models.ShortenResponse, error) {
	shorten, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if shorten == nil {
		return nil, errors.New("short URL not found")
	}

	shorten.OriginalURL = req.OriginalURL
	shorten.UpdatedAt = time.Now()

	// Update expiration if provided
	if req.ExpiresAfter > 0 {
		shorten.ExpiresAt = time.Now().AddDate(0, 0, req.ExpiresAfter)
	}

	updatedShorten, err := s.repo.Update(ctx, shorten)
	if err != nil {
		return nil, err
	}

	return &models.ShortenResponse{
		ShortCode:   updatedShorten.ShortCode,
		OriginalURL: req.OriginalURL,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, code),
		ExpiresAt:   shorten.ExpiresAt,
	}, nil
}

func (s *shortenService) Delete(ctx context.Context, code string) error {
	log := utils.LoggerFromContext(ctx)
	shorten, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if shorten == nil {
		return errors.New("short URL not found")
	}

	id, err := s.repo.Delete(ctx, code)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Error deleting")
		return err
	}
	return nil
}

func (s *shortenService) ShortCodeExists(ctx context.Context, code string) (bool, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().Str("code", code).Msg("Checking if short code exists")

	_, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		log.Error().Err(err).Str("code", code).Msg("Error finding short code")
		return false, nil
	}

	return true, nil
}
