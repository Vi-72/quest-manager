//go:build integration

package repository

import (
	"context"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

func (s *Suite) TestQuestRepository_Save_Success() {
	ctx := context.Background()

	// Pre-condition - create a valid quest
	q := s.createTestQuest("Test Quest", "easy")

	// Act - save quest
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)

	// Assert
	s.Require().NoError(err)

	// Verify quest was saved by retrieving it
	saved, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)
	s.assertQuestEquals(q, saved)
}

func (s *Suite) TestQuestRepository_Save_Update() {
	ctx := context.Background()

	// Pre-condition - save initial quest
	q := s.createTestQuest("Original Title", "easy")
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	// Modify quest (simulate domain operations)
	originalID := q.ID()
	// Update the existing quest instead of creating a new one
	q.Title = "Updated Title"
	q.Description = "Updated Description"
	q.Difficulty = "medium"

	// Act - save updated quest
	err = s.TestDIContainer.QuestRepository.Save(ctx, q)

	// Assert
	s.Require().NoError(err)

	// Verify quest was updated
	updated, err := s.TestDIContainer.QuestRepository.GetByID(ctx, originalID)
	s.Require().NoError(err)
	s.Equal("Updated Title", updated.Title)
	s.Equal("Updated Description", updated.Description)
	s.Equal(quest.Difficulty("medium"), updated.Difficulty)
}

func (s *Suite) TestQuestRepository_GetByID_Success() {
	ctx := context.Background()

	// Pre-condition - save a quest
	q := s.createTestQuest("Test Quest", "medium")
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	// Act - get quest by ID
	found, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())

	// Assert
	s.Require().NoError(err)
	s.assertQuestEquals(q, found)
}

func (s *Suite) TestQuestRepository_GetByID_NotFound() {
	ctx := context.Background()

	// Pre-condition - use non-existent ID
	nonExistentID := uuid.New()

	// Act - try to get quest by non-existent ID
	_, err := s.TestDIContainer.QuestRepository.GetByID(ctx, nonExistentID)

	// Assert - should return error
	s.Require().Error(err)
	s.Contains(err.Error(), "not found")
}

func (s *Suite) TestQuestRepository_FindAll_Empty() {
	ctx := context.Background()

	// Act - find all quests when database is empty
	quests, err := s.TestDIContainer.QuestRepository.FindAll(ctx)

	// Assert
	s.Require().NoError(err)
	s.Empty(quests)
}

func (s *Suite) TestQuestRepository_FindAll_Success() {
	ctx := context.Background()

	// Pre-condition - save multiple quests
	quest1 := s.createTestQuest("Quest 1", "easy")
	quest2 := s.createTestQuest("Quest 2", "medium")
	quest3 := s.createTestQuest("Quest 3", "hard")

	err := s.TestDIContainer.QuestRepository.Save(ctx, quest1)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest2)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest3)
	s.Require().NoError(err)

	// Act - find all quests
	quests, err := s.TestDIContainer.QuestRepository.FindAll(ctx)

	// Assert
	s.Require().NoError(err)
	s.Len(quests, 3)

	// Verify all quests are present
	questIDs := make(map[uuid.UUID]bool)
	for _, q := range quests {
		questIDs[q.ID()] = true
	}
	s.True(questIDs[quest1.ID()])
	s.True(questIDs[quest2.ID()])
	s.True(questIDs[quest3.ID()])
}

func (s *Suite) TestQuestRepository_FindByStatus_Success() {
	ctx := context.Background()

	// Pre-condition - save quests with different statuses
	quest1 := s.createTestQuest("Quest 1", "easy")
	quest2 := s.createTestQuest("Quest 2", "medium")
	quest3 := s.createTestQuest("Quest 3", "hard")

	err := s.TestDIContainer.QuestRepository.Save(ctx, quest1)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest2)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest3)
	s.Require().NoError(err)

	// Change status of quest2
	quest2.ChangeStatus(quest.StatusPosted)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest2)
	s.Require().NoError(err)

	// Act - find quests by status
	createdQuests, err := s.TestDIContainer.QuestRepository.FindByStatus(ctx, quest.StatusCreated)
	s.Require().NoError(err)

	postedQuests, err := s.TestDIContainer.QuestRepository.FindByStatus(ctx, quest.StatusPosted)
	s.Require().NoError(err)

	// Assert
	s.Len(createdQuests, 2) // quest1 and quest3
	s.Len(postedQuests, 1)  // quest2

	s.Equal(quest.StatusPosted, postedQuests[0].Status)
}

