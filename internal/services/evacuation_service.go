package services

import (
	"database/sql"
	"fmt"
)

// EvacuationService encapsulates logic for determining a safe zone and fetching an evacuation route.
type EvacuationService struct {
	DB *sql.DB
	GH GraphHopperServiceInterface
}

// NewEvacuationService creates a new EvacuationService.
func NewEvacuationService(db *sql.DB, gh GraphHopperServiceInterface) *EvacuationService {
	return &EvacuationService{
		DB: db,
		GH: gh,
	}
}

// getNearestSafeZone queries the safe_zone table for the nearest safe zone matching the incident type.
// This query uses the haversine formula calculation directly in SQL for more precise distance computation.
func (s *EvacuationService) getNearestSafeZone(dangerPoint [2]float64, incidentTypeID int) (zoneLat, zoneLon float64, err error) {
	query := `
        SELECT zone_lat, zone_lon
        FROM safe_zone
        WHERE incident_type_id = $1
        ORDER BY acos(
            sin($2 * pi()/180) * sin(zone_lat * pi()/180) +
            cos($2 * pi()/180) * cos(zone_lat * pi()/180) *
            cos(($3 - zone_lon) * pi()/180)
        ) * 6371000 ASC
        LIMIT 1
    `
	row := s.DB.QueryRow(query, incidentTypeID, dangerPoint[0], dangerPoint[1])
	if err := row.Scan(&zoneLat, &zoneLon); err != nil {
		return 0, 0, fmt.Errorf("failed to find nearest safe zone: %v", err)
	}
	return zoneLat, zoneLon, nil
}

// GetEvacuationRoute returns the evacuation route from the danger point to a safe zone.
// If safePoint is nil, it finds the nearest safe zone for the given incident type.
func (s *EvacuationService) GetEvacuationRoute(dangerPoint [2]float64, incidentTypeID int, safePoint *[2]float64) (EvacuationRouteResponse, error) {
	var destination [2]float64
	if safePoint == nil {
		lat, lon, err := s.getNearestSafeZone(dangerPoint, incidentTypeID)
		if err != nil {
			return EvacuationRouteResponse{}, err
		}
		destination = [2]float64{lat, lon}
	} else {
		destination = *safePoint
	}
	return s.GH.GetEvacuationRoute(dangerPoint, destination)
}
