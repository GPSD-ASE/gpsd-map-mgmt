package router

import (
	"github.com/gin-gonic/gin"
	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/pkg/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public Endpoints
	r.GET("/health", handlers.HealthCheck)

	// Protected Endpoints
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())

	api.POST("/disaster-zones", handlers.CreateDisasterZone)
	api.GET("/disaster-zones", handlers.GetDisasterZones)
	api.DELETE("/disaster-zones/:id", handlers.DeleteDisasterZone)

	return r
}
