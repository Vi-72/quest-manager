package location

import (
	"quest-manager/internal/core/domain/model/kernel"
	"time"

	"github.com/google/uuid"
)

// LocationCreated представляет событие создания локации
type LocationCreated struct {
	ID         uuid.UUID            `json:"id"`       // ID локации
	EventID    uuid.UUID            `json:"event_id"` // ID события
	Name       string               `json:"name"`
	Coordinate kernel.GeoCoordinate `json:"coordinate"`
	Timestamp  time.Time            `json:"timestamp"`
}

func (e LocationCreated) GetID() uuid.UUID {
	return e.EventID
}

func (e LocationCreated) GetName() string {
	return "location.created"
}

// LocationUpdated представляет событие обновления локации
type LocationUpdated struct {
	ID         uuid.UUID            `json:"id"`       // ID локации
	EventID    uuid.UUID            `json:"event_id"` // ID события
	Name       string               `json:"name"`
	Coordinate kernel.GeoCoordinate `json:"coordinate"`
	Timestamp  time.Time            `json:"timestamp"`
}

func (e LocationUpdated) GetID() uuid.UUID {
	return e.EventID
}

func (e LocationUpdated) GetName() string {
	return "location.updated"
}

// LocationDeleted представляет событие удаления локации
type LocationDeleted struct {
	ID        uuid.UUID `json:"id"`       // ID локации
	EventID   uuid.UUID `json:"event_id"` // ID события
	Timestamp time.Time `json:"timestamp"`
}

func (e LocationDeleted) GetID() uuid.UUID {
	return e.EventID
}

func (e LocationDeleted) GetName() string {
	return "location.deleted"
}
