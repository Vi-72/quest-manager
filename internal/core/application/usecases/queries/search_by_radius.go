package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// SearchQuestsByRadiusQueryHandler defines the interface for handling quest search by radius.
type SearchQuestsByRadiusQueryHandler interface {
	Handle(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error)
}

type searchQuestsByRadiusHandler struct {
	unitOfWorkFactory ports.UnitOfWorkFactory
}

// NewSearchQuestsByRadiusQueryHandler creates a new SearchQuestsByRadiusQueryHandler instance.
func NewSearchQuestsByRadiusQueryHandler(factory ports.UnitOfWorkFactory) SearchQuestsByRadiusQueryHandler {
	return &searchQuestsByRadiusHandler{unitOfWorkFactory: factory}
}

// Handle retrieves quests within the specified radius from the center coordinate.
// This method contains the business logic for geospatial search:
// 1. Calculate bounding box for efficient database query
// 2. Filter results by exact radius using Haversine distance
func (h *searchQuestsByRadiusHandler) Handle(ctx context.Context, center kernel.GeoCoordinate, radiusKm float64) ([]quest.Quest, error) {
	// Step 1: Calculate bounding box using domain model
	bbox := center.BoundingBoxForRadius(radiusKm)

	// Step 2: Get candidates from repository using simple bounding box query
	unitOfWork, _, err := h.unitOfWorkFactory()
	if err != nil {
		return nil, errs.WrapInfrastructureError("failed to create unit of work for radius search", err)
	}

	candidates, err := unitOfWork.QuestRepository().FindByBoundingBox(ctx, bbox)
	if err != nil {
		return nil, err
	}

	// Step 3: Apply business logic - filter by exact radius using accurate Haversine distance
	var result []quest.Quest
	for _, q := range candidates {
		// Check if either target location OR execution location is within radius
		if center.DistanceTo(q.TargetLocation) <= radiusKm ||
			center.DistanceTo(q.ExecutionLocation) <= radiusKm {
			result = append(result, q)
		}
	}

	return result, nil
}
