package casesteps

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"

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

// ListAssignedQuestsStep gets list of quests assigned to a user
func ListAssignedQuestsStep(
	ctx context.Context,
	handler queries.ListAssignedQuestsQueryHandler,
	userID string,
) ([]quest.Quest, error) {
	return handler.Handle(ctx, userID)
}

// SearchQuestsByRadiusStep searches for quests within a radius from center coordinates
func SearchQuestsByRadiusStep(
	ctx context.Context,
	handler queries.SearchQuestsByRadiusQueryHandler,
	center kernel.GeoCoordinate,
	radiusKm float64,
) ([]quest.Quest, error) {
	return handler.Handle(ctx, center, radiusKm)
}
