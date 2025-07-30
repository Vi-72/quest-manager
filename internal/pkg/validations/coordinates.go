package validations

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

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

// validateCoordinate валидирует координаты (используется внутренне)
func validateCoordinate(coord servers.Coordinate, fieldName string) *ValidationError {
	// Валидация latitude
	if coord.Latitude < -90 || coord.Latitude > 90 {
		return NewValidationError(fieldName+".latitude", "must be between -90 and 90")
	}

	// Валидация longitude
	if coord.Longitude < -180 || coord.Longitude > 180 {
		return NewValidationError(fieldName+".longitude", "must be between -180 and 180")
	}

	return nil
}
