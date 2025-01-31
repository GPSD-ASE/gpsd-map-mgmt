package main

import (
	"log"
	"disaster-response-map-api/config"
    "disaster-response-map-api/pkg/database"
    "disaster-response-map-api/pkg/router"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Connect to PostgreSQL
	database.ConnectDB()

	// Start API Server
	r := router.SetupRouter()
	log.Println("Starting Disaster Response API on port 8080")
	r.Run(":8080")
}
