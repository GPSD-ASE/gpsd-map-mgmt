// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package services

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type TrafficServiceInterface interface {
	GetTrafficData(lat, lon string) ([]byte, error)
}

type TrafficService struct {
	APIKey  string
	BaseURL string
}

type TrafficResponse struct {
	FlowSegmentData struct {
		Coordinates        [][]float64 `json:"coordinates"`
		CurrentSpeed       float64     `json:"currentSpeed"`
		FreeFlowSpeed      float64     `json:"freeFlowSpeed"`
		CurrentTravelTime  int         `json:"currentTravelTime"`
		FreeFlowTravelTime int         `json:"freeFlowTravelTime"`
		Confidence         int         `json:"confidence"`
		RoadClosure        bool        `json:"roadClosure"`
	} `json:"flowSegmentData"`
}

func NewTrafficService(url string, apiKey string) *TrafficService {
	return &TrafficService{
		APIKey:  apiKey,
		BaseURL: url,
	}
}

func (s *TrafficService) GetTrafficData(lat, lon string) ([]byte, error) {
	url := fmt.Sprintf("%s?key=%s&point=%s,%s", s.BaseURL, s.APIKey, lat, lon)
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("TomTom API error: %s", resp.Status())
	}
	return resp.Body(), nil
}
