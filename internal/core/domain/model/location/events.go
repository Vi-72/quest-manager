package location

import (
	"quest-manager/internal/core/domain/model/kernel"
	"time"

	"github.com/google/uuid"
)

// LocationCoordinate представляет координаты с ID локации для событий
type LocationCoordinate struct {
	kernel.GeoCoordinate
	LocationID uuid.UUID `json:"location_id"` // ID локации
}

// LocationCreated представляет событие создания локации
type LocationCreated struct {
	ID         uuid.UUID          `json:"id"` // ID события
	Coordinate LocationCoordinate `json:"coordinate"`
	Timestamp  time.Time          `json:"timestamp"`
}

func (e LocationCreated) GetID() uuid.UUID {
	return e.ID
}

func (e LocationCreated) GetName() string {
	return "location.created"
}

// LocationUpdated представляет событие обновления локации
type LocationUpdated struct {
	ID         uuid.UUID          `json:"id"` // ID события
	Coordinate LocationCoordinate `json:"coordinate"`
	Timestamp  time.Time          `json:"timestamp"`
}

func (e LocationUpdated) GetID() uuid.UUID {
	return e.ID
}

func (e LocationUpdated) GetName() string {
	return "location.updated"
}
