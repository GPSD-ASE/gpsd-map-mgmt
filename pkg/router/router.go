package router

import (
	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/services"
	"disaster-response-map-api/pkg/database"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and routes.
// It now accepts both a database and a GraphHopper service.
func SetupRouter(db *database.Database, ghService *services.GraphHopperService, tfService *services.TrafficService) *gin.Engine {
	r := gin.Default()
	dzService := services.NewDisasterZoneService(db.DB)
	// Create disaster zone handler (using db)
	disasterZoneHandler := handlers.NewDisasterZoneHandler(dzService)
	r.GET("/zones", disasterZoneHandler.GetDisasterZones)
	// Traffic handler (using tfService)
	trafficHandler := handlers.NewTrafficHandler(tfService)
	r.GET("/traffic", trafficHandler.GetTrafficData)
	// Routing handler
	routingHandler := handlers.NewRoutingHandler(ghService, dzService)
	r.GET("/routing", routingHandler.GetSafeRouting)

	r.GET("/route", routingHandler.GetDefaultRoute)
	// Evacuation endpoint (POST)
	evacService := services.NewEvacuationService(db.DB, ghService) // assuming db.DB is *sql.DB
	evacuationHandler := handlers.NewEvacuationHandler(evacService)
	r.POST("/evacuation", evacuationHandler.GetEvacuationRoute)

	safeZoneService := services.NewSafeZoneService(db.DB)
	safeZoneHandler := handlers.NewSafeZoneHandler(safeZoneService)
	r.POST("/safezones", safeZoneHandler.CreateSafeZone)
	return r
}
