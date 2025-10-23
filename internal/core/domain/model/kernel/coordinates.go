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
	epsilon       = 1e-12
)

type GeoCoordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// BoundingBox represents a geographical bounding box.
type BoundingBox struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
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

// BoundingBoxForRadius calculates a bounding box for the given radius in kilometers.
// This provides an approximate rectangular area that encompasses all points within the radius.
//
// Near the poles (|latitude| ≈ 90°), the longitude radius approaches infinity as
// meridians converge. In these cases the longitude radius is clamped to cover the
// full range of longitudes (180°).
func (g GeoCoordinate) BoundingBoxForRadius(radiusKm float64) BoundingBox {
	// Calculate latitude radius in degrees (1 degree latitude ≈ 111.32 km)
	latRadiusInDegrees := radiusKm / 111.32

	// Calculate longitude radius in degrees, adjusted for latitude distortion
	// At latitude, 1 degree longitude = cos(latitude) * 111.32 km
	cosLat := math.Cos(g.Lat * math.Pi / 180.0)
	var lonRadiusInDegrees float64
	if math.Abs(cosLat) < epsilon {
		lonRadiusInDegrees = 180.0
	} else {
		lonRadiusInDegrees = radiusKm / (111.32 * cosLat)
	}

	return BoundingBox{
		MinLat: g.Lat - latRadiusInDegrees,
		MaxLat: g.Lat + latRadiusInDegrees,
		MinLon: g.Lon - lonRadiusInDegrees,
		MaxLon: g.Lon + lonRadiusInDegrees,
	}
}

func (g GeoCoordinate) Equals(other GeoCoordinate) bool {
	return g.Lat == other.Lat && g.Lon == other.Lon
}
