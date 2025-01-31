package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_URL     string
	JWT_SECRET string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables")
	}

	DB_URL = os.Getenv("DATABASE_URL")
	JWT_SECRET = os.Getenv("JWT_SECRET")

	if DB_URL == "" || JWT_SECRET == "" {
		log.Fatal("Missing required environment variables")
	}
}
