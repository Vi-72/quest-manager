package ports

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"

	"github.com/google/uuid"
)

// LocationRepository defines access methods for locations.
type LocationRepository interface {
	GetByID(ctx context.Context, locationID uuid.UUID) (*location.Location, error)
	Save(ctx context.Context, l *location.Location) error

	// FindAll retrieves all locations without filters.
	FindAll(ctx context.Context) ([]*location.Location, error)

	// FindByBoundingBox returns all locations within the specified bounding box area.
	FindByBoundingBox(ctx context.Context, bbox kernel.BoundingBox) ([]*location.Location, error)

	// FindByName searches locations by name (partial match).
	FindByName(ctx context.Context, namePattern string) ([]*location.Location, error)
}
