package questrepo

import (
	"context"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
)

var _ ports.QuestRepository = &Repository{}

type Repository struct {
	tracker ports.Tracker
}

func NewRepository(tracker ports.Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	return &Repository{tracker: tracker}, nil
}

// Save saves a single quest.
func (r *Repository) Save(ctx context.Context, q quest.Quest) error {
	dto := DomainToDTO(q)
	if err := r.tracker.Db().WithContext(ctx).Save(&dto).Error; err != nil {
		return errs.WrapInfrastructureError("failed to save quest", err)
	}
	return nil
}

// GetByID retrieves a quest by its ID.
func (r *Repository) GetByID(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	var dto QuestDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("id = ?", questID.String()).
		First(&dto).Error; err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to get quest by ID", err)
	}
	return DtoToDomain(dto)
}

const earthRadiusKm = 6371.0

// FindByLocation retrieves quests within a radius (in km) around the given coordinate.
func (r *Repository) FindByLocation(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error) {
	var dtos []QuestDTO

	// Bounding box in degrees (1 degree ~ 111 km)
	radiusDeg := radiusKm / 111.0
	minLat := center.Latitude() - radiusDeg
	maxLat := center.Latitude() + radiusDeg
	minLon := center.Longitude() - radiusDeg
	maxLon := center.Longitude() + radiusDeg

	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("(target_latitude BETWEEN ? AND ? AND target_longitude BETWEEN ? AND ?) OR "+
			"(execution_latitude BETWEEN ? AND ? AND execution_longitude BETWEEN ? AND ?)",
			minLat, maxLat, minLon, maxLon,
			minLat, maxLat, minLon, maxLon).
		Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get quests by location", err)
	}

	// Filter by exact radius
	var result []quest.Quest
	for _, dto := range dtos {
		q, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}

		if center.DistanceTo(q.TargetLocation) <= radiusKm ||
			center.DistanceTo(q.ExecutionLocation) <= radiusKm {
			result = append(result, q)
		}
	}

	return result, nil
}

// FindByAssignee retrieves all quests assigned to a specific user.
func (r *Repository) FindByAssignee(ctx context.Context, userID string) ([]quest.Quest, error) {
	var dtos []QuestDTO

	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("assignee = ?", userID).
		Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get quests by assignee", err)
	}

	quests := make([]quest.Quest, len(dtos))
	for i, dto := range dtos {
		q, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}
		quests[i] = q
	}

	return quests, nil
}

// FindAll retrieves all quests without any filter.
func (r *Repository) FindAll(ctx context.Context) ([]quest.Quest, error) {
	var dtos []QuestDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get all quests", err)
	}

	quests := make([]quest.Quest, len(dtos))
	for i, dto := range dtos {
		q, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}
		quests[i] = q
	}

	return quests, nil
}

// FindByStatus retrieves all quests with the specified status.
func (r *Repository) FindByStatus(ctx context.Context, status quest.Status) ([]quest.Quest, error) {
	var dtos []QuestDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("status = ?", string(status)).
		Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get quests by status", err)
	}

	quests := make([]quest.Quest, len(dtos))
	for i, dto := range dtos {
		q, err := DtoToDomain(dto)
		if err != nil {
			return nil, errs.WrapInfrastructureError("failed to convert dto to domain", err)
		}
		quests[i] = q
	}

	return quests, nil
}
