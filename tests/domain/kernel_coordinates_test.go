package domain

// DOMAIN LAYER UNIT TESTS
// Tests for domain model business rules and validation logic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
)

func TestNewGeoCoordinate_Valid(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"Moscow center", 55.7558, 37.6176},
		{"North pole", 90.0, 0.0},
		{"South pole", -90.0, 0.0},
		{"Date line", 0.0, 180.0},
		{"Antimeridian", 0.0, -180.0},
		{"Zero coordinates", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord, err := kernel.NewGeoCoordinate(tt.lat, tt.lon)
			assert.NoError(t, err)
			assert.Equal(t, tt.lat, coord.Latitude())
			assert.Equal(t, tt.lon, coord.Longitude())
		})
	}
}

func TestNewGeoCoordinate_InvalidLatitude(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"Latitude too high", 91.0, 0.0},
		{"Latitude too low", -91.0, 0.0},
		{"Latitude way too high", 200.0, 0.0},
		{"Latitude way too low", -200.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kernel.NewGeoCoordinate(tt.lat, tt.lon)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "latitude")
			assert.Contains(t, err.Error(), "out of range")
		})
	}
}

func TestNewGeoCoordinate_InvalidLongitude(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"Longitude too high", 0.0, 181.0},
		{"Longitude too low", 0.0, -181.0},
		{"Longitude way too high", 0.0, 360.0},
		{"Longitude way too low", 0.0, -360.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := kernel.NewGeoCoordinate(tt.lat, tt.lon)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "longitude")
			assert.Contains(t, err.Error(), "out of range")
		})
	}
}

func TestGeoCoordinate_DistanceTo(t *testing.T) {
	// Moscow coordinates
	moscow := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}

	tests := []struct {
		name           string
		from           kernel.GeoCoordinate
		to             kernel.GeoCoordinate
		expectedDistKm float64
		tolerance      float64
	}{
		{
			name:           "Same coordinates",
			from:           moscow,
			to:             moscow,
			expectedDistKm: 0.0,
			tolerance:      0.001,
		},
		{
			name:           "Moscow to St. Petersburg",
			from:           moscow,
			to:             kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609},
			expectedDistKm: 635.0, // Approximate distance
			tolerance:      50.0,  // 50km tolerance
		},
		{
			name:           "Moscow to New York",
			from:           moscow,
			to:             kernel.GeoCoordinate{Lat: 40.7128, Lon: -74.0060},
			expectedDistKm: 7500.0, // Approximate distance
			tolerance:      200.0,  // 200km tolerance
		},
		{
			name:           "Equator points",
			from:           kernel.GeoCoordinate{Lat: 0.0, Lon: 0.0},
			to:             kernel.GeoCoordinate{Lat: 0.0, Lon: 1.0},
			expectedDistKm: 111.32, // 1 degree at equator
			tolerance:      5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			distance := tt.from.DistanceTo(tt.to)
			assert.InDelta(t, tt.expectedDistKm, distance, tt.tolerance,
				"Distance from %v to %v should be approximately %.1f km, got %.1f km",
				tt.from, tt.to, tt.expectedDistKm, distance)

			// Distance should be symmetric
			reverseDistance := tt.to.DistanceTo(tt.from)
			assert.InDelta(t, distance, reverseDistance, 0.001,
				"Distance should be symmetric")
		})
	}
}

func TestGeoCoordinate_BoundingBoxForRadius(t *testing.T) {
	tests := []struct {
		name     string
		center   kernel.GeoCoordinate
		radiusKm float64
	}{
		{
			name:     "Moscow 10km radius",
			center:   kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176},
			radiusKm: 10.0,
		},
		{
			name:     "Equator 5km radius",
			center:   kernel.GeoCoordinate{Lat: 0.0, Lon: 0.0},
			radiusKm: 5.0,
		},
		{
			name:     "High latitude 1km radius",
			center:   kernel.GeoCoordinate{Lat: 80.0, Lon: 0.0},
			radiusKm: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbox := tt.center.BoundingBoxForRadius(tt.radiusKm)

			// Basic validations
			assert.True(t, bbox.MinLat <= tt.center.Lat, "MinLat should be <= center latitude")
			assert.True(t, bbox.MaxLat >= tt.center.Lat, "MaxLat should be >= center latitude")
			assert.True(t, bbox.MinLon <= tt.center.Lon, "MinLon should be <= center longitude")
			assert.True(t, bbox.MaxLon >= tt.center.Lon, "MaxLon should be >= center longitude")

			// Ensure bounding box stays within valid coordinate ranges
			assert.True(t, bbox.MinLat >= -90.0, "MinLat should be >= -90")
			assert.True(t, bbox.MaxLat <= 90.0, "MaxLat should be <= 90")
			assert.True(t, bbox.MinLon >= -180.0, "MinLon should be >= -180")
			assert.True(t, bbox.MaxLon <= 180.0, "MaxLon should be <= 180")

			// Check that corners of bounding box are approximately within radius
			corners := []kernel.GeoCoordinate{
				{Lat: bbox.MinLat, Lon: bbox.MinLon},
				{Lat: bbox.MinLat, Lon: bbox.MaxLon},
				{Lat: bbox.MaxLat, Lon: bbox.MinLon},
				{Lat: bbox.MaxLat, Lon: bbox.MaxLon},
			}

			for i, corner := range corners {
				distance := tt.center.DistanceTo(corner)
				// Corner should be within reasonable distance (allowing some tolerance for rectangular approximation)
				assert.True(t, distance <= tt.radiusKm*1.5,
					"Corner %d distance %.2f km should be reasonably close to radius %.2f km", i, distance, tt.radiusKm)
			}
		})
	}
}

func TestGeoCoordinate_BoundingBoxForRadius_PolarRegions(t *testing.T) {
	// Test near poles where longitude calculations become extreme
	nearNorthPole := kernel.GeoCoordinate{Lat: 89.0, Lon: 0.0}
	bbox := nearNorthPole.BoundingBoxForRadius(100.0) // 100km radius

	// At high latitudes, longitude range should be very large or clamped
	lonRange := bbox.MaxLon - bbox.MinLon
	assert.True(t, lonRange > 90.0, "Longitude range should be large near poles")

	// Verify bounding box is still valid
	assert.True(t, bbox.MinLat >= -90.0 && bbox.MaxLat <= 90.0, "Latitude range should be valid")
	assert.True(t, bbox.MinLon >= -180.0 && bbox.MaxLon <= 180.0, "Longitude range should be valid")
}

func TestGeoCoordinate_Equals(t *testing.T) {
	coord1 := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	coord2 := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	coord3 := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6176} // Slightly different

	assert.True(t, coord1.Equals(coord2), "Identical coordinates should be equal")
	assert.True(t, coord2.Equals(coord1), "Equals should be symmetric")
	assert.False(t, coord1.Equals(coord3), "Different coordinates should not be equal")
}

func TestGeoCoordinate_DistanceTo_EdgeCases(t *testing.T) {
	// Test distance across date line
	west := kernel.GeoCoordinate{Lat: 0.0, Lon: 179.0}
	east := kernel.GeoCoordinate{Lat: 0.0, Lon: -179.0}

	distance := west.DistanceTo(east)
	expected := 2.0 * 111.32 // 2 degrees at equator
	assert.InDelta(t, expected, distance, 10.0, "Distance across date line should be calculated correctly")

	// Test very small distances
	coord1 := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	coord2 := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177} // 0.0001 degree difference

	smallDistance := coord1.DistanceTo(coord2)
	assert.True(t, smallDistance > 0.0 && smallDistance < 1.0, "Small distance should be between 0 and 1 km")
}
