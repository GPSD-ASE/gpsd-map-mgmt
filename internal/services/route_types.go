// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package services

type GeoJSON struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"` // can be []float64 or [][]float64 depending on the geometry
}

type Instruction struct {
	Distance   float64 `json:"distance"`
	Heading    float64 `json:"heading"`
	Sign       int     `json:"sign"`
	Interval   []int   `json:"interval"`
	Text       string  `json:"text"`
	Time       int     `json:"time"`
	StreetName string  `json:"street_name"`
}

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

type RouteResponse struct {
	Hints map[string]interface{} `json:"hints" example:"{\"visited_nodes.sum\": 100, \"visited_nodes.average\": 100}"`
	Info  map[string]interface{} `json:"info" example:"{\"took\": 3, \"copyrights\": [\"GraphHopper\", \"OpenStreetMap contributors\"]}"`
	Paths []RoutePath            `json:"paths"`
}

type EvacuationRouteResponse struct {
	Hints map[string]interface{} `json:"hints"`
	Info  map[string]interface{} `json:"info"`
	Paths []RoutePath            `json:"paths"`
}
