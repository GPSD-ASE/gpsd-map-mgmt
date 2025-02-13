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

// EvacuationServiceInterface defines the method used by the handler.
type EvacuationServiceInterface interface {
	GetEvacuationRoute(dangerPoint [2]float64, incidentTypeID int, safePoint *[2]float64) (services.EvacuationRouteResponse, error)
}

// EvacuationHandler handles evacuation route requests.
type EvacuationHandler struct {
	Service EvacuationServiceInterface
}

// NewEvacuationHandler creates a new instance of EvacuationHandler.
// @Summary Create Evacuation Handler
// @Description Returns a new instance of EvacuationHandler.
// @Tags Evacuation
func NewEvacuationHandler(service EvacuationServiceInterface) *EvacuationHandler {
	return &EvacuationHandler{Service: service}
}

// EvacuationRequest defines the expected JSON payload for an evacuation request.
// swagger:model EvacuationRequest
type EvacuationRequest struct {
	// Coordinates of the danger zone in [latitude, longitude] format.
	// example: [53.349805, -6.26031]
	DangerPoint [2]float64 `json:"danger_point" example:"[53.349805, -6.26031]"`
	// Incident type ID used to match a safe zone in the database.
	// example: 3
	IncidentTypeID int `json:"incident_type_id" example:"3"`
	// (Optional) Coordinates of the safe zone in [latitude, longitude] format.
	// example: [53.3440, -6.2670]
	SafePoint *[2]float64 `json:"safe_point,omitempty" example:"[53.3440, -6.2670]"`
}

// GetEvacuationRoute godoc
// @Summary      Calculate Evacuation Route
// @Description  Calculates an evacuation route from a danger point to a safe zone. If safe_point is omitted, the API determines the nearest safe zone matching the incident type.
// @Tags         Evacuation
// @Accept       json
// @Produce      json
// @Param        evacuationRequest  body      EvacuationRequest  true  "Evacuation Request"
// @Success      200  {object}  services.EvacuationRouteResponse
// @Failure      400  {object}  map[string]string  "Invalid request payload"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /evacuation [post]
func (h *EvacuationHandler) GetEvacuationRoute(c *gin.Context) {
	var req EvacuationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	route, err := h.Service.GetEvacuationRoute(req.DangerPoint, req.IncidentTypeID, req.SafePoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}
