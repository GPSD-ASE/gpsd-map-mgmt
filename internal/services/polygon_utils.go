// @title Disaster Response Map API
// @version 1.0.0
// @description API for disaster response, including retrieval of disaster zones, routing between two points, and calculating evacuation routes.
// @contact.name Rokas Paulauskas
// @contact.email paulausr@tcd.ie
// @BasePath /
package services

import "math"

func BuildCirclePolygon(centerLat, centerLon, radius float64) [][][]float64 {
	const (
		earthRadius = 6371000
		numPoints   = 36
	)
	angleStep := 2 * math.Pi / numPoints
	var ring [][]float64

	latRad := centerLat * math.Pi / 180.0
	lonRad := centerLon * math.Pi / 180.0
	angularDistance := radius / earthRadius

	for i := 0; i < numPoints; i++ {
		bearing := float64(i) * angleStep
		destLatRad := math.Asin(math.Sin(latRad)*math.Cos(angularDistance) +
			math.Cos(latRad)*math.Sin(angularDistance)*math.Cos(bearing))
		destLonRad := lonRad + math.Atan2(math.Sin(bearing)*math.Sin(angularDistance)*math.Cos(latRad),
			math.Cos(angularDistance)-math.Sin(latRad)*math.Sin(destLatRad))
		destLat := destLatRad * 180.0 / math.Pi
		destLon := destLonRad * 180.0 / math.Pi
		ring = append(ring, []float64{destLon, destLat})
	}

	// Close the ring.
	if len(ring) > 0 {
		ring = append(ring, ring[0])
	}
	return [][][]float64{ring}
}