func (s *Suite) TestQuestRepository_FindByBoundingBox_Success() {
	ctx := context.Background()

	// Pre-condition - create quests at different locations
	// Moscow center
	quest1 := s.createTestQuestAtLocation("Moscow Quest", "easy",
		kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6173})

	// Saint Petersburg
	quest2 := s.createTestQuestAtLocation("SPB Quest", "medium",
		kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609})

	// Another Moscow location
	quest3 := s.createTestQuestAtLocation("Moscow Quest 2", "hard",
		kernel.GeoCoordinate{Lat: 55.7500, Lon: 37.6200})

	err := s.TestDIContainer.QuestRepository.Save(ctx, quest1)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest2)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest3)
	s.Require().NoError(err)

	// Act - find quests in Moscow area (bounding box)
	moscowBoundingBox := kernel.BoundingBox{
		MinLat: 55.7000, // South boundary
		MaxLat: 55.8000, // North boundary
		MinLon: 37.5000, // West boundary
		MaxLon: 37.7000, // East boundary
	}
	moscowQuests, err := s.TestDIContainer.QuestRepository.FindByBoundingBox(ctx, moscowBoundingBox)

	// Assert
	s.Require().NoError(err)
	s.Len(moscowQuests, 2) // Should find both Moscow quests

	foundIDs := make(map[uuid.UUID]bool)
	for _, q := range moscowQuests {
		foundIDs[q.ID()] = true
	}
	s.True(foundIDs[quest1.ID()])
	s.True(foundIDs[quest3.ID()])
	s.False(foundIDs[quest2.ID()]) // SPB quest should not be found
}

func (s *Suite) TestQuestRepository_FindByAssignee_Success() {
	ctx := context.Background()

	// Pre-condition - create and assign quests
	quest1 := s.createTestQuest("Quest 1", "easy")
	quest2 := s.createTestQuest("Quest 2", "medium")
	quest3 := s.createTestQuest("Quest 3", "hard")

	// Post quests first (required for assignment)
	quest1.ChangeStatus(quest.StatusPosted)
	quest2.ChangeStatus(quest.StatusPosted)
	quest3.ChangeStatus(quest.StatusPosted)

	// Assign quests
	err := quest1.AssignTo("user1")
	s.Require().NoError(err)
	err = quest2.AssignTo("user1")
	s.Require().NoError(err)
	err = quest3.AssignTo("user2")
	s.Require().NoError(err)

	// Save assigned quests
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest1)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest2)
	s.Require().NoError(err)
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest3)
	s.Require().NoError(err)

	// Act - find quests assigned to user1
	user1Quests, err := s.TestDIContainer.QuestRepository.FindByAssignee(ctx, "user1")
	s.Require().NoError(err)

	user2Quests, err := s.TestDIContainer.QuestRepository.FindByAssignee(ctx, "user2")
	s.Require().NoError(err)

	// Assert
	s.Len(user1Quests, 2) // quest1 and quest2
	s.Len(user2Quests, 1) // quest3

	for _, q := range user1Quests {
		s.Equal("user1", *q.Assignee)
		s.Equal(quest.StatusAssigned, q.Status)
	}

	s.Equal("user2", *user2Quests[0].Assignee)
}

// ==========================================
// POSTGRESQL-SPECIFIC TESTS
// ==========================================

func (s *Suite) TestQuestRepository_PostgreSQL_JSONArrays() {
	ctx := context.Background()

	// Test PostgreSQL JSON array handling
	q := s.createTestQuest("PostgreSQL Test", "hard")

	// Set complex arrays (PostgreSQL should handle these as JSON)
	q.Equipment = []string{"GPS device", "Waterproof camera", "Climbing gear", "Emergency beacon"}
	q.Skills = []string{"Mountain climbing", "Photography", "Navigation", "First aid"}

	// Save and retrieve
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	found, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)

	// Verify JSON arrays are preserved correctly in PostgreSQL
	s.Equal(q.Equipment, found.Equipment)
	s.Equal(q.Skills, found.Skills)
	s.Len(found.Equipment, 4)
	s.Len(found.Skills, 4)
}

func (s *Suite) TestQuestRepository_PostgreSQL_Transactions() {
	ctx := context.Background()

	// Test that changes within transaction are isolated
	quest1 := s.createTestQuest("Transaction Test 1", "easy")
	quest2 := s.createTestQuest("Transaction Test 2", "medium")

	// Save quest1
	err := s.TestDIContainer.QuestRepository.Save(ctx, quest1)
	s.Require().NoError(err)

	// Within current transaction, save quest2
	err = s.TestDIContainer.QuestRepository.Save(ctx, quest2)
	s.Require().NoError(err)

	// Both should be findable within this transaction
	quests, err := s.TestDIContainer.QuestRepository.FindAll(ctx)
	s.Require().NoError(err)
	s.Len(quests, 2)

	// But transaction will be rolled back in TearDownTest
	// so in next test they won't be there (tested implicitly)
}

func (s *Suite) TestQuestRepository_PostgreSQL_EmptyEquipmentAndSkills() {
	ctx := context.Background()

	// Test PostgreSQL handling of empty Equipment and Skills arrays
	q := s.createTestQuestWithEmptyArrays("Empty Arrays Test", "medium")

	// Save and retrieve
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	found, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)

	// Verify empty arrays are preserved correctly in PostgreSQL
	s.Equal(q.Equipment, found.Equipment)
	s.Equal(q.Skills, found.Skills)
	s.Len(found.Equipment, 0)
	s.Len(found.Skills, 0)
	s.NotNil(found.Equipment) // Теперь всегда [], а не nil
	s.NotNil(found.Skills)    // Теперь всегда [], а не nil
}

