package main

import (
	"disaster-response-map-api/pkg/database"
	"disaster-response-map-api/pkg/router"
	"log"
)

func main() {
	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
	}

	// Ensure the database connection is closed when the app exits
	defer db.Close()

	// Initialize router and pass database instance
	r := router.SetupRouter(db)

	// Start the server
	log.Println("Server running on port 7000")
	r.Run(":7000")
}
