package ports

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

// QuestRepository defines access methods for quests.
type QuestRepository interface {
	GetByID(ctx context.Context, questID uuid.UUID) (quest.Quest, error)
	Save(ctx context.Context, q quest.Quest) error

	// FindAll retrieves all quests without filters.
	FindAll(ctx context.Context) ([]quest.Quest, error)

	// FindByStatus retrieves all quests with the specified status.
	FindByStatus(ctx context.Context, status quest.Status) ([]quest.Quest, error)

	// FindByBoundingBox returns all quests within the specified bounding box area.
	// This is a simple database query without business logic.
	FindByBoundingBox(ctx context.Context, bbox kernel.BoundingBox) ([]quest.Quest, error)

	// FindByAssignee returns all quests assigned to a specific user.
	FindByAssignee(ctx context.Context, userID uuid.UUID) ([]quest.Quest, error)
}
