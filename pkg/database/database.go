package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DatabaseInterface defines the expected database operations
type DatabaseInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

// Database struct for PostgreSQL implementation
type Database struct {
	DB *sql.DB
}

// Exec executes an SQL statement
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

// Query runs an SQL query and returns rows
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// NewDatabase initializes a new PostgreSQL database
func NewDatabase() (*Database, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return &Database{DB: db}, nil
}

// Migrate runs SQL migrations
func (d *Database) Migrate() error {
	script, err := os.ReadFile("tables.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = d.DB.Exec(string(script))
	if err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
