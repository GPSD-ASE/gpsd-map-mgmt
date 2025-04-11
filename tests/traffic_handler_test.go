package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disaster-response-map-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockTrafficService struct{}

func (m *MockTrafficService) GetTrafficData(lat, lon string) ([]byte, error) {
	sample := map[string]interface{}{
		"flowSegmentData": map[string]interface{}{
			"coordinates":        [][]float64{{-6.26031, 53.349805}},
			"currentSpeed":       50.0,
			"freeFlowSpeed":      60.0,
			"currentTravelTime":  100,
			"freeFlowTravelTime": 90,
			"confidence":         80,
			"roadClosure":        false,
		},
	}
	return json.Marshal(sample)
}

func TestGetTrafficDataHandler_Happy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockTrafficService{}
	handler := handlers.NewTrafficHandler(mockService)

	router := gin.Default()
	router.GET("/traffic", handler.GetTrafficData)

	req, err := http.NewRequest(http.MethodGet, "/traffic?lat=53.349805&lon=-6.26031", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	flowData, exists := response["flowSegmentData"].(map[string]interface{})
	assert.True(t, exists)
	assert.Equal(t, float64(50.0), flowData["currentSpeed"])
}
