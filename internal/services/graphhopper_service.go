package services

import (
	"bytes"
	"disaster-response-map-api/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GraphHopperServiceInterface defines the methods provided by a GraphHopper service.
type GraphHopperServiceInterface interface {
	GetEvacuationRoute(dangerPoint, safePoint [2]float64) (EvacuationRouteResponse, error)
	GetSafeRoute(origin, destination string, zones []models.DisasterZone) (RouteResponse, error)
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

func parseCoordinates(coord string) (float64, float64, error) {
	parts := strings.Split(coord, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate format: %s", coord)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %s", parts[0])
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %s", parts[1])
	}
	return lat, lon, nil
}
func buildPoints(origin, destination string) ([][]float64, error) {
	lat1, lon1, err := parseCoordinates(origin)
	if err != nil {
		return nil, err
	}
	lat2, lon2, err := parseCoordinates(destination)
	if err != nil {
		return nil, err
	}
	// Create an array of points in [longitude, latitude] order.
	points := [][]float64{
		{lon1, lat1},
		{lon2, lat2},
	}
	return points, nil
}

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

// GetSafeRoute calculates a route that avoids disaster zones using a custom model.
// The custom model penalizes any road segments that fall within a defined disaster zone.
func (s *GraphHopperService) GetSafeRoute(origin, destination string, zones []models.DisasterZone) (RouteResponse, error) {
	// Build the custom model based on provided disaster zones.
	customModel := BuildDisasterZonesCustomModel(zones)
	custjsonb, err := json.Marshal(customModel)
	if err != nil {
		return RouteResponse{}, err
	}
	log.Println("Custom Model:", string(custjsonb))
	// Use buildPoints to parse origin and destination.
	points, err := buildPoints(origin, destination)
	if err != nil {
		return RouteResponse{}, err
	}

	// Build the request payload including the custom model and disable speed mode.
	requestPayload := map[string]interface{}{
		"points":         points,
		"profile":        "car",
		"locale":         "en",
		"custom_model":   customModel,
		"ch.disable":     true,
		"instructions":   true,
		"calc_points":    true,
		"points_encoded": false,
	}

	jsonBytes, _ := json.Marshal(requestPayload)
	if err != nil {
		return RouteResponse{}, err
	}

	// Optionally log the payload for debugging.
	// log.Println("Request Payload:", string(jsonBytes))

	// Construct the URL with the API key as a query parameter.
	url := fmt.Sprintf("%s?key=%s", s.BaseURL, s.APIKey)
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return RouteResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return RouteResponse{}, fmt.Errorf("GraphHopper API error: %s - %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RouteResponse{}, err
	}

	var routeResp RouteResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return RouteResponse{}, err
	}

	return routeResp, nil
}
