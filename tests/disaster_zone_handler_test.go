package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockDisasterZoneService struct{}

func (m *MockDisasterZoneService) GetDisasterZones() ([]models.DisasterZone, error) {
	zones := []models.DisasterZone{
		{IncidentID: 1, IncidentName: "Flood Zone", Latitude: 53.349805, Longitude: -6.26031, Radius: 30.5},
		{IncidentID: 2, IncidentName: "Earthquake Zone", Latitude: 53.3478, Longitude: -6.2597, Radius: 25.0},
	}
	return zones, nil
}
func (m *MockDisasterZoneService) GetActiveDisasterZones() ([]models.DisasterZone, error) {
	zones := []models.DisasterZone{
		{IncidentID: 1, IncidentName: "Flood Zone", Latitude: 53.349805, Longitude: -6.26031, Radius: 30.5},
		{IncidentID: 2, IncidentName: "Earthquake Zone", Latitude: 53.3478, Longitude: -6.2597, Radius: 25.0},
	}
	return zones, nil
}

func TestGetDisasterZonesHandler_Happy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockDisasterZoneService{}
	handler := handlers.NewDisasterZoneHandler(mockService)

	router := gin.Default()
	router.GET("/zones", handler.GetDisasterZones)

	req, err := http.NewRequest(http.MethodGet, "/zones", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var zones []models.DisasterZone
	err = json.Unmarshal(recorder.Body.Bytes(), &zones)
	assert.NoError(t, err)
	assert.Len(t, zones, 2)
}
