package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// SearchQuestsByRadiusQuery represents the input parameters for searching quests within a radius.
type SearchQuestsByRadiusQuery struct {
	Center   kernel.GeoCoordinate
	RadiusKm float64
}

// SearchQuestsByRadiusResult represents the result of searching quests by radius.
type SearchQuestsByRadiusResult struct {
	Quests []quest.Quest
}

// SearchQuestsByRadiusQueryHandler defines the interface for handling SearchQuestsByRadiusQuery.
type SearchQuestsByRadiusQueryHandler interface {
	Handle(ctx context.Context, query SearchQuestsByRadiusQuery) (SearchQuestsByRadiusResult, error)
}

type searchQuestsByRadiusHandler struct {
	repo ports.QuestRepository
}

// NewSearchQuestsByRadiusQueryHandler creates a new instance of SearchQuestsByRadiusQueryHandler.
func NewSearchQuestsByRadiusQueryHandler(repo ports.QuestRepository) SearchQuestsByRadiusQueryHandler {
	return &searchQuestsByRadiusHandler{repo: repo}
}

// Handle executes the query to search for quests within a given radius.
func (h *searchQuestsByRadiusHandler) Handle(ctx context.Context, query SearchQuestsByRadiusQuery) (SearchQuestsByRadiusResult, error) {
	quests, err := h.repo.FindByLocation(ctx, query.Center, query.RadiusKm)
	if err != nil {
		return SearchQuestsByRadiusResult{}, err
	}
	return SearchQuestsByRadiusResult{Quests: quests}, nil
}