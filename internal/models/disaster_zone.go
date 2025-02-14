package models

// DisasterZone represents a disaster zone.
type DisasterZone struct {
	IncidentID   int     `json:"incident_id" example:"1"`
	IncidentName string  `json:"incident_name" example:"Flood Zone"`
	Latitude     float64 `json:"latitude" example:"53.349805"`
	Longitude    float64 `json:"longitude" example:"-6.26031"`
	Radius       float64 `json:"radius" example:"30.5"`
}
