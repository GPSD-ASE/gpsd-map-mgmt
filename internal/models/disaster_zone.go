// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package models

// swagger:model DisasterZone
type DisasterZone struct {
	IncidentID   int     `json:"incident_id" example:"int"`
	IncidentName string  `json:"incident_name" example:"Flood Zone"`
	Latitude     float64 `json:"latitude" example:"53.349805"`
	Longitude    float64 `json:"longitude" example:"-6.26031"`
	Radius       float64 `json:"radius" example:"30.5"`
}
