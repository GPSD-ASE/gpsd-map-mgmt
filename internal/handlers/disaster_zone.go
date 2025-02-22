// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package handlers

import (
	"disaster-response-map-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DisasterZoneHandler struct {
	DZService services.DisasterZoneServiceInterface
}

// NewDisasterZoneHandler creates a new instance of DisasterZoneHandler.
// @Summary      Create DisasterZoneHandler
// @Description  Returns a new instance of DisasterZoneHandler.
// @Tags         DisasterZone
func NewDisasterZoneHandler(dzService services.DisasterZoneServiceInterface) *DisasterZoneHandler {
	return &DisasterZoneHandler{DZService: dzService}
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
	zones, err := h.DZService.GetDisasterZones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch disaster zones"})
		return
	}
	c.JSON(http.StatusOK, zones)
}
