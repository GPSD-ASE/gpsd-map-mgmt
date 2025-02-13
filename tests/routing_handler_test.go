package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"disaster-response-map-api/internal/handlers"
	"disaster-response-map-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRouting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock service instance.
	mockService := &MockGraphHopperService{}

	// Create the routing handler using the mock service.
	routingHandler := handlers.NewRoutingHandler(mockService)

	// Set up the Gin router and endpoint.
	router := gin.Default()
	router.GET("/routing", routingHandler.GetRouting)

	// Simulate a request with query parameters for origin and destination.
	req, _ := http.NewRequest("GET", "/routing?origin=53.349805,-6.26031&destination=53.3478,-6.2597", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect HTTP 200 OK.
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse and validate the JSON response.
	var response services.RouteResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1000), response.Distance)
	assert.Equal(t, 600, response.Time)
}