func (s *Suite) TestQuestRepository_PostgreSQL_EmptyEquipmentOnly() {
	ctx := context.Background()

	// Test quest with empty Equipment but filled Skills
	q := s.createTestQuestWithEmptyEquipment("Empty Equipment Test", "hard")

	// Save and retrieve
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	found, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)

	// Verify mixed empty/filled arrays
	s.Len(found.Equipment, 0)
	s.Len(found.Skills, 2)
	s.Equal([]string{}, found.Equipment)
	s.Equal([]string{"Navigation", "Survival"}, found.Skills)
}

func (s *Suite) TestQuestRepository_PostgreSQL_EmptySkillsOnly() {
	ctx := context.Background()

	// Test quest with empty Skills but filled Equipment
	q := s.createTestQuestWithEmptySkills("Empty Skills Test", "easy")

	// Save and retrieve
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	found, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)

	// Verify mixed filled/empty arrays
	s.Len(found.Equipment, 2)
	s.Len(found.Skills, 0)
	s.Equal([]string{"Map", "Compass"}, found.Equipment)
	s.True(len(found.Skills) == 0, "Skills should be empty")
}

func (s *Suite) TestQuestRepository_PostgreSQL_ArrayUpdates() {
	ctx := context.Background()

	// Test updating quest from filled arrays to empty arrays
	q := s.createTestQuest("Array Update Test", "medium")

	// Initially save with filled arrays
	err := s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	// Verify initial state
	found, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)
	s.Len(found.Equipment, 1)
	s.Len(found.Skills, 1)

	// Update to empty arrays
	q.Equipment = []string{}
	q.Skills = []string{}
	err = s.TestDIContainer.QuestRepository.Save(ctx, q)
	s.Require().NoError(err)

	// Verify arrays are now empty
	updated, err := s.TestDIContainer.QuestRepository.GetByID(ctx, q.ID())
	s.Require().NoError(err)
	s.Len(updated.Equipment, 0)
	s.Len(updated.Skills, 0)
	s.True(len(updated.Equipment) == 0, "Equipment should be empty")
	s.True(len(updated.Skills) == 0, "Skills should be empty")
}

// ==========================================
// HELPER METHODS
// ==========================================

func (s *Suite) createTestQuest(title, difficulty string) quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6173}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7560, Lon: 37.6175}

	q, err := quest.NewQuest(
		title,
		"Test Description for "+title,
		difficulty,
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"equipment"},
		[]string{"skill"},
	)
	s.Require().NoError(err)
	return q
}

func (s *Suite) createTestQuestAtLocation(title, difficulty string, location kernel.GeoCoordinate) quest.Quest {
	// Same location for target and execution for simplicity
	q, err := quest.NewQuest(
		title,
		"Test Description for "+title,
		difficulty,
		3,
		60,
		location,
		location,
		"test-creator",
		[]string{"equipment"},
		[]string{"skill"},
	)
	s.Require().NoError(err)
	return q
}

func (s *Suite) createTestQuestWithEmptyArrays(title, difficulty string) quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6173}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7560, Lon: 37.6175}

	q, err := quest.NewQuest(
		title,
		"Test Description for "+title,
		difficulty,
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{}, // Empty equipment
		[]string{}, // Empty skills
	)
	s.Require().NoError(err)
	return q
}

func (s *Suite) createTestQuestWithEmptyEquipment(title, difficulty string) quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6173}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7560, Lon: 37.6175}

	q, err := quest.NewQuest(
		title,
		"Test Description for "+title,
		difficulty,
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{},                         // Empty equipment
		[]string{"Navigation", "Survival"}, // Filled skills
	)
	s.Require().NoError(err)
	return q
}

func (s *Suite) createTestQuestWithEmptySkills(title, difficulty string) quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6173}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7560, Lon: 37.6175}

	q, err := quest.NewQuest(
		title,
		"Test Description for "+title,
		difficulty,
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"Map", "Compass"}, // Filled equipment
		[]string{},                 // Empty skills
	)
	s.Require().NoError(err)
	return q
}

func (s *Suite) assertQuestEquals(expected, actual quest.Quest) {
	s.Equal(expected.ID(), actual.ID())
	s.Equal(expected.Title, actual.Title)
	s.Equal(expected.Description, actual.Description)
	s.Equal(expected.Difficulty, actual.Difficulty)
	s.Equal(expected.Reward, actual.Reward)
	s.Equal(expected.DurationMinutes, actual.DurationMinutes)
	s.Equal(expected.TargetLocation, actual.TargetLocation)
	s.Equal(expected.ExecutionLocation, actual.ExecutionLocation)
	s.Equal(expected.Creator, actual.Creator)
	s.Equal(expected.Equipment, actual.Equipment)
	s.Equal(expected.Skills, actual.Skills)
	s.Equal(expected.Status, actual.Status)
}
