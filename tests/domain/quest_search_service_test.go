package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/domain/services"
)

func TestQuestSearchService_FilterByRadius(t *testing.T) {
	searchService := services.NewQuestSearchService()

	// Test data setup
	center := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176} // Moscow coordinates
	radiusKm := 10.0

	// Create test quests with different locations
	quest1, err1 := createTestQuestWithLocation("Quest 1", 55.7558, 37.6176) // Same as center
	if err1 != nil {
		t.Fatalf("Failed to create quest1: %v", err1)
	}
	quest2, err2 := createTestQuestWithLocation("Quest 2", 55.7658, 37.6276) // ~1.4km from center
	if err2 != nil {
		t.Fatalf("Failed to create quest2: %v", err2)
	}
	quest3, err3 := createTestQuestWithLocation("Quest 3", 55.8558, 37.7176) // ~15km from center (outside radius)
	if err3 != nil {
		t.Fatalf("Failed to create quest3: %v", err3)
	}
	quest4, err4 := createTestQuestWithLocation("Quest 4", 55.7458, 37.6076) // ~1.4km from center
	if err4 != nil {
		t.Fatalf("Failed to create quest4: %v", err4)
	}

	// Quest with execution location within radius (target outside)
	quest5, err5 := createTestQuestWithExecutionLocation("Quest 5", 55.7558, 37.6176, 55.8558, 37.7176)
	if err5 != nil {
		t.Fatalf("Failed to create quest5: %v", err5)
	}

	allQuests := []quest.Quest{quest1, quest2, quest3, quest4, quest5}

	t.Run("should filter quests within radius", func(t *testing.T) {
		result := searchService.FilterByRadius(allQuests, center, radiusKm)

		// Should include quests 1, 2, 4 (target location within radius)
		// and quest 5 (execution location within radius)
		assert.Len(t, result, 4, "Should return 4 quests within radius")

		questTitles := make([]string, len(result))
		for i, q := range result {
			questTitles[i] = q.Title
		}

		assert.Contains(t, questTitles, "Quest 1", "Should include Quest 1")
		assert.Contains(t, questTitles, "Quest 2", "Should include Quest 2")
		assert.Contains(t, questTitles, "Quest 4", "Should include Quest 4")
		assert.Contains(t, questTitles, "Quest 5", "Should include Quest 5")
		assert.NotContains(t, questTitles, "Quest 3", "Should not include Quest 3 (outside radius)")
	})

	t.Run("should return empty slice for no matching quests", func(t *testing.T) {
		// All quests outside radius
		farQuests := []quest.Quest{quest3}
		result := searchService.FilterByRadius(farQuests, center, 1.0) // Very small radius

		assert.Empty(t, result, "Should return empty slice when no quests within radius")
	})

	t.Run("should handle empty input", func(t *testing.T) {
		result := searchService.FilterByRadius([]quest.Quest{}, center, radiusKm)

		assert.Empty(t, result, "Should return empty slice for empty input")
	})

	t.Run("should handle zero radius", func(t *testing.T) {
		// Create a quest list with only Quest 1 at exact center
		centerQuests := []quest.Quest{quest1}
		result := searchService.FilterByRadius(centerQuests, center, 0.0)

		// Only quest at exact center should match
		assert.Len(t, result, 1, "Should return only quest at exact center")
		assert.Equal(t, "Quest 1", result[0].Title, "Should return Quest 1 at exact center")
	})

	t.Run("should handle very large radius", func(t *testing.T) {
		result := searchService.FilterByRadius(allQuests, center, 1000.0) // Very large radius

		assert.Len(t, result, len(allQuests), "Should return all quests for very large radius")
	})

	t.Run("should consider both target and execution locations", func(t *testing.T) {
		// Quest with target outside but execution inside radius
		questTargetOutside, _ := createTestQuestWithExecutionLocation("Target Outside", 55.8558, 37.7176, 55.7558, 37.6176)
		quests := []quest.Quest{questTargetOutside}

		result := searchService.FilterByRadius(quests, center, radiusKm)

		assert.Len(t, result, 1, "Should include quest with execution location within radius")
		assert.Equal(t, "Target Outside", result[0].Title)
	})
}

func TestQuestSearchService_EdgeCases(t *testing.T) {
	searchService := services.NewQuestSearchService()
	center := kernel.GeoCoordinate{Lat: 0.0, Lon: 0.0} // Equator, Prime Meridian

	t.Run("should handle quests at exact same coordinates", func(t *testing.T) {
		quest1, _ := createTestQuestWithLocation("Quest 1", 0.0, 0.0)
		quest2, _ := createTestQuestWithLocation("Quest 2", 0.0, 0.0)
		quests := []quest.Quest{quest1, quest2}

		result := searchService.FilterByRadius(quests, center, 1.0)

		assert.Len(t, result, 2, "Should return both quests at exact same coordinates")
	})

	t.Run("should handle negative coordinates", func(t *testing.T) {
		center := kernel.GeoCoordinate{Lat: -55.7558, Lon: -37.6176} // Negative coordinates
		quest1, _ := createTestQuestWithLocation("Quest 1", -55.7558, -37.6176)
		quest2, _ := createTestQuestWithLocation("Quest 2", -55.7658, -37.6276)
		quests := []quest.Quest{quest1, quest2}

		result := searchService.FilterByRadius(quests, center, 10.0)

		assert.Len(t, result, 2, "Should handle negative coordinates correctly")
	})

	t.Run("should handle extreme coordinates", func(t *testing.T) {
		// North Pole
		center := kernel.GeoCoordinate{Lat: 90.0, Lon: 0.0}
		quest1, _ := createTestQuestWithLocation("Quest 1", 89.9, 0.0)
		quests := []quest.Quest{quest1}

		result := searchService.FilterByRadius(quests, center, 100.0)

		assert.Len(t, result, 1, "Should handle extreme coordinates correctly")
	})
}

// Helper functions to create test quests

func createTestQuestWithLocation(title string, lat, lon float64) (quest.Quest, error) {
	// Use direct struct creation like in other tests
	location := kernel.GeoCoordinate{Lat: lat, Lon: lon}

	q, err := quest.NewQuest(
		title,
		"Test Description",
		"easy",
		3, // Valid reward (1-5)
		60,
		location,
		location,
		"test-creator",
		[]string{"test"},
		[]string{"test"},
	)

	if err != nil {
		return quest.Quest{}, err
	}

	return q, nil
}

func createTestQuestWithExecutionLocation(title string, targetLat, targetLon, execLat, execLon float64) (quest.Quest, error) {
	// Use direct struct creation like in other tests
	targetLocation := kernel.GeoCoordinate{Lat: targetLat, Lon: targetLon}
	executionLocation := kernel.GeoCoordinate{Lat: execLat, Lon: execLon}

	return quest.NewQuest(
		title,
		"Test Description",
		"easy",
		3, // Valid reward (1-5)
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"test"},
		[]string{"test"},
	)
}
