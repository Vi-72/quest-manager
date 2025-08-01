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
	Address    *string              `json:"address,omitempty"`
}

func NewLocationCreated(locationID uuid.UUID, coordinate kernel.GeoCoordinate, address *string) LocationCreated {
	return LocationCreated{
		BaseEvent:  ddd.NewBaseEvent(locationID, "location.created"),
		Coordinate: coordinate,
		Address:    address,
	}
}

// LocationUpdated represents location update event
type LocationUpdated struct {
	ddd.BaseEvent
	Coordinate kernel.GeoCoordinate `json:"coordinate"`
	Address    *string              `json:"address,omitempty"`
}

func NewLocationUpdated(locationID uuid.UUID, coordinate kernel.GeoCoordinate, address *string) LocationUpdated {
	return LocationUpdated{
		BaseEvent:  ddd.NewBaseEvent(locationID, "location.updated"),
		Coordinate: coordinate,
		Address:    address,
	}
}
