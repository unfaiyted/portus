// router/router.go
package router

import (
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

	healthService := services.NewHealthService(db)

	// Register all routes
	RegisterConfigRoutes(v1, configService)
	RegisterHealthRoutes(v1, healthService)

	return r
}
