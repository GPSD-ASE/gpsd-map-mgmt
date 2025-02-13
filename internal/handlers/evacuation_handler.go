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
func NewEvacuationHandler(service EvacuationServiceInterface) *EvacuationHandler {
	return &EvacuationHandler{Service: service}
}

// EvacuationRequest defines the expected JSON payload for an evacuation request.
// danger_point: coordinates of the danger zone.
// incident_type_id: used to match a safe zone in the database.
// safe_point: optional; if not provided, the nearest safe zone matching the incident type is used.
type EvacuationRequest struct {
	DangerPoint    [2]float64  `json:"danger_point"`
	IncidentTypeID int         `json:"incident_type_id"`
	SafePoint      *[2]float64 `json:"safe_point,omitempty"`
}

// GetEvacuationRoute handles the HTTP POST request to get an evacuation route.
// It validates the incoming JSON payload, calls the EvacuationService, and returns the route.
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
