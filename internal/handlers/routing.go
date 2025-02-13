package handlers

import (
	"net/http"

	"disaster-response-map-api/internal/services"

	"github.com/gin-gonic/gin"
)

type RoutingHandler struct {
	Service services.GraphHopperServiceInterface
}

func NewRoutingHandler(service services.GraphHopperServiceInterface) *RoutingHandler {
	return &RoutingHandler{Service: service}
}

func (h *RoutingHandler) GetRouting(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	// Delegate the responsibility of getting the route to the service.
	route, err := h.Service.GetRoute(origin, destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch route"})
		return
	}

	c.JSON(http.StatusOK, route)
}
