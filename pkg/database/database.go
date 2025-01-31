package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	// Load Database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}

	log.Println("Connected to PostgreSQL successfully")

	// Run SQL migration script
	if err := runMigrations(); err != nil {
		log.Fatal("Database migration failed:", err)
	}
}

// Run SQL migrations
func runMigrations() error {
	script, err := os.ReadFile("tables.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = DB.Exec(string(script))
	if err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
