package validations

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

// ValidatedSearchByRadiusData содержит валидированные данные для поиска по радиусу
type ValidatedSearchByRadiusData struct {
	Center   kernel.GeoCoordinate
	RadiusKm float64
}

// ValidateSearchByRadiusParams валидирует параметры поиска по радиусу
func ValidateSearchByRadiusParams(lat, lon, radiusKm float32) (*ValidatedSearchByRadiusData, *ValidationError) {
	// Валидация широты
	latF64 := float64(lat)
	if latF64 < -90 || latF64 > 90 {
		return nil, NewValidationError("lat", "must be between -90 and 90 degrees")
	}

	// Валидация долготы
	lonF64 := float64(lon)
	if lonF64 < -180 || lonF64 > 180 {
		return nil, NewValidationError("lon", "must be between -180 and 180 degrees")
	}

	// Валидация радиуса
	radiusF64 := float64(radiusKm)
	if radiusF64 <= 0 {
		return nil, NewValidationError("radius_km", "must be greater than 0 kilometers")
	}
	if radiusF64 > 20000 { // Earth's circumference is ~40000km, so max half of that
		return nil, NewValidationError("radius_km", "must be less than 20000 kilometers")
	}

	// Создаем доменную координату
	center, err := kernel.NewGeoCoordinate(latF64, lonF64)
	if err != nil {
		return nil, NewValidationErrorWithCause("coordinates", "invalid coordinate values", err)
	}

	return &ValidatedSearchByRadiusData{
		Center:   center,
		RadiusKm: radiusF64,
	}, nil
}

// ConvertAPICoordinateToKernel конвертирует API координаты в доменные с валидацией
func ConvertAPICoordinateToKernel(apiCoord servers.Coordinate) (kernel.GeoCoordinate, *ValidationError) {
	// Валидация широты
	if apiCoord.Latitude < -90 || apiCoord.Latitude > 90 {
		return kernel.GeoCoordinate{}, NewValidationError("latitude", "must be between -90 and 90")
	}

	// Валидация долготы
	if apiCoord.Longitude < -180 || apiCoord.Longitude > 180 {
		return kernel.GeoCoordinate{}, NewValidationError("longitude", "must be between -180 and 180")
	}

	// Создаем доменную координату (конвертируем float32 в float64)
	coord, err := kernel.NewGeoCoordinate(float64(apiCoord.Latitude), float64(apiCoord.Longitude))
	if err != nil {
		return kernel.GeoCoordinate{}, NewValidationErrorWithCause("coordinate", "invalid coordinate values", err)
	}

	return coord, nil
}

// ConvertKernelCoordinateToAPI конвертирует доменные координаты в API формат
func ConvertKernelCoordinateToAPI(kernelCoord kernel.GeoCoordinate) servers.Coordinate {
	return servers.Coordinate{
		Latitude:  float32(kernelCoord.Latitude()),
		Longitude: float32(kernelCoord.Longitude()),
	}
}
