package handlers

import (
	"disaster-response-map-api/pkg/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DisasterZoneHandler holds a reference to the database interface
type DisasterZoneHandler struct {
	DB database.DatabaseInterface
}

// NewDisasterZoneHandler initializes a new handler with dependency injection
func NewDisasterZoneHandler(db database.DatabaseInterface) *DisasterZoneHandler {
	return &DisasterZoneHandler{DB: db}
}

// CreateDisasterZone handles disaster zone creation
func (h *DisasterZoneHandler) CreateDisasterZone(c *gin.Context) {
	var zone struct {
		Name   string  `json:"name"`
		Lat    float64 `json:"latitude"`
		Lng    float64 `json:"longitude"`
		Radius float64 `json:"radius"`
	}

	if err := c.BindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := h.DB.Exec(`
		INSERT INTO disaster_zones (name, latitude, longitude, radius) 
		VALUES ($1, $2, $3, $4)`,
		zone.Name, zone.Lat, zone.Lng, zone.Radius)

	if err != nil {
		log.Println("Database Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create disaster zone"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Disaster zone created successfully"})
}

// GetDisasterZones retrieves all disaster zones
func (h *DisasterZoneHandler) GetDisasterZones(c *gin.Context) {
	rows, err := h.DB.Query("SELECT id, name, latitude, longitude, radius FROM disaster_zones")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch disaster zones"})
		return
	}
	defer rows.Close()

	var zones []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var latitude, longitude, radius float64
		if err := rows.Scan(&id, &name, &latitude, &longitude, &radius); err != nil {
			continue
		}
		zones = append(zones, map[string]interface{}{
			"id":        id,
			"name":      name,
			"latitude":  latitude,
			"longitude": longitude,
			"radius":    radius,
		})
	}

	c.JSON(http.StatusOK, zones)
}

// DeleteDisasterZone deletes a disaster zone by ID
func (h *DisasterZoneHandler) DeleteDisasterZone(c *gin.Context) {
	id := c.Param("id")

	// Execute delete query
	result, err := h.DB.Exec("DELETE FROM disaster_zones WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete disaster zone"})
		return
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Disaster zone not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Disaster zone deleted"})
}
