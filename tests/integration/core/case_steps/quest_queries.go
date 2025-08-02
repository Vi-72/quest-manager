package casesteps

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/quest"
)

// GetQuestByIDStep gets quest by ID
func GetQuestByIDStep(
	ctx context.Context,
	handler queries.GetQuestByIDQueryHandler,
	questID uuid.UUID,
) (quest.Quest, error) {
	return handler.Handle(ctx, questID)
}

// ListQuestsStep gets list of quests
func ListQuestsStep(
	ctx context.Context,
	handler queries.ListQuestsQueryHandler,
	status *quest.Status,
) ([]quest.Quest, error) {
	return handler.Handle(ctx, status)
}
