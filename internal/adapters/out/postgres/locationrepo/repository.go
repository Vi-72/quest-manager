package locationrepo

import (
	"context"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/location"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ports.LocationRepository = &Repository{}

type Repository struct {
	tracker ports.Tracker
}

func NewRepository(tracker ports.Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	return &Repository{tracker: tracker}, nil
}

// Save saves a single location.
func (r *Repository) Save(ctx context.Context, l *location.Location) error {
	dto := DomainToDTO(l)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		if err := r.tracker.Begin(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to begin location transaction", err)
		}
	}
	tx := r.tracker.Tx()

	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		if !isInTransaction {
			if rollbackErr := r.tracker.Rollback(); rollbackErr != nil {
				// Log rollback error but don't override the original error
				_ = rollbackErr
			}
		}
		return errs.WrapInfrastructureError("failed to save location", err)
	}

	if !isInTransaction {
		if err := r.tracker.Commit(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to commit location transaction", err)
		}
	}
	return nil
}

// GetByID retrieves a location by its ID.
func (r *Repository) GetByID(ctx context.Context, locationID uuid.UUID) (*location.Location, error) {
	var dto LocationDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("id = ?", locationID.String()).
		First(&dto).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get location by ID", err)
	}
	return DtoToDomain(dto)
}

// FindAll retrieves all locations without filters.
func (r *Repository) FindAll(ctx context.Context) ([]*location.Location, error) {
	var dtos []LocationDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get all locations", err)
	}

	locations := make([]*location.Location, len(dtos))
	for i, dto := range dtos {
		l, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}
		locations[i] = l
	}

	return locations, nil
}

// FindByBoundingBox retrieves locations within a bounding box area.
func (r *Repository) FindByBoundingBox(ctx context.Context, bbox kernel.BoundingBox) ([]*location.Location, error) {
	var dtos []LocationDTO

	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
			bbox.MinLat, bbox.MaxLat, bbox.MinLon, bbox.MaxLon).
		Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get locations by bounding box", err)
	}

	locations := make([]*location.Location, len(dtos))
	for i, dto := range dtos {
		l, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}
		locations[i] = l
	}

	return locations, nil
}

// FindByName searches locations by name (partial match).
func (r *Repository) FindByName(ctx context.Context, namePattern string) ([]*location.Location, error) {
	var dtos []LocationDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("name ILIKE ?", "%"+namePattern+"%").
		Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get locations by name", err)
	}

	locations := make([]*location.Location, len(dtos))
	for i, dto := range dtos {
		l, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}
		locations[i] = l
	}

	return locations, nil
}
