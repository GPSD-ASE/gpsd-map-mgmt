package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockSafeZoneService struct{}

func (m *MockSafeZoneService) CreateSafeZone(sz models.SafeZoneCreate) (int, error) {
	// Just pretend we inserted row ID 42
	return 42, nil
}
func (m *MockSafeZoneService) GetSafeZones() ([]models.SafeZone, error) {
	zones := []models.SafeZone{
		{ZoneID: 1, ZoneName: "Flood Safe Zone", ZoneLat: 53.349805, ZoneLon: -6.26031, IncidentTypeID: 1},
		{ZoneID: 1, ZoneName: "EarthQuake Safe Zone", ZoneLat: 53.349805, ZoneLon: -6.26031, IncidentTypeID: 1},
	}

	return zones, nil
}
func TestCreateSafeZone_Happy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockSafeZoneService{}
	handler := handlers.NewSafeZoneHandler(mockSvc)

	router := gin.Default()
	router.POST("/safezones", handler.CreateSafeZone)

	payload := map[string]interface{}{
		"zone_name":        "Temporary Shelter 1",
		"zone_lat":         53.12345,
		"zone_lon":         -6.98765,
		"incident_type_id": 3,
	}
	reqBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/safezones", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}
