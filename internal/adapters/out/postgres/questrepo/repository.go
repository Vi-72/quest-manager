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
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Save saves a single quest.
func (r *Repository) Save(ctx context.Context, q quest.Quest) error {
	dto := DomainToDTO(q)
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
}

// GetByID retrieves a quest by its ID.
func (r *Repository) GetByID(ctx context.Context, questID uuid.UUID) (quest.Quest, error) {
	var dto QuestWithAddressDTO
	if err := r.db.WithContext(ctx).
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

	if err := r.db.WithContext(ctx).
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
func (r *Repository) FindByAssignee(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error) {
	var dtos []QuestDTO

	if err := r.db.WithContext(ctx).
		Where("assignee = ?", userID.String()).
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
	if err := r.db.WithContext(ctx).Find(&dtos).Error; err != nil {
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
	if err := r.db.WithContext(ctx).
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
