package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/models"
	"disaster-response-map-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockGraphHopperService implements the GraphHopperServiceInterface for testing happy scenarios.
type MockGraphHopperService struct{}

func (m *MockGraphHopperService) GetRoute(origin, destination string) (services.RouteResponse, error) {
	return services.RouteResponse{
		Hints: map[string]interface{}{"sample_hint": "default"},
		Info:  map[string]interface{}{"took": 2},
		Paths: []services.RoutePath{
			{
				Distance: 500.0,
				Time:     600,
				Points: services.GeoJSON{
					Type:        "LineString",
					Coordinates: [][]float64{{-6.26031, 53.349805}, {-6.2597, 53.3478}},
				},
				Instructions: []services.Instruction{
					{
						Text: "Proceed straight",
					},
				},
			},
		},
	}, nil
}

func (m *MockGraphHopperService) GetEvacuationRoute(dangerPoint, safePoint [2]float64) (services.EvacuationRouteResponse, error) {
	// Not used in the safe routing happy test below.
	return services.EvacuationRouteResponse{}, nil
}

func (m *MockGraphHopperService) GetSafeRoute(origin, destination string, zones []models.DisasterZone) (services.RouteResponse, error) {
	return services.RouteResponse{
		Hints: map[string]interface{}{"sample_hint": "safe"},
		Info:  map[string]interface{}{"took": 3},
		Paths: []services.RoutePath{
			{
				Distance: 800.0,
				Time:     900,
				Points: services.GeoJSON{
					Type:        "LineString",
					Coordinates: [][]float64{{-6.26031, 53.349805}, {-6.2597, 53.3478}},
				},
				Instructions: []services.Instruction{
					{
						Text: "Take the safe route",
					},
				},
			},
		},
	}, nil
}

// MockDisasterZoneServiceForActive implements GetActiveDisasterZones for safe routing.
type MockDisasterZoneServiceForActive struct{}

func (m *MockDisasterZoneServiceForActive) GetDisasterZones() ([]models.DisasterZone, error) {
	return nil, nil
}

func (m *MockDisasterZoneServiceForActive) GetActiveDisasterZones() ([]models.DisasterZone, error) {
	zones := []models.DisasterZone{
		{IncidentID: 1, IncidentName: "Flood Zone", Latitude: 53.349805, Longitude: -6.26031, Radius: 30.5},
	}
	return zones, nil
}

func TestGetDefaultRouteHandler_Happy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockGHService := &MockGraphHopperService{}
	// For default route, the disaster zone service is not used.
	mockDZService := &MockDisasterZoneServiceForActive{}

	handler := handlers.NewRoutingHandler(mockGHService, mockDZService)

	router := gin.Default()
	router.GET("/route", handler.GetDefaultRoute)

	req, err := http.NewRequest(http.MethodGet, "/route?origin=53.349805,-6.26031&destination=53.3478,-6.2597", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var route services.RouteResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &route)
	assert.NoError(t, err)
	assert.Greater(t, len(route.Paths), 0)
}

func TestGetSafeRoutingHandler_Happy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockGHService := &MockGraphHopperService{}
	mockDZService := &MockDisasterZoneServiceForActive{}

	handler := handlers.NewRoutingHandler(mockGHService, mockDZService)

	router := gin.Default()
	router.GET("/routing", handler.GetSafeRouting)

	req, err := http.NewRequest(http.MethodGet, "/routing?origin=53.349805,-6.26031&destination=53.3478,-6.2597", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var route services.RouteResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &route)
	assert.NoError(t, err)
	assert.Greater(t, len(route.Paths), 0)
}
