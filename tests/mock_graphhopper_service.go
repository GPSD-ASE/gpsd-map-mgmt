package tests

import (
	"disaster-response-map-api/internal/services"
)

// MockGraphHopperService implements the GraphHopperServiceInterface.
type MockGraphHopperService struct{}

// GetRoute returns a dummy simple route response.
func (m *MockGraphHopperService) GetRoute(origin, destination string) (services.RouteResponse, error) {
	return services.RouteResponse{
		// Distance: 1000,
		// Time:     600,
	}, nil
}

// GetEvacuationRoute returns a dummy full evacuation route response.
func (m *MockGraphHopperService) GetEvacuationRoute(dangerPoint, safePoint [2]float64) (services.EvacuationRouteResponse, error) {
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
