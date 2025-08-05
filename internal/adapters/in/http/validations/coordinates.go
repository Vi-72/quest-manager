package validations

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

// ValidatedSearchByRadiusData contains validated data for radius search
type ValidatedSearchByRadiusData struct {
	Center   kernel.GeoCoordinate
	RadiusKm float64
}

// ValidateSearchByRadiusParams validates radius search parameters
func ValidateSearchByRadiusParams(lat, lon, radiusKm float32) (*ValidatedSearchByRadiusData, *ValidationError) {
	// Convert to float64 for domain layer
	latF64 := float64(lat)
	lonF64 := float64(lon)
	radiusF64 := float64(radiusKm)

	// Validate radius (API level validation for business logic)
	if radiusF64 <= 0 {
		return nil, NewValidationError("radius_km", "must be greater than 0 kilometers")
	}
	if radiusF64 > 20000 { // Earth's circumference is ~40000km, so max half of that
		return nil, NewValidationError("radius_km", "must be less than 20000 kilometers")
	}

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
func ConvertAPICoordinateToKernel(apiCoord servers.Coordinate) (kernel.GeoCoordinate, *ValidationError) {
	// Create domain coordinate (convert float32 to float64)
	// Domain layer handles all coordinate validation
	coord, err := kernel.NewGeoCoordinate(float64(apiCoord.Latitude), float64(apiCoord.Longitude))
	if err != nil {
		return kernel.GeoCoordinate{}, NewValidationErrorWithCause("coordinate", "invalid coordinate values", err)
	}

	return coord, nil
}

// ConvertKernelCoordinateToAPI converts domain coordinates to API format
func ConvertKernelCoordinateToAPI(kernelCoord kernel.GeoCoordinate, address *string) servers.Coordinate {
	coord := servers.Coordinate{
		Latitude:  float32(kernelCoord.Latitude()),
		Longitude: float32(kernelCoord.Longitude()),
		Address:   address,
	}

	return coord
}
