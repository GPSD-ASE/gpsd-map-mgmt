package router

import (
	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/services"
	"disaster-response-map-api/pkg/database"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and routes.
// It now accepts both a database and a GraphHopper service.
func SetupRouter(db *database.Database, ghService *services.GraphHopperService) *gin.Engine {
	r := gin.Default()
	dzService := services.NewDisasterZoneService(db.DB)
	// Create disaster zone handler (using db)
	disasterZoneHandler := handlers.NewDisasterZoneHandler(db)
	r.GET("/zones", disasterZoneHandler.GetDisasterZones)

	// Routing handler
	routingHandler := handlers.NewRoutingHandler(ghService, dzService)
	r.GET("/routing", routingHandler.GetSafeRouting)

	// Evacuation endpoint (POST)
	evacService := services.NewEvacuationService(db.DB, ghService) // assuming db.DB is *sql.DB
	evacuationHandler := handlers.NewEvacuationHandler(evacService)
	r.POST("/evacuation", evacuationHandler.GetEvacuationRoute)

	return r
}
