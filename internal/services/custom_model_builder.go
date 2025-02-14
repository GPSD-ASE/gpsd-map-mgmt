package services

import (
	"strconv"

	"disaster-response-map-api/internal/models"
)

// BuildDisasterZonesCustomModel builds a custom model for GraphHopper based on a slice of disaster zones.
// Each disaster zone is represented as a GeoJSON polygon, and a priority rule is added to block routes through that area.
func BuildDisasterZonesCustomModel(zones []models.DisasterZone) map[string]interface{} {
	features := []map[string]interface{}{}
	priorityRules := []map[string]string{}

	for _, zone := range zones {
		polygon := BuildCirclePolygon(zone.Latitude, zone.Longitude, zone.Radius)
		if len(polygon) == 0 {
			continue
		}
		featureID := "disaster_zone_" + strconv.Itoa(zone.IncidentID)
		feature := map[string]interface{}{
			"id":         featureID,
			"type":       "Feature",
			"properties": map[string]interface{}{},
			"geometry": map[string]interface{}{
				"type":        "Polygon",
				"coordinates": polygon[0],
			},
		}
		features = append(features, feature)
		priorityRules = append(priorityRules, map[string]string{
			"if":          "in_" + featureID,
			"multiply_by": "0",
		})
	}

	return map[string]interface{}{
		"priority": priorityRules,
		"areas": map[string]interface{}{
			"type":     "FeatureCollection",
			"features": features,
		},
	}
}
