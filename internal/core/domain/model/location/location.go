package location

import (
	"time"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/timeprovider"

	"github.com/google/uuid"
)

// Location represents a geographic location that can be reused across quests
type Location struct {
	*ddd.BaseAggregate[uuid.UUID]
	Coordinate   kernel.GeoCoordinate
	Address      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TimeProvider timeprovider.TimeProvider
}

// NewLocation creates a new location with validation
func NewLocation(coordinate kernel.GeoCoordinate, address *string) (*Location, error) {
	id := uuid.New()
	tp := timeprovider.RealTimeProvider{}
	now := tp.Now()

	location := &Location{
		BaseAggregate: ddd.NewBaseAggregate(id),
		Coordinate:    coordinate,
		Address:       address,
		CreatedAt:     now,
		UpdatedAt:     now,
		TimeProvider:  tp,
	}

	// Raise domain event
	location.RaiseDomainEvent(NewLocationCreated(id, coordinate, address))

	return location, nil
}

// Update updates location information
func (l *Location) Update(coordinate kernel.GeoCoordinate, address *string) error {
	l.Coordinate = coordinate
	l.Address = address
	l.UpdatedAt = l.TimeProvider.Now()

	// Raise domain event
	l.RaiseDomainEvent(NewLocationUpdated(l.ID(), coordinate, address))

	return nil
}

// SetTimeProvider overrides the time provider used by the location (primarily for tests).
func (l *Location) SetTimeProvider(tp timeprovider.TimeProvider) {
	l.TimeProvider = tp
}
