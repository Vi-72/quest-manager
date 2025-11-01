package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/domain/services"
	"quest-manager/internal/core/ports"
)

// SearchQuestsByRadiusQueryHandler defines the interface for handling quest search by radius.
type SearchQuestsByRadiusQueryHandler interface {
	Handle(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error)
}

type searchQuestsByRadiusHandler struct {
	questRepo     ports.QuestRepository
	searchService services.QuestSearchService
}

// NewSearchQuestsByRadiusQueryHandler creates a new SearchQuestsByRadiusQueryHandler instance.
func NewSearchQuestsByRadiusQueryHandler(questRepo ports.QuestRepository) SearchQuestsByRadiusQueryHandler {
	return &searchQuestsByRadiusHandler{
		questRepo:     questRepo,
		searchService: services.NewQuestSearchService(),
	}
}

// Handle retrieves quests within the specified radius from the center coordinate.
// This method orchestrates the geospatial search:
// 1. Calculate bounding box for efficient database query
// 2. Get candidates from database (infrastructure concern)
// 3. Apply domain service for precise filtering
func (h *searchQuestsByRadiusHandler) Handle(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error) {
	// Step 1: Calculate bounding box using domain model
	bbox := center.BoundingBoxForRadius(radiusKm)

	// Step 2: Get candidates from repository using simple bounding box query
	candidates, err := h.questRepo.FindByBoundingBox(ctx, bbox)
	if err != nil {
		return nil, err
	}

	// Step 3: Apply domain service for precise filtering
	return h.searchService.FilterByRadius(candidates, center, radiusKm), nil
}
