package router

import (
	"portus/handlers"
	"portus/services"

	"github.com/gin-gonic/gin"
)

func RegisterShortenRoutes(rg *gin.RouterGroup, service services.ShortenService) {
	shortenHandlers := handlers.NewShortenHandler(service)
	shorts := rg.Group("/shorten")
	{

		shorts.POST("", shortenHandlers.Create)
		shorts.PUT("/:code", shortenHandlers.Update)
		shorts.DELETE("/:code", shortenHandlers.Delete)
		shorts.GET("/:code", shortenHandlers.Redirect)

	}
}
