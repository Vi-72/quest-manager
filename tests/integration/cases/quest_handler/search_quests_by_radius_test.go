package quest_handler

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

func (s *Suite) TestSearchQuestsByRadius() {
	ctx := context.Background()

	// Pre-condition - create quest at specific location
	centerLocation := testdatagenerators.DefaultTestCoordinate() // Moscow center: 55.7558, 37.6176
	nearLocation := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.01,  // ~1km away
		Lon: centerLocation.Longitude() + 0.01, // ~1km away
	}
	farLocation := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.1,  // ~10km away
		Lon: centerLocation.Longitude() + 0.1, // ~10km away
	}

	// Create quest near the center (within 5km radius)
	nearQuestData := testdatagenerators.SimpleQuestData(
		"Near Quest", "Quest near center", "easy", 1, 30, nearLocation, nearLocation)
	nearQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, nearQuestData)
	s.Require().NoError(err)

	// Create quest far from center (outside 5km radius)
	farQuestData := testdatagenerators.SimpleQuestData(
		"Far Quest", "Quest far from center", "medium", 2, 60, farLocation, farLocation)
	farQuest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, farQuestData)
	s.Require().NoError(err)

	// Act - search for quests within 5km radius
	radiusKm := 5.0
	foundQuests, err := casesteps.SearchQuestsByRadiusStep(ctx, s.TestDIContainer.SearchQuestsByRadiusHandler,
		centerLocation, radiusKm)

	// Assert
	s.Require().NoError(err)

	// Should find the near quest but not the far quest
	listAssertions := assertions.NewQuestListAssertions(s.Assert())
	listAssertions.QuestWithIDExists(foundQuests, nearQuest.ID().String())
	listAssertions.QuestWithIDNotExists(foundQuests, farQuest.ID().String())
}

func (s *Suite) TestSearchQuestsByRadiusEmpty() {
	ctx := context.Background()

	// Pre-condition - use location far from any existing quests
	remoteLocation := kernel.GeoCoordinate{
		Lat: -89.0, // Near South Pole
		Lon: 0.0,
	}

	// Act - search for quests in remote location
	radiusKm := 10.0
	foundQuests, err := casesteps.SearchQuestsByRadiusStep(ctx, s.TestDIContainer.SearchQuestsByRadiusHandler,
		remoteLocation, radiusKm)

	// Assert
	s.Require().NoError(err)
	s.Assert().Len(foundQuests, 0, "Should return empty list for remote location")
}

func (s *Suite) TestSearchQuestsByRadiusMultipleQuests() {
	ctx := context.Background()

	// Pre-condition - create multiple quests at different distances
	centerLocation := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}

	// Quest 1: Very close (within 1km)
	quest1Location := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.005,  // ~0.5km away
		Lon: centerLocation.Longitude() + 0.005, // ~0.5km away
	}
	quest1Data := testdatagenerators.SimpleQuestData(
		"Quest 1", "Very close quest", "easy", 1, 30, quest1Location, quest1Location)
	quest1, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, quest1Data)
	s.Require().NoError(err)

	// Quest 2: Medium distance (within 5km)
	quest2Location := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.03,  // ~3km away
		Lon: centerLocation.Longitude() + 0.03, // ~3km away
	}
	quest2Data := testdatagenerators.SimpleQuestData(
		"Quest 2", "Medium distance quest", "medium", 2, 60, quest2Location, quest2Location)
	quest2, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, quest2Data)
	s.Require().NoError(err)

	// Quest 3: Far away (outside 5km)
	quest3Location := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.1,  // ~10km away
		Lon: centerLocation.Longitude() + 0.1, // ~10km away
	}
	quest3Data := testdatagenerators.SimpleQuestData(
		"Quest 3", "Far away quest", "hard", 3, 90, quest3Location, quest3Location)
	quest3, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, quest3Data)
	s.Require().NoError(err)

	// Act - search within 5km radius
	radiusKm := 5.0
	foundQuests, err := casesteps.SearchQuestsByRadiusStep(ctx, s.TestDIContainer.SearchQuestsByRadiusHandler,
		centerLocation, radiusKm)

	// Assert
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(len(foundQuests), 2, "Should find at least 2 quests within radius")

	listAssertions := assertions.NewQuestListAssertions(s.Assert())
	listAssertions.QuestWithIDExists(foundQuests, quest1.ID().String())
	listAssertions.QuestWithIDExists(foundQuests, quest2.ID().String())
	listAssertions.QuestWithIDNotExists(foundQuests, quest3.ID().String())
}

func (s *Suite) TestSearchQuestsByRadiusWithTargetAndExecutionLocations() {
	ctx := context.Background()

	// Pre-condition - create quest where target is near but execution is far
	centerLocation := kernel.GeoCoordinate{Lat: 40.0, Lon: 20.0}

	// Target location near center
	targetLocation := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.01,  // ~1km away
		Lon: centerLocation.Longitude() + 0.01, // ~1km away
	}

	// Execution location far from center
	executionLocation := kernel.GeoCoordinate{
		Lat: centerLocation.Latitude() + 0.2,  // ~20km away
		Lon: centerLocation.Longitude() + 0.2, // ~20km away
	}

	questData := testdatagenerators.SimpleQuestData(
		"Mixed Location Quest", "Quest with near target, far execution", "medium", 2, 60, targetLocation, executionLocation)
	quest, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questData)
	s.Require().NoError(err)

	// Act - search within 5km radius (should find quest because target is within radius)
	radiusKm := 5.0
	foundQuests, err := casesteps.SearchQuestsByRadiusStep(ctx, s.TestDIContainer.SearchQuestsByRadiusHandler,
		centerLocation, radiusKm)

	// Assert
	s.Require().NoError(err)

	listAssertions := assertions.NewQuestListAssertions(s.Assert())
	listAssertions.QuestWithIDExists(foundQuests, quest.ID().String())
}
