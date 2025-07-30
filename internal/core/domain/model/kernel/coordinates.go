package kernel

import (
	"fmt"
	"math"
)

const (
	MinLatitude   = -90.0
	MaxLatitude   = 90.0
	MinLongitude  = -180.0
	MaxLongitude  = 180.0
	earthRadiusKm = 6371.0
)

type GeoCoordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// NewGeoCoordinate creates a new coordinate with validation.
func NewGeoCoordinate(lat, lon float64) (GeoCoordinate, error) {
	if lat < MinLatitude || lat > MaxLatitude {
		return GeoCoordinate{}, fmt.Errorf("latitude %.6f is out of range (%f–%f)", lat, MinLatitude, MaxLatitude)
	}
	if lon < MinLongitude || lon > MaxLongitude {
		return GeoCoordinate{}, fmt.Errorf("longitude %.6f is out of range (%f–%f)", lon, MinLongitude, MaxLongitude)
	}
	return GeoCoordinate{Lat: lat, Lon: lon}, nil
}

func (g GeoCoordinate) Latitude() float64 {
	return g.Lat
}

func (g GeoCoordinate) Longitude() float64 {
	return g.Lon
}

// DistanceTo calculates the great-circle distance to another coordinate in kilometers.
func (g GeoCoordinate) DistanceTo(other GeoCoordinate) float64 {
	dLat := (other.Lat - g.Lat) * math.Pi / 180.0
	dLon := (other.Lon - g.Lon) * math.Pi / 180.0

	lat1Rad := g.Lat * math.Pi / 180.0
	lat2Rad := other.Lat * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func (g GeoCoordinate) Equals(other GeoCoordinate) bool {
	return g.Lat == other.Lat && g.Lon == other.Lon
}
