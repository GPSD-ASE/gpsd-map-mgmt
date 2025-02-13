package handlers

import (
	"disaster-response-map-api/pkg/database"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DisasterZoneHandler struct {
	DB database.DatabaseInterface
}

func NewDisasterZoneHandler(db database.DatabaseInterface) *DisasterZoneHandler {
	return &DisasterZoneHandler{DB: db}
}

func (h *DisasterZoneHandler) GetDisasterZones(c *gin.Context) {
	rows, err := h.DB.Query("SELECT i.incident_id as incident_id, t.type_name AS incident_name, i.latitude as latitude, i.longitude as longitude, i.severity_id as severity_id FROM incident i JOIN incident_type t ON i.type_id = t.type_id; ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch disaster zones"})
		return
	}
	defer rows.Close()
	var zones []map[string]interface{}
	for rows.Next() {
		var incident_id int
		var incident_name string
		var latitude, longitude float64
		var severity_id int
		if err := rows.Scan(&incident_id, &incident_name, &latitude, &longitude, &severity_id); err != nil {
			continue
		}
		zones = append(zones, map[string]interface{}{
			"incident_id":    incident_id,
			"inncident_name": incident_name,
			"latitude":       latitude,
			"longitude":      longitude,
			"radius":         float64(severity_id) * (8 + rand.Float64()*(15-8)),
		})
	}

	c.JSON(http.StatusOK, zones)
}
