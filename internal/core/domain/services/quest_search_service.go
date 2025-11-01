package services

import (
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

// QuestSearchService encapsulates domain logic for quest search operations.
// This is a Domain Service because the logic doesn't naturally belong to
// any single aggregate.
type QuestSearchService interface {
	// FilterByRadius filters quests to only those within the specified radius
	// from the center coordinate. Uses Haversine formula for accurate distance.
	FilterByRadius(quests []quest.Quest, center kernel.GeoCoordinate, radiusKm float64) []quest.Quest
}

type questSearchService struct{}

func NewQuestSearchService() QuestSearchService {
	return &questSearchService{}
}

func (s *questSearchService) FilterByRadius(quests []quest.Quest, center kernel.GeoCoordinate, radiusKm float64) []quest.Quest {
	result := make([]quest.Quest, 0, len(quests))

	for _, q := range quests {
		// Business rule: Quest is within radius if EITHER target OR execution location matches
		if center.DistanceTo(q.TargetLocation) <= radiusKm ||
			center.DistanceTo(q.ExecutionLocation) <= radiusKm {
			result = append(result, q)
		}
	}

	return result
}
