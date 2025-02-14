package services

import "math"

// BuildCirclePolygon returns a GeoJSON polygon approximating a circle centered at (centerLat, centerLon)
// with the given radius in meters. Coordinates are in [longitude, latitude] order.
func BuildCirclePolygon(centerLat, centerLon, radius float64) [][][]float64 {
	const (
		earthRadius = 6371000 // in meters
		numPoints   = 36      // Increase for smoother circle
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
