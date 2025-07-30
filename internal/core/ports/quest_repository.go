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

	// FindByLocation returns all quests by target or execution location.
	FindByLocation(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error)

	// FindByAssignee returns all quests assigned to a specific user.
	FindByAssignee(ctx context.Context, userID string) ([]quest.Quest, error)
}
