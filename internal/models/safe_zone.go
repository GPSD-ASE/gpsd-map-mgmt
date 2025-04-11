package models

type SafeZoneCreate struct {
	ZoneName       string  `json:"zone_name"`
	ZoneLat        float64 `json:"zone_lat"`
	ZoneLon        float64 `json:"zone_lon"`
	IncidentTypeID int     `json:"incident_type_id"`
}
