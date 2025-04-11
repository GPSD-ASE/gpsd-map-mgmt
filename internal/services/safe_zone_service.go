package services

import (
	"database/sql"
	"disaster-response-map-api/internal/models"
	"fmt"
)

type SafeZoneServiceInterface interface {
	CreateSafeZone(safeZone models.SafeZoneCreate) (int, error)
}

type SafeZoneService struct {
	DB *sql.DB
}

func NewSafeZoneService(db *sql.DB) *SafeZoneService {
	return &SafeZoneService{DB: db}
}

func (s *SafeZoneService) CreateSafeZone(safeZone models.SafeZoneCreate) (int, error) {
	var newID int
	query := `
        INSERT INTO safe_zone (zone_name, zone_lat, zone_lon, incident_type_id)
        VALUES ($1, $2, $3, $4)
        RETURNING zone_id
    `
	err := s.DB.QueryRow(query, safeZone.ZoneName, safeZone.ZoneLat, safeZone.ZoneLon, safeZone.IncidentTypeID).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new safe zone: %w", err)
	}
	return newID, nil
}
