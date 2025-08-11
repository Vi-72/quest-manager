package questrepo

import (
	"context"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		if err := r.tracker.Begin(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to begin quest transaction", err)
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
		return errs.WrapInfrastructureError("failed to save quest", err)
	}

	if !isInTransaction {
		if err := r.tracker.Commit(ctx); err != nil {
			return errs.WrapInfrastructureError("failed to commit quest transaction", err)
		}
	}
	return nil
}

// GetByID retrieves a quest by its ID.
func (r *Repository) GetByID(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	var dto QuestWithAddressDTO
	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Select("quests.*, target_loc.address as target_address, exec_loc.address as execution_address").
		Table("quests").
		Joins("LEFT JOIN locations target_loc ON quests.target_location_id = target_loc.id").
		Joins("LEFT JOIN locations exec_loc ON quests.execution_location_id = exec_loc.id").
		Where("quests.id = ?", questID.String()).
		First(&dto).Error; err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to get quest by ID", err)
	}

	return DtoToDomainWithAddress(dto)
}

// FindByBoundingBox retrieves quests within a bounding box area.
// Simple database query without business logic.
func (r *Repository) FindByBoundingBox(ctx context.Context, bbox kernel.BoundingBox) ([]quest.Quest, error) {
	var dtos []QuestDTO

	db := r.tracker.Db()
	if err := db.WithContext(ctx).
		Where("(target_latitude BETWEEN ? AND ? AND target_longitude BETWEEN ? AND ?) OR "+
			"(execution_latitude BETWEEN ? AND ? AND execution_longitude BETWEEN ? AND ?)",
			bbox.MinLat, bbox.MaxLat, bbox.MinLon, bbox.MaxLon,
			bbox.MinLat, bbox.MaxLat, bbox.MinLon, bbox.MaxLon).
		Find(&dtos).Error; err != nil {
		return nil, errs.WrapInfrastructureError("failed to get quests by bounding box", err)
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
