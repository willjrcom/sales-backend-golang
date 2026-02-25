package addressentity

import "math"

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// CalculateDistance returns the distance between two coordinates in kilometers.
func (c Coordinates) CalculateDistance(other Coordinates) float64 {
	const earthRadius = 6371.0 // Radius of the Earth in kilometers

	lat1 := c.Latitude * math.Pi / 180
	lon1 := c.Longitude * math.Pi / 180
	lat2 := other.Latitude * math.Pi / 180
	lon2 := other.Longitude * math.Pi / 180

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	cVal := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * cVal
}
