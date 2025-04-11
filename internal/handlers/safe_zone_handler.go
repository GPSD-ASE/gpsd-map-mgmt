package handlers

import (
	"net/http"

	"disaster-response-map-api/internal/models"
	"disaster-response-map-api/internal/services"

	"github.com/gin-gonic/gin"
)

type SafeZoneHandler struct {
	Service services.SafeZoneServiceInterface
}

func NewSafeZoneHandler(service services.SafeZoneServiceInterface) *SafeZoneHandler {
	return &SafeZoneHandler{Service: service}
}

type CreateSafeZoneRequest struct {
	ZoneName       string  `json:"zone_name"       example:"Safe Zone 1"`
	ZoneLat        float64 `json:"zone_lat"        example:"53.12345"`
	ZoneLon        float64 `json:"zone_lon"        example:"-6.98765"`
	IncidentTypeID int     `json:"incident_type_id" example:"3"`
}

// CreateSafeZone godoc
// @Summary      Create new safe zone
// @Description  Inserts a new safe zone record into the DB.
// @Tags         SafeZone
// @Accept       json
// @Produce      json
// @Param        safeZoneRequest  body      CreateSafeZoneRequest  true  "New Safe Zone Data"
// @Success      201  {object}  map[string]interface{}  "Creation success"
// @Failure      400  {object}  map[string]string       "Invalid payload"
// @Failure      500  {object}  map[string]string       "Internal server error"
// @Router       /safezones [post]
func (h *SafeZoneHandler) CreateSafeZone(c *gin.Context) {
	var req CreateSafeZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	safeZone := models.SafeZoneCreate{
		ZoneName:       req.ZoneName,
		ZoneLat:        req.ZoneLat,
		ZoneLon:        req.ZoneLon,
		IncidentTypeID: req.IncidentTypeID,
	}

	newID, err := h.Service.CreateSafeZone(safeZone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Safe zone created successfully",
		"zone_id":   newID,
		"zone_name": safeZone.ZoneName,
	})
}
