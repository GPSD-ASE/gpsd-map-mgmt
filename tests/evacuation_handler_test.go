package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockEvacuationService implements the EvacuationServiceInterface.
type MockEvacuationService struct{}

// GetEvacuationRoute returns a dummy full route response.
func (m *MockEvacuationService) GetEvacuationRoute(dangerPoint [2]float64, incidentTypeID int, safePoint *[2]float64) (services.EvacuationRouteResponse, error) {
	return services.EvacuationRouteResponse{
		Hints: map[string]interface{}{"sample_hint": "value"},
		Info:  map[string]interface{}{"took": 1},
		Paths: []services.RoutePath{
			{
				Distance:      1060.843,
				Weight:        307.85,
				Time:          780435,
				Transfers:     0,
				PointsEncoded: false,
				BBox:          []float64{11.539424, 48.118343, 11.558901, 48.122364},
				Points: services.GeoJSON{
					Type: "LineString",
					Coordinates: [][]float64{
						{11.539424, 48.118352},
						{11.540387, 48.118368},
					},
				},
				Instructions: []services.Instruction{
					{
						Distance:   661.29,
						Heading:    88.83,
						Sign:       0,
						Interval:   []int{0, 10},
						Text:       "Continue onto Lindenschmitstraße",
						Time:       238065,
						StreetName: "Lindenschmitstraße",
					},
				},
				Legs:             []interface{}{},
				Details:          map[string]interface{}{"road_class": "residential"},
				Ascend:           10.0,
				Descend:          5.0,
				SnappedWaypoints: services.GeoJSON{Type: "LineString", Coordinates: [][]float64{{11.539424, 48.118352}, {11.558901, 48.122364}}},
			},
		},
	}, nil
}

func TestGetEvacuationRouteHandler_Success(t *testing.T) {
	// Set Gin to test mode.
	gin.SetMode(gin.TestMode)

	// Create a mock evacuation service.
	mockService := &MockEvacuationService{}

	// Create the evacuation handler using the interface.
	evacuationHandler := handlers.NewEvacuationHandler(mockService)

	// Setup a Gin router and register the evacuation endpoint.
	router := gin.Default()
	router.POST("/evacuation", evacuationHandler.GetEvacuationRoute)

	// Prepare a sample request payload.
	payload := map[string]interface{}{
		"danger_point":     []float64{53.349805, -6.26031},
		"incident_type_id": 3,
		// safe_point is omitted so that the service uses the nearest safe zone logic.
	}
	jsonPayload, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Create a new HTTP request.
	req, err := http.NewRequest("POST", "/evacuation", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response.
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response status is 200 OK.
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response JSON.
	var response services.EvacuationRouteResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert that we received at least one route path with details.
	assert.Greater(t, len(response.Paths), 0)
	path := response.Paths[0]
	assert.NotEmpty(t, path.Points.Type)
	assert.NotNil(t, path.Instructions)
	// Additional assertions can check for expected values in the response.
}
