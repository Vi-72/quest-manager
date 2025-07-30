package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// SearchQuestsByRadiusQueryHandler defines the interface for handling quest search by radius.
type SearchQuestsByRadiusQueryHandler interface {
	Handle(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error)
}

type searchQuestsByRadiusHandler struct {
	repo ports.QuestRepository
}

// NewSearchQuestsByRadiusQueryHandler creates a new SearchQuestsByRadiusQueryHandler instance.
func NewSearchQuestsByRadiusQueryHandler(repo ports.QuestRepository) SearchQuestsByRadiusQueryHandler {
	return &searchQuestsByRadiusHandler{repo: repo}
}

// Handle retrieves quests within the specified radius from the center coordinate.
func (h *searchQuestsByRadiusHandler) Handle(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error) {
	return h.repo.FindByLocation(ctx, center, radiusKm)
}
