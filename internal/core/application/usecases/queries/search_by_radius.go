package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// SearchQuestsByRadiusQuery represents the input parameters for searching quests by radius.
type SearchQuestsByRadiusQuery struct {
	Center   kernel.GeoCoordinate
	RadiusKm float64
}

// SearchQuestsByRadiusQueryHandler defines the interface for handling SearchQuestsByRadiusQuery.
type SearchQuestsByRadiusQueryHandler interface {
	Handle(ctx context.Context, query SearchQuestsByRadiusQuery) ([]quest.Quest, error)
}

type searchQuestsByRadiusHandler struct {
	repo ports.QuestRepository
}

// NewSearchQuestsByRadiusQueryHandler creates a new SearchQuestsByRadiusQueryHandler instance.
func NewSearchQuestsByRadiusQueryHandler(repo ports.QuestRepository) SearchQuestsByRadiusQueryHandler {
	return &searchQuestsByRadiusHandler{repo: repo}
}

// Handle retrieves quests within the specified radius from the center coordinate.
func (h *searchQuestsByRadiusHandler) Handle(ctx context.Context, query SearchQuestsByRadiusQuery) ([]quest.Quest, error) {
	return h.repo.FindByLocation(ctx, query.Center, query.RadiusKm)
}
