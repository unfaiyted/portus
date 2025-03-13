// router/router.go
package router

import (
	"portus/repository"
	"portus/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, configService services.ConfigService) *gin.Engine {
	r := gin.Default()

	// CORS config
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	r.Use(cors.New(config))

	// Setup API v1 routes
	v1 := r.Group("/api/v1")

	// TODO: should I fix this? It doesent technically need a repo, but ti does interact with the database?
	healthService := services.NewHealthService(db)

	shortenRepo := repository.NewShortenRepository(db)
	shortenService := services.NewShortenService(shortenRepo, configService.GetConfig().App.AppURL)

	// Register all routes
	RegisterConfigRoutes(v1, configService)
	RegisterHealthRoutes(v1, healthService)
	RegisterShortenRoutes(v1, shortenService)

	return r
}
