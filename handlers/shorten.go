package handlers

import (
	"net/http"
	"portus/models"
	"portus/services"

	"github.com/gin-gonic/gin"
)

type ShortenHandler struct {
	service services.ShortenService
}

func NewShortenHandler(service services.ShortenService) *ShortenHandler {
	return &ShortenHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create a shortened URL
// @Description Creates a new shortened URL
// @Tags shorten
// @Accept json
// @Produce json
// @Param request body models.ShortenRequest true "URL to shorten"
// @Success 201 {object} models.ShortenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /shorten [post]
func (h *ShortenHandler) Create(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorTypeBadRequest,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: models.ErrorTypeInternalError,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Update godoc
// @Summary Update a shortened URL
// @Description Updates an existing shortened URL
// @Tags shorten
// @Accept json
// @Produce json
// @Param code path string true "Short code"
// @Param request body models.ShortenRequest true "Updated URL data"
// @Success 200 {object} models.ShortenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /shorten/{code} [put]
func (h *ShortenHandler) Update(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "short code is required"})
		return
	}

	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: models.ErrorTypeBadRequest,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	result, err := h.service.Update(c.Request.Context(), code, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "short URL not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.ErrorResponse{
			Error: models.ErrorTypeInternalError,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete godoc
// @Summary Delete a shortened URL
// @Description Deletes an existing shortened URL
// @Tags shorten
// @Param code path string true "Short code"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /shorten/{code} [delete]
func (h *ShortenHandler) Delete(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "short code is required"})
		return
	}

	err := h.service.Delete(c.Request.Context(), code)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "short URL not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.ErrorResponse{
			Error: models.ErrorTypeNotFound,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects to the original URL from a short code
// @Tags shorten
// @Param code path string true "Short code"
// @Success 302 "Redirect to original URL"
// @Failure 404 {object} models.ErrorResponse
// @Router /{code} [get]
func (h *ShortenHandler) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "short code is required"})
		return
	}

	url, err := h.service.GetOriginalURL(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: models.ErrorTypeNotFound,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})

		return
	}

	c.Redirect(http.StatusFound, url)
}
