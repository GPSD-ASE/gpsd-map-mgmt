package services

import (
	"database/sql"
	"disaster-response-map-api/internal/models"
	"log"
)

// DisasterZoneServiceInterface defines methods to fetch disaster zones.
type DisasterZoneServiceInterface interface {
	GetDisasterZones() ([]models.DisasterZone, error)
}

// DisasterZoneService implements DisasterZoneServiceInterface.
type DisasterZoneService struct {
	DB *sql.DB
}

// NewDisasterZoneService creates a new instance of DisasterZoneService.
func NewDisasterZoneService(db *sql.DB) *DisasterZoneService {
	return &DisasterZoneService{DB: db}
}

// GetDisasterZones retrieves disaster zones from the database.
func (s *DisasterZoneService) GetDisasterZones() ([]models.DisasterZone, error) {
	rows, err := s.DB.Query("SELECT i.incident_id as incident_id, t.type_name AS incident_name, i.latitude as latitude, i.longitude as longitude, i.severity_level_id as severity_id FROM gpsd_inc.incident i JOIN gpsd_inc.incident_type t ON i.incident_type_id = t.type_id; ")
	if err != nil {
		// Log the error for debugging.
		log.Printf("Error querying disaster zones: %v", err)
		return nil, err
	}
	defer rows.Close()
	var zones []models.DisasterZone
	for rows.Next() {
		var dz models.DisasterZone
		var severityID int
		if err := rows.Scan(&dz.IncidentID, &dz.IncidentName, &dz.Latitude, &dz.Longitude, &severityID); err != nil {
			// Log the error but continue.
			log.Printf("Error scanning row: %v", err)
			continue
		}
		// Compute radius (modify as needed).
		dz.Radius = float64(severityID) * (8 + 10)
		zones = append(zones, dz)
	}
	return zones, nil
}
