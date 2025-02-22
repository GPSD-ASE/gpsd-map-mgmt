// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package services

import (
	"disaster-response-map-api/internal/models"
	"log"
	"strconv"
)

func BuildDisasterZonesCustomModel(zones []models.DisasterZone) map[string]interface{} {
	features := []map[string]interface{}{}
	priorityRules := []map[string]interface{}{}
	for _, zone := range zones {
		log.Println("Building zone for incident: ", zone.IncidentID)
		polygon := BuildCirclePolygon(zone.Latitude, zone.Longitude, zone.Radius)
		if len(polygon) == 0 {
			continue
		}
		featureID := "disaster_zone_" + strconv.Itoa(zone.IncidentID)
		feature := map[string]interface{}{
			"id":   featureID,
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        "Polygon",
				"coordinates": polygon,
			},
		}
		features = append(features, feature)
		priorityRules = append(priorityRules, map[string]interface{}{
			"if":          "in_" + featureID,
			"multiply_by": 0,
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
