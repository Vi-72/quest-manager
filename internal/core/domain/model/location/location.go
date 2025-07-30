package location

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/pkg/ddd"
	"time"

	"github.com/google/uuid"
)

// Location represents a named geographic location that can be reused across quests
type Location struct {
	*ddd.BaseAggregate[uuid.UUID]
	Name        string
	Coordinate  kernel.GeoCoordinate
	Address     string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewLocation creates a new location with validation
func NewLocation(name string, coordinate kernel.GeoCoordinate, address, description string) (*Location, error) {
	id := uuid.New()
	now := time.Now()

	location := &Location{
		BaseAggregate: ddd.NewBaseAggregate(id),
		Name:          name,
		Coordinate:    coordinate,
		Address:       address,
		Description:   description,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Raise domain event
	location.RaiseDomainEvent(LocationCreated{
		ID:         id,
		EventID:    uuid.New(),
		Name:       name,
		Coordinate: coordinate,
		Timestamp:  now,
	})

	return location, nil
}

// Update updates location information
func (l *Location) Update(name string, coordinate kernel.GeoCoordinate, address, description string) error {
	l.Name = name
	l.Coordinate = coordinate
	l.Address = address
	l.Description = description
	l.UpdatedAt = time.Now()

	// Raise domain event
	l.RaiseDomainEvent(LocationUpdated{
		ID:         l.ID(),
		EventID:    uuid.New(),
		Name:       name,
		Coordinate: coordinate,
		Timestamp:  l.UpdatedAt,
	})

	return nil
}
