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

type MockEvacuationService struct{}

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
	gin.SetMode(gin.TestMode)

	mockService := &MockEvacuationService{}

	evacuationHandler := handlers.NewEvacuationHandler(mockService)

	router := gin.Default()
	router.POST("/evacuation", evacuationHandler.GetEvacuationRoute)

	payload := map[string]interface{}{
		"danger_point":     []float64{53.349805, -6.26031},
		"incident_type_id": 3,
	}
	jsonPayload, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/evacuation", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response services.EvacuationRouteResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Greater(t, len(response.Paths), 0)
	path := response.Paths[0]
	assert.NotEmpty(t, path.Points.Type)
	assert.NotNil(t, path.Instructions)
}
