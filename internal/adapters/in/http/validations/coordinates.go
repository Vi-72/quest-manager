package validations

import (
	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/core/domain/model/kernel"
)

// ValidatedSearchByRadiusData contains validated data for radius search
type ValidatedSearchByRadiusData struct {
	Center   kernel.GeoCoordinate
	RadiusKm float64
}

// ValidateSearchByRadiusParams validates and converts radius search parameters
// Note: Basic validations (ranges, coordinate bounds) are now handled by OpenAPI
func ValidateSearchByRadiusParams(lat, lon, radiusKm float32) (*ValidatedSearchByRadiusData, *ValidationError) {
	// Convert to float64 for domain layer
	latF64 := float64(lat)
	lonF64 := float64(lon)
	radiusF64 := float64(radiusKm)

	// OpenAPI now handles:
	// - lat: minimum -90, maximum 90
	// - lon: minimum -180, maximum 180
	// - radius_km: minimum 0.1, maximum 20000

	// Create domain coordinate - domain layer handles coordinate validation
	center, err := kernel.NewGeoCoordinate(latF64, lonF64)
	if err != nil {
		return nil, NewValidationErrorWithCause("coordinates", "invalid coordinate values", err)
	}

	return &ValidatedSearchByRadiusData{
		Center:   center,
		RadiusKm: radiusF64,
	}, nil
}

// ConvertAPICoordinateToKernel converts API coordinates to domain coordinates with validation
// Validation is performed by the domain layer (kernel.NewGeoCoordinate)
func ConvertAPICoordinateToKernel(apiCoord v1.Coordinate) (kernel.GeoCoordinate, *ValidationError) {
	// Create domain coordinate (convert float32 to float64)
	// Domain layer handles all coordinate validation
	coord, err := kernel.NewGeoCoordinate(float64(apiCoord.Latitude), float64(apiCoord.Longitude))
	if err != nil {
		return kernel.GeoCoordinate{}, NewValidationErrorWithCause("coordinate", "invalid coordinate values", err)
	}

	return coord, nil
}

// ConvertKernelCoordinateToAPI converts domain coordinates to API format
func ConvertKernelCoordinateToAPI(kernelCoord kernel.GeoCoordinate, address *string) v1.Coordinate {
	coord := v1.Coordinate{
		Latitude:  float32(kernelCoord.Latitude()),
		Longitude: float32(kernelCoord.Longitude()),
		Address:   address,
	}

	return coord
}
