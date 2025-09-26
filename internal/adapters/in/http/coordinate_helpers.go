package http

import (
	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/core/domain/model/kernel"
)

func convertAPICoordinateToKernel(coord v1.Coordinate) (kernel.GeoCoordinate, error) {
	return kernel.NewGeoCoordinate(float64(coord.Latitude), float64(coord.Longitude))
}

func convertKernelCoordinateToAPI(coord kernel.GeoCoordinate, address *string) v1.Coordinate {
	return v1.Coordinate{
		Latitude:  float32(coord.Latitude()),
		Longitude: float32(coord.Longitude()),
		Address:   address,
	}
}
