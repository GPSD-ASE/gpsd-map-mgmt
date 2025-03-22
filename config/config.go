package config

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

var (
	MAP_MGMT_ENV      string
	MAP_MGMT_DB_HOST  string
	MAP_MGMT_DB_NAME  string
	MAP_MGMT_DB_PORT  string
	MAP_MGMT_DB_PASS  string
	MAP_MGMT_DB_USER  string
	MAP_MGMT_APP_PORT string
	JWT_SECRET        string
	GRAPHHOPPER_KEY   string
	GRAPHHOPPER_URL   string
	TOMTOM_API_KEY    string
	TOMTOM_URL        string
)

func LoadConfig() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Warning: No .env file found, using system env variables")
	// }
	if os.Getenv("VAULT_TOKEN") != "" {
		vaultSecrets, err := loadVaultSecrets()
		if err != nil {
			log.Printf("Error loading secrets from Vault: %v", err)
		} else {
			MAP_MGMT_DB_HOST = getString(vaultSecrets, "MAP_MGMT_DB_HOST", os.Getenv("MAP_MGMT_DB_HOST"))
			MAP_MGMT_DB_NAME = getString(vaultSecrets, "MAP_MGMT_DB_NAME", os.Getenv("MAP_MGMT_DB_NAME"))
			MAP_MGMT_DB_PORT = getString(vaultSecrets, "MAP_MGMT_DB_PORT", os.Getenv("MAP_MGMT_DB_PORT"))
			MAP_MGMT_DB_PASS = getString(vaultSecrets, "MAP_MGMT_DB_PASS", os.Getenv("MAP_MGMT_DB_PASS"))
			MAP_MGMT_DB_USER = getString(vaultSecrets, "MAP_MGMT_DB_USER", os.Getenv("MAP_MGMT_DB_USER"))
			MAP_MGMT_APP_PORT = getString(vaultSecrets, "MAP_MGMT_APP_PORT", os.Getenv("MAP_MGMT_APP_PORT"))
			JWT_SECRET = getString(vaultSecrets, "JWT_SECRET", os.Getenv("JWT_SECRET"))
			GRAPHHOPPER_KEY = getString(vaultSecrets, "GRAPHHOPPER_KEY", os.Getenv("GRAPHHOPPER_KEY"))
			TOMTOM_API_KEY = getString(vaultSecrets, "TOMTOM_API_KEY", os.Getenv("TOMTOM_API_KEY"))
			GRAPHHOPPER_URL = getString(vaultSecrets, "GRAPHHOPPER_URL", os.Getenv("GRAPHHOPPER_URL"))
			TOMTOM_URL = getString(vaultSecrets, "TOMTOM_URL", os.Getenv("TOMTOM_URL"))
		}
		log.Printf("DEBUG - All vault secrets : %v", vaultSecrets)
	}
	if MAP_MGMT_DB_HOST == "" {
		MAP_MGMT_DB_HOST = os.Getenv("MAP_MGMT_DB_HOST")
	}
	if MAP_MGMT_DB_PASS == "" {
		MAP_MGMT_DB_PASS = os.Getenv("MAP_MGMT_DB_PASS")
	}
	if MAP_MGMT_DB_NAME == "" {
		MAP_MGMT_DB_NAME = os.Getenv("MAP_MGMT_DB_NAME")
	}
	if MAP_MGMT_DB_PORT == "" {
		MAP_MGMT_DB_PORT = os.Getenv("MAP_MGMT_DB_PORT")
	}
	if MAP_MGMT_DB_USER == "" {
		MAP_MGMT_DB_USER = os.Getenv("MAP_MGMT_DB_USER")
	}
	if JWT_SECRET == "" {
		JWT_SECRET = os.Getenv("JWT_SECRET")
	}
	if GRAPHHOPPER_KEY == "" {
		GRAPHHOPPER_KEY = os.Getenv("GRAPHHOPPER_KEY")
	}

	if TOMTOM_API_KEY == "" {
		TOMTOM_API_KEY = os.Getenv("TOMTOM_API_KEY")
	}
	if GRAPHHOPPER_URL == "" {
		GRAPHHOPPER_URL = os.Getenv("GRAPHHOPPER_URL")
	}
	if TOMTOM_URL == "" {
		TOMTOM_URL = os.Getenv("TOMTOM_URL")
	}
	if MAP_MGMT_DB_HOST == "" || MAP_MGMT_DB_PASS == "" || MAP_MGMT_DB_NAME == "" || MAP_MGMT_DB_PORT == "" || MAP_MGMT_DB_USER == "" || JWT_SECRET == "" || GRAPHHOPPER_KEY == "" || TOMTOM_API_KEY == "" || GRAPHHOPPER_URL == "" || TOMTOM_URL == "" {
		log.Fatal("Missing environment variables")
	}
}

func loadVaultSecrets() (map[string]interface{}, error) {
	config := vault.DefaultConfig()
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		config.Address = addr
	}
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("missing vault_token")
	}
	client.SetToken(token)

	secret, err := client.KVv2("secret").Get(context.Background(), "gpsd/map-mgmt")
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from vault: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found in vault")
	}
	return secret.Data, nil
}

func getString(secrets map[string]interface{}, key, def string) string {
	if val, ok := secrets[key].(string); ok && val != "" {
		return val
	}
	return def
}
