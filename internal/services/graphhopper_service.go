// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package services

import (
	"bytes"
	"disaster-response-map-api/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type GraphHopperServiceInterface interface {
	GetEvacuationRoute(dangerPoint, safePoint [2]float64) (EvacuationRouteResponse, error)
	GetSafeRoute(origin, destination string, zones []models.DisasterZone) (RouteResponse, error)
}

type GraphHopperService struct {
	APIKey  string
	BaseURL string
}

func NewGraphHopperService(apiKey string, url string) *GraphHopperService {
	return &GraphHopperService{
		APIKey:  apiKey,
		BaseURL: url,
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
	points := [][]float64{
		{lon1, lat1},
		{lon2, lat2},
	}
	return points, nil
}

func (s *GraphHopperService) GetEvacuationRoute(dangerPoint, safePoint [2]float64) (EvacuationRouteResponse, error) {
	points := []interface{}{
		[]float64{dangerPoint[1], dangerPoint[0]}, // [lon, lat]
		[]float64{safePoint[1], safePoint[0]},     // [lon, lat]
	}

	requestPayload := map[string]interface{}{
		"points":           points,
		"snap_preventions": []string{"motorway", "ferry", "tunnel"},
		"details":          []string{"road_class", "surface"},
		"profile":          "foot",
		"locale":           "en",
		"instructions":     true,
		"calc_points":      true,
		"points_encoded":   false,
	}

	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return EvacuationRouteResponse{}, err
	}

	url := fmt.Sprintf("%s?key=%s", s.BaseURL, s.APIKey)
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return EvacuationRouteResponse{}, err
	}
	defer resp.Body.Close()

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

func (s *GraphHopperService) GetSafeRoute(origin, destination string, zones []models.DisasterZone) (RouteResponse, error) {
	customModel := BuildDisasterZonesCustomModel(zones)
	points, err := buildPoints(origin, destination)
	if err != nil {
		return RouteResponse{}, err
	}

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

	jsonBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return RouteResponse{}, err
	}

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
