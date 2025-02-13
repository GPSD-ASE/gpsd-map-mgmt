package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GraphHopperServiceInterface defines the methods provided by a GraphHopper service.
type GraphHopperServiceInterface interface {
	GetRoute(origin, destination string) (RouteResponse, error)
	GetEvacuationRoute(dangerPoint, safePoint [2]float64) (EvacuationRouteResponse, error)
}

// GraphHopperService implements GraphHopperServiceInterface.
type GraphHopperService struct {
	APIKey  string
	BaseURL string
}

// NewGraphHopperService creates a new GraphHopperService with the provided API key.
func NewGraphHopperService(apiKey string) *GraphHopperService {
	return &GraphHopperService{
		APIKey:  apiKey,
		BaseURL: "https://graphhopper.com/api/1/route",
	}
}

// GetRoute is a dummy implementation for demonstration purposes.
func (s *GraphHopperService) GetRoute(origin, destination string) (RouteResponse, error) {
	return RouteResponse{
		Distance: 1000,
		Time:     600,
	}, nil
}

// GetEvacuationRoute calls the GraphHopper API with a POST request to get an evacuation route.
func (s *GraphHopperService) GetEvacuationRoute(dangerPoint, safePoint [2]float64) (EvacuationRouteResponse, error) {
	// IMPORTANT: The GraphHopper POST API expects points as [longitude, latitude].
	// We assume the user sends coordinates as [latitude, longitude],
	// so we swap the order here.
	points := []interface{}{
		[]float64{dangerPoint[1], dangerPoint[0]}, // [lon, lat]
		[]float64{safePoint[1], safePoint[0]},     // [lon, lat]
	}

	// Build the request payload using the provided documentation.
	requestPayload := map[string]interface{}{
		"points":           points,
		"snap_preventions": []string{"motorway", "ferry", "tunnel"},
		"details":          []string{"road_class", "surface"},
		"profile":          "foot", // using foot profile for a person evacuating
		"locale":           "en",
		"instructions":     true,
		"calc_points":      true,
		"points_encoded":   false,
	}

	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return EvacuationRouteResponse{}, err
	}

	// Construct the URL with the API key as a query parameter.
	url := fmt.Sprintf("%s?key=%s", s.BaseURL, s.APIKey)
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return EvacuationRouteResponse{}, err
	}
	defer resp.Body.Close()

	// If the response is not OK, read the body for debugging information.
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return EvacuationRouteResponse{}, fmt.Errorf("GraphHopper API error: %s - %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EvacuationRouteResponse{}, err
	}

	var routeResp EvacuationRouteResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return EvacuationRouteResponse{}, err
	}

	return routeResp, nil
}
