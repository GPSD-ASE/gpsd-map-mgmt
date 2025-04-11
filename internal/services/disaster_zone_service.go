// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package services

import (
	"database/sql"
	"disaster-response-map-api/internal/models"
	"log"
)

type DisasterZoneServiceInterface interface {
	GetDisasterZones() ([]models.DisasterZone, error)
	GetActiveDisasterZones() ([]models.DisasterZone, error)
}

type DisasterZoneService struct {
	DB *sql.DB
}

func NewDisasterZoneService(db *sql.DB) *DisasterZoneService {
	return &DisasterZoneService{DB: db}
}

func (s *DisasterZoneService) GetDisasterZones() ([]models.DisasterZone, error) {
	rows, err := s.DB.Query("SELECT t.type_name AS incident_name, i.latitude as latitude, i.longitude as longitude, i.severity_id as severity_id FROM incident i JOIN incident_type t ON i.type_id = t.type_id; ")
	if err != nil {
		log.Printf("Error querying disaster zones: %v", err)
		return nil, err
	}
	defer rows.Close()
	var zones []models.DisasterZone
	increment := 1
	for rows.Next() {
		var dz models.DisasterZone
		var severityID int
		if err := rows.Scan(&dz.IncidentName, &dz.Latitude, &dz.Longitude, &severityID); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		dz.IncidentID = increment
		increment++

		dz.Radius = float64(severityID) * (8 + 10)
		zones = append(zones, dz)
	}
	return zones, nil
}

func (s *DisasterZoneService) GetActiveDisasterZones() ([]models.DisasterZone, error) {
	rows, err := s.DB.Query("SELECT t.type_name AS incident_name, i.latitude as latitude, i.longitude as longitude, i.severity_id as severity_id FROM incident i JOIN incident_type t ON i.type_id = t.type_id WHERE status_id = 3; ")
	if err != nil {
		log.Printf("Error querying disaster zones: %v", err)
		return nil, err
	}
	defer rows.Close()
	var zones []models.DisasterZone
	increment := 1
	for rows.Next() {
		var dz models.DisasterZone
		var severityID int
		if err := rows.Scan(&dz.IncidentName, &dz.Latitude, &dz.Longitude, &severityID); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		dz.IncidentID = increment
		increment++

		dz.Radius = float64(severityID) * (8 + 10)
		zones = append(zones, dz)
	}
	return zones, nil
}
