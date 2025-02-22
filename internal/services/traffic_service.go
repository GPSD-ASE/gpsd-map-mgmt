package services

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// TrafficServiceInterface defines the methods for retrieving traffic data.
type TrafficServiceInterface interface {
	GetTrafficData(lat, lon string) ([]byte, error)
}

// TrafficService implements TrafficServiceInterface.
type TrafficService struct {
	APIKey  string
	BaseURL string
}
type TrafficResponse struct {
}

// NewTrafficService creates a new instance of TrafficService.
func NewTrafficService(url string, apiKey string) *TrafficService {
	return &TrafficService{
		APIKey:  apiKey,
		BaseURL: url,
	}
}

// GetTrafficData fetches traffic data from the TomTom API using provided latitude and longitude.
func (s *TrafficService) GetTrafficData(lat, lon string) ([]byte, error) {
	// Build the TomTom API URL.
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
