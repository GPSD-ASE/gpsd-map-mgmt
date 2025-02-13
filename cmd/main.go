package main

import (
	"log"

	"disaster-response-map-api/config"
	"disaster-response-map-api/internal/services"
	"disaster-response-map-api/pkg/database"
	"disaster-response-map-api/pkg/router"
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

	// Create GraphHopper service using the API key from config
	ghService := services.NewGraphHopperService(config.GRAPHHOPPER_KEY)

	// Initialize router with both database and GraphHopper service
	r := router.SetupRouter(db, ghService)

	log.Println("Server running on port 7000")
	r.Run(":7001")
}
