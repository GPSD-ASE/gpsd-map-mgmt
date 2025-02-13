// File: internal/services/route_types.go
package services

// GeoJSON represents a GeoJSON geometry.
type GeoJSON struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"` // can be []float64 or [][]float64 depending on the geometry
}

// Instruction holds turn-by-turn navigation instructions.
type Instruction struct {
	Distance   float64 `json:"distance"`
	Heading    float64 `json:"heading"`
	Sign       int     `json:"sign"`
	Interval   []int   `json:"interval"`
	Text       string  `json:"text"`
	Time       int     `json:"time"`
	StreetName string  `json:"street_name"`
}

// RoutePath holds the complete details for a route path.
type RoutePath struct {
	Distance         float64                `json:"distance"`
	Weight           float64                `json:"weight"`
	Time             int                    `json:"time"`
	Transfers        int                    `json:"transfers"`
	PointsEncoded    bool                   `json:"points_encoded"`
	BBox             []float64              `json:"bbox"`
	Points           GeoJSON                `json:"points"`
	Instructions     []Instruction          `json:"instructions"`
	Legs             []interface{}          `json:"legs"`
	Details          map[string]interface{} `json:"details"`
	Ascend           float64                `json:"ascend"`
	Descend          float64                `json:"descend"`
	SnappedWaypoints GeoJSON                `json:"snapped_waypoints"`
}

// RouteResponse is the response from GraphHopper for a simple route request.
type RouteResponse struct {
	Distance float64 `json:"distance"`
	Time     int     `json:"time"`
}

// EvacuationRouteResponse represents the full response from GraphHopper for an evacuation route.
type EvacuationRouteResponse struct {
	Hints map[string]interface{} `json:"hints"`
	Info  map[string]interface{} `json:"info"`
	Paths []RoutePath            `json:"paths"`
}
