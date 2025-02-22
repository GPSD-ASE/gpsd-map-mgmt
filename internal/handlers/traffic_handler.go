// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package handlers

import (
	"net/http"

	"disaster-response-map-api/internal/services"

	"github.com/gin-gonic/gin"
)

type TrafficHandler struct {
	Service services.TrafficServiceInterface
}

func NewTrafficHandler(service services.TrafficServiceInterface) *TrafficHandler {
	return &TrafficHandler{Service: service}
}

// GetTrafficData godoc
// @Summary      Get real-time traffic data
// @Description  Fetches traffic data from TomTom API based on latitude and longitude.
// @Tags         Traffic
// @Produce      json
// @Param        lat query string true "Latitude" example("53.349805")
// @Param        lon query string true "Longitude" example("-6.26031")
// @Success      200 {object} services.TrafficResponse
// @Failure      400 {object} map[string]string "Latitude and Longitude are required"
// @Failure      500 {object} map[string]string "Failed to fetch traffic data"
// @Router       /traffic [get]
func (h *TrafficHandler) GetTrafficData(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")
	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and Longitude are required"})
		return
	}

	data, err := h.Service.GetTrafficData(lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch traffic data", "details": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
