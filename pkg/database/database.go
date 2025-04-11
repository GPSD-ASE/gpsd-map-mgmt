package database

import (
	"database/sql"
	"disaster-response-map-api/config"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type DatabaseInterface interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

type Database struct {
	DB *sql.DB
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func NewDatabase() (*Database, error) {
	config.LoadConfig()

	dbURL := config.MAP_MGMT_DB_HOST
	println(dbURL)
	if !strings.Contains(dbURL, "@") {
		dbURL = fmt.Sprintf("postgres://%s:%s@%s", config.MAP_MGMT_DB_USER, config.MAP_MGMT_DB_PASS, dbURL)
	}
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
