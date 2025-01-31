package handlers

import (
	"log"
	"net/http"
	"disaster-response-map-api/pkg/database"
	"github.com/gin-gonic/gin"
)

type DisasterZone struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"` // In meters
}

// Create Disaster Zone
// Create a disaster zone
func CreateDisasterZone(c *gin.Context) {
	var zone struct {
		Name    string  `json:"name"`
		Lat     float64 `json:"latitude"`
		Lng     float64 `json:"longitude"`
		Radius  float64 `json:"radius"`
	}

	if err := c.BindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Insert into PostgreSQL
	_, err := database.DB.Exec(`
		INSERT INTO disaster_zones (name, latitude, longitude, radius) 
		VALUES ($1, $2, $3, $4)`, 
		zone.Name, zone.Lat, zone.Lng, zone.Radius)

	if err != nil {
		// Log the exact error to debug
		log.Println("Database Insert Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create disaster zone", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Disaster zone created successfully"})
}


// Health Check Endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Get All Disaster Zones
func GetDisasterZones(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, latitude, longitude, radius FROM disaster_zones")
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

// Delete Disaster Zone
func DeleteDisasterZone(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM disaster_zones WHERE id = $1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete disaster zone"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Disaster zone deleted"})
}
