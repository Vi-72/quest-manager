package location

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/pkg/ddd"

	"github.com/google/uuid"
)

// LocationCreated represents location creation event
type LocationCreated struct {
	ddd.BaseEvent
	Coordinate kernel.GeoCoordinate `json:"coordinate"`
}

func NewLocationCreated(locationID uuid.UUID, coordinate kernel.GeoCoordinate, address string) LocationCreated {
	return LocationCreated{
		BaseEvent:  ddd.NewBaseEvent(locationID, "location.created"),
		Coordinate: coordinate,
	}
}

// LocationUpdated represents location update event
type LocationUpdated struct {
	ddd.BaseEvent
	Coordinate kernel.GeoCoordinate `json:"coordinate"`
}

func NewLocationUpdated(locationID uuid.UUID, coordinate kernel.GeoCoordinate, address string) LocationUpdated {
	return LocationUpdated{
		BaseEvent:  ddd.NewBaseEvent(locationID, "location.updated"),
		Coordinate: coordinate,
	}
}
