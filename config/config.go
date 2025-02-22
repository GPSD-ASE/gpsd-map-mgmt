package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_URL          string
	JWT_SECRET      string
	DB_USERNAME     string
	DB_PASS         string
	GRAPHHOPPER_KEY string
	PORT            string
	TOMTOM_API_KEY  string
	GRAPHHOPPER_URL string
	TOMTOM_URL      string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables")
	}

	DB_URL = os.Getenv("DATABASE_URL")
	DB_PASS = os.Getenv("DATABASE_PASS")
	DB_USERNAME = os.Getenv("DATABASE_USERNAME")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	GRAPHHOPPER_KEY = os.Getenv("GRAPHHOPPER_KEY")
	PORT = os.Getenv("PORT")
	TOMTOM_API_KEY = os.Getenv("TOMTOM_API_KEY")
	GRAPHHOPPER_URL = os.Getenv("GRAPHHOPPER_URL")
	TOMTOM_URL = os.Getenv("TOMTOM_URL")

	if DB_URL == "" || JWT_SECRET == "" {
		log.Fatal("Missing required environment variables")
	}
}
