// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @host localhost:7000
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package main

import (
	"log"
	"net/http"

	"disaster-response-map-api/config"
	"disaster-response-map-api/internal/services"
	"disaster-response-map-api/pkg/database"
	"disaster-response-map-api/pkg/router"

	// Swagger UI packages:
	_ "disaster-response-map-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load environment configuration including GRAPHHOPPER_KEY
	config.LoadConfig()

	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ghService := services.NewGraphHopperService(config.GRAPHHOPPER_KEY, config.GRAPHHOPPER_URL)
	tfService := services.NewTrafficService(config.TOMTOM_URL, config.TOMTOM_API_KEY)

	r := router.SetupRouter(db, ghService, tfService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server running on port ", config.PORT)
	http.ListenAndServe(":"+config.PORT, r)
}
