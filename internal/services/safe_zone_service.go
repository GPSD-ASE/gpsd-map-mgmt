package services

import (
	"database/sql"
	"disaster-response-map-api/internal/models"
	"fmt"
)

type SafeZoneServiceInterface interface {
	CreateSafeZone(safeZone models.SafeZoneCreate) (int, error)
	GetSafeZones() ([]models.SafeZone, error)
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

func (s *SafeZoneService) GetSafeZones() ([]models.SafeZone, error) {
	query := `SELECT zone_id, zone_name, zone_lat, zone_lon, incident_type_id FROM safe_zone`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query safe zones: %v", err)
	}
	defer rows.Close()

	var safeZones []models.SafeZone
	for rows.Next() {
		var sz models.SafeZone
		if err := rows.Scan(&sz.ZoneID, &sz.ZoneName, &sz.ZoneLat, &sz.ZoneLon, &sz.IncidentTypeID); err != nil {
			return nil, fmt.Errorf("failed to scan safe zone: %v", err)
		}
		safeZones = append(safeZones, sz)
	}
	return safeZones, nil
}
