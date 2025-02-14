// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package handlers

import (
	"disaster-response-map-api/pkg/database"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DisasterZone represents a disaster zone.
// swagger:model DisasterZone
type DisasterZone struct {
	// The unique identifier of the incident.
	IncidentID int `json:"incident_id" example:"1"`
	// The name of the incident.
	IncidentName string `json:"incident_name" example:"Flood Zone"`
	// The latitude of the incident location.
	Latitude float64 `json:"latitude" example:"53.349805"`
	// The longitude of the incident location.
	Longitude float64 `json:"longitude" example:"-6.26031"`
	// The calculated radius for the disaster zone.
	Radius float64 `json:"radius" example:"30.5"`
}

// DisasterZoneHandler handles disaster zone requests.
type DisasterZoneHandler struct {
	DB database.DatabaseInterface
}

// NewDisasterZoneHandler creates a new instance of DisasterZoneHandler.
// @Summary      Create DisasterZoneHandler
// @Description  Returns a new instance of DisasterZoneHandler.
// @Tags         DisasterZone
func NewDisasterZoneHandler(db database.DatabaseInterface) *DisasterZoneHandler {
	return &DisasterZoneHandler{DB: db}
}

// GetDisasterZones godoc
// @Summary      Retrieve Disaster Zones
// @Description  Retrieves a list of disaster zones from the database.
// @Tags         DisasterZone
// @Produce      json
// @Success      200  {array}   DisasterZone
// @Failure      500  {object}  map[string]string  "Internal Server Error"
// @Router       /zones [get]
func (h *DisasterZoneHandler) GetDisasterZones(c *gin.Context) {
	rows, err := h.DB.Query("SELECT i.incident_id as incident_id, t.type_name AS incident_name, i.latitude as latitude, i.longitude as longitude, i.severity_level_id as severity_id FROM gpsd_inc.incident i JOIN gpsd_inc.incident_type t ON i.incident_type_id = t.type_id; ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch disaster zones"})
		return
	}
	defer rows.Close()
	var zones []DisasterZone
	for rows.Next() {
		var dz DisasterZone
		var severityID int
		if err := rows.Scan(&dz.IncidentID, &dz.IncidentName, &dz.Latitude, &dz.Longitude, &severityID); err != nil {
			continue
		}
		dz.Radius = float64(severityID) * (8 + rand.Float64()*(15-8))
		zones = append(zones, dz)
	}

	c.JSON(http.StatusOK, zones)
}
