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

// RoutingHandler handles route requests.
type RoutingHandler struct {
	Service services.GraphHopperServiceInterface
}

// NewRoutingHandler creates a new instance of RoutingHandler.
// @Summary Create Routing Handler
// @Description Returns a new instance of RoutingHandler.
// @Tags Routing
func NewRoutingHandler(service services.GraphHopperServiceInterface) *RoutingHandler {
	return &RoutingHandler{Service: service}
}

// GetRouting godoc
// @Summary      Calculate Route
// @Description  Calculates a route between two points using the GraphHopper service.
// @Tags         Routing
// @Produce      json
// @Param        origin       query     string  true  "Origin coordinates in latitude,longitude format"  example("53.349805,-6.26031")
// @Param        destination  query     string  true  "Destination coordinates in latitude,longitude format"  example("53.3478,-6.2597")
// @Success      200  {object}  services.RouteResponse
// @Failure      400  {object}  map[string]string  "Missing required parameters"
// @Failure      500  {object}  map[string]string  "Failed to fetch route"
// @Router       /routing [get]
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
