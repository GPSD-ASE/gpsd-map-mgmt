package models

type SafeZoneCreate struct {
	ZoneName       string  `json:"zone_name"`
	ZoneLat        float64 `json:"zone_lat"`
	ZoneLon        float64 `json:"zone_lon"`
	IncidentTypeID int     `json:"incident_type_id"`
}
type SafeZone struct {
	ZoneID         int     `json:"zone_id" example:"1"`
	ZoneName       string  `json:"zone_name" example:"Safe Zone 1"`
	ZoneLat        float64 `json:"zone_lat" example:"53.12345"`
	ZoneLon        float64 `json:"zone_lon" example:"-6.98765"`
	IncidentTypeID int     `json:"incident_type_id" example:"3"`
}
