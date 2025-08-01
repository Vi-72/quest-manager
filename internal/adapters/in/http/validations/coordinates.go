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
	// Validate latitude
	latF64 := float64(lat)
	if latF64 < -90 || latF64 > 90 {
		return nil, NewValidationError("lat", "must be between -90 and 90 degrees")
	}

	// Validate longitude
	lonF64 := float64(lon)
	if lonF64 < -180 || lonF64 > 180 {
		return nil, NewValidationError("lon", "must be between -180 and 180 degrees")
	}

	// Validate radius
	radiusF64 := float64(radiusKm)
	if radiusF64 <= 0 {
		return nil, NewValidationError("radius_km", "must be greater than 0 kilometers")
	}
	if radiusF64 > 20000 { // Earth's circumference is ~40000km, so max half of that
		return nil, NewValidationError("radius_km", "must be less than 20000 kilometers")
	}

	// Create domain coordinate
	center, err := kernel.NewGeoCoordinate(latF64, lonF64, nil)
	if err != nil {
		return nil, NewValidationErrorWithCause("coordinates", "invalid coordinate values", err)
	}

	return &ValidatedSearchByRadiusData{
		Center:   center,
		RadiusKm: radiusF64,
	}, nil
}

// ConvertAPICoordinateToKernel converts API coordinates to domain coordinates with validation
func ConvertAPICoordinateToKernel(apiCoord servers.Coordinate) (kernel.GeoCoordinate, *ValidationError) {
	// Validate latitude
	if apiCoord.Latitude < -90 || apiCoord.Latitude > 90 {
		return kernel.GeoCoordinate{}, NewValidationError("latitude", "must be between -90 and 90")
	}

	// Validate longitude
	if apiCoord.Longitude < -180 || apiCoord.Longitude > 180 {
		return kernel.GeoCoordinate{}, NewValidationError("longitude", "must be between -180 and 180")
	}

	// Create domain coordinate (convert float32 to float64)
	coord, err := kernel.NewGeoCoordinate(float64(apiCoord.Latitude), float64(apiCoord.Longitude), apiCoord.Address)
	if err != nil {
		return kernel.GeoCoordinate{}, NewValidationErrorWithCause("coordinate", "invalid coordinate values", err)
	}

	return coord, nil
}

// ConvertKernelCoordinateToAPI converts domain coordinates to API format
func ConvertKernelCoordinateToAPI(kernelCoord kernel.GeoCoordinate) servers.Coordinate {
	coord := servers.Coordinate{
		Latitude:  float32(kernelCoord.Latitude()),
		Longitude: float32(kernelCoord.Longitude()),
		Address:   kernelCoord.GetAddress(),
	}

	return coord
}
