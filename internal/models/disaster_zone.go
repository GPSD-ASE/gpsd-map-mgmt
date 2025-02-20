package models

// DisasterZone represents a disaster zone.
// swagger:model DisasterZone
type DisasterZone struct {
	// The unique identifier of the incident.
	IncidentID int `json:"incident_id" example:"int"`
	// The name of the incident.
	IncidentName string `json:"incident_name" example:"Flood Zone"`
	// The latitude of the incident location.
	Latitude float64 `json:"latitude" example:"53.349805"`
	// The longitude of the incident location.
	Longitude float64 `json:"longitude" example:"-6.26031"`
	// The calculated radius for the disaster zone.
	Radius float64 `json:"radius" example:"30.5"`
}
