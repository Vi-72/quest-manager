package location

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/pkg/ddd"
	"time"

	"github.com/google/uuid"
)

// Location represents a geographic location that can be reused across quests
type Location struct {
	*ddd.BaseAggregate[uuid.UUID]
	Coordinate kernel.GeoCoordinate
	Address    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewLocation creates a new location with validation
func NewLocation(coordinate kernel.GeoCoordinate, address string) (*Location, error) {
	id := uuid.New()
	now := time.Now()

	location := &Location{
		BaseAggregate: ddd.NewBaseAggregate(id),
		Coordinate:    coordinate,
		Address:       address,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Raise domain event
	location.RaiseDomainEvent(LocationCreated{
		ID: uuid.New(),
		Coordinate: LocationCoordinate{
			GeoCoordinate: coordinate,
			LocationID:    id,
		},
		Timestamp: now,
	})

	return location, nil
}

// Update updates location information
func (l *Location) Update(coordinate kernel.GeoCoordinate, address string) error {
	l.Coordinate = coordinate
	l.Address = address
	l.UpdatedAt = time.Now()

	// Raise domain event
	l.RaiseDomainEvent(LocationUpdated{
		ID: uuid.New(),
		Coordinate: LocationCoordinate{
			GeoCoordinate: coordinate,
			LocationID:    l.ID(),
		},
		Timestamp: l.UpdatedAt,
	})

	return nil
}
