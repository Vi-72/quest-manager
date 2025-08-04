package repository

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

type QuestRepositoryTestSuite struct {
	RepositoryTestSuite
}

func (suite *QuestRepositoryTestSuite) SetupTest() {
	// Call parent setup to clean database
	suite.RepositoryTestSuite.SetupTest()
}

func (suite *QuestRepositoryTestSuite) TestSave_Success() {
	ctx := context.Background()

	// Pre-condition - create a valid quest
	q := suite.createTestQuest("Test Quest", "easy")

	// Act - save quest
	err := suite.questRepo.Save(ctx, q)

	// Assert
	suite.Require().NoError(err)

	// Verify quest was saved by retrieving it
	saved, err := suite.questRepo.GetByID(ctx, q.ID())
	suite.Require().NoError(err)
	suite.assertQuestEquals(q, saved)
}

func (suite *QuestRepositoryTestSuite) TestSave_Update() {
	ctx := context.Background()

	// Pre-condition - save initial quest
	q := suite.createTestQuest("Original Title", "easy")
	err := suite.questRepo.Save(ctx, q)
	suite.Require().NoError(err)

	// Modify quest
	q.Title = "Updated Title"
	q.Description = "Updated Description"

	// Act - save updated quest
	err = suite.questRepo.Save(ctx, q)

	// Assert
	suite.Require().NoError(err)

	// Verify quest was updated
	updated, err := suite.questRepo.GetByID(ctx, q.ID())
	suite.Require().NoError(err)
	suite.Equal("Updated Title", updated.Title)
	suite.Equal("Updated Description", updated.Description)
}

func (suite *QuestRepositoryTestSuite) TestGetByID_Success() {
	ctx := context.Background()

	// Pre-condition - save a quest
	q := suite.createTestQuest("Test Quest", "medium")
	err := suite.questRepo.Save(ctx, q)
	suite.Require().NoError(err)

	// Act - get quest by ID
	found, err := suite.questRepo.GetByID(ctx, q.ID())

	// Assert
	suite.Require().NoError(err)
	suite.assertQuestEquals(q, found)
}

func (suite *QuestRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()

	// Act - try to get non-existent quest
	nonExistentID := uuid.New()
	_, err := suite.questRepo.GetByID(ctx, nonExistentID)

	// Assert
	suite.Require().Error(err)
	suite.Contains(err.Error(), "not found")
}

func (suite *QuestRepositoryTestSuite) TestFindAll_Success() {
	ctx := context.Background()

	// Pre-condition - save multiple quests
	quest1 := suite.createTestQuest("Quest 1", "easy")
	quest2 := suite.createTestQuest("Quest 2", "hard")

	err := suite.questRepo.Save(ctx, quest1)
	suite.Require().NoError(err)
	err = suite.questRepo.Save(ctx, quest2)
	suite.Require().NoError(err)

	// Act - find all quests
	quests, err := suite.questRepo.FindAll(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.Len(quests, 2)

	// Verify both quests are present
	questIDs := make(map[string]bool)
	for _, q := range quests {
		questIDs[q.ID().String()] = true
	}
	suite.True(questIDs[quest1.ID().String()])
	suite.True(questIDs[quest2.ID().String()])
}

func (suite *QuestRepositoryTestSuite) TestFindAll_Empty() {
	ctx := context.Background()

	// Act - find all quests in empty database
	quests, err := suite.questRepo.FindAll(ctx)

	// Assert
	suite.Require().NoError(err)
	suite.Len(quests, 0)
}

func (suite *QuestRepositoryTestSuite) TestFindByStatus_Success() {
	ctx := context.Background()

	// Pre-condition - save quests with different statuses
	quest1 := suite.createTestQuest("Created Quest", "easy")
	quest2 := suite.createTestQuest("Posted Quest", "medium")

	// Change status of second quest
	err := quest2.ChangeStatus(quest.StatusPosted)
	suite.Require().NoError(err)

	err = suite.questRepo.Save(ctx, quest1)
	suite.Require().NoError(err)
	err = suite.questRepo.Save(ctx, quest2)
	suite.Require().NoError(err)

	// Act - find quests by status
	createdQuests, err := suite.questRepo.FindByStatus(ctx, quest.StatusCreated)
	suite.Require().NoError(err)

	postedQuests, err := suite.questRepo.FindByStatus(ctx, quest.StatusPosted)
	suite.Require().NoError(err)

	// Assert
	suite.Len(createdQuests, 1)
	suite.Equal(quest1.ID(), createdQuests[0].ID())

	suite.Len(postedQuests, 1)
	suite.Equal(quest2.ID(), postedQuests[0].ID())
}

func (suite *QuestRepositoryTestSuite) TestFindByAssignee_Success() {
	ctx := context.Background()

	// Pre-condition - create quests and assign them
	quest1 := suite.createTestQuest("Quest 1", "easy")
	quest2 := suite.createTestQuest("Quest 2", "medium")
	quest3 := suite.createTestQuest("Quest 3", "hard")

	userID1 := "user-123"
	userID2 := "user-456"

	err := quest1.AssignTo(userID1)
	suite.Require().NoError(err)
	err = quest2.AssignTo(userID1)
	suite.Require().NoError(err)
	err = quest3.AssignTo(userID2)
	suite.Require().NoError(err)

	err = suite.questRepo.Save(ctx, quest1)
	suite.Require().NoError(err)
	err = suite.questRepo.Save(ctx, quest2)
	suite.Require().NoError(err)
	err = suite.questRepo.Save(ctx, quest3)
	suite.Require().NoError(err)

	// Act - find quests by assignee
	user1Quests, err := suite.questRepo.FindByAssignee(ctx, userID1)
	suite.Require().NoError(err)

	user2Quests, err := suite.questRepo.FindByAssignee(ctx, userID2)
	suite.Require().NoError(err)

	// Assert
	suite.Len(user1Quests, 2)
	suite.Len(user2Quests, 1)

	// Verify correct assignment
	user1QuestIDs := make(map[string]bool)
	for _, q := range user1Quests {
		user1QuestIDs[q.ID().String()] = true
	}
	suite.True(user1QuestIDs[quest1.ID().String()])
	suite.True(user1QuestIDs[quest2.ID().String()])

	suite.Equal(quest3.ID(), user2Quests[0].ID())
}

func (suite *QuestRepositoryTestSuite) TestFindByBoundingBox_Success() {
	ctx := context.Background()

	// Pre-condition - create quests at different locations
	moscowCenter := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	moscowSuburb := kernel.GeoCoordinate{Lat: 55.8000, Lon: 37.7000}
	stPetersburg := kernel.GeoCoordinate{Lat: 59.9311, Lon: 30.3609}

	quest1 := suite.createTestQuestAtLocation("Moscow Center Quest", moscowCenter)
	quest2 := suite.createTestQuestAtLocation("Moscow Suburb Quest", moscowSuburb)
	quest3 := suite.createTestQuestAtLocation("St Petersburg Quest", stPetersburg)

	err := suite.questRepo.Save(ctx, quest1)
	suite.Require().NoError(err)
	err = suite.questRepo.Save(ctx, quest2)
	suite.Require().NoError(err)
	err = suite.questRepo.Save(ctx, quest3)
	suite.Require().NoError(err)

	// Act - find quests in Moscow area bounding box
	moscowBBox := kernel.BoundingBox{
		MinLat: 55.5, MaxLat: 56.0,
		MinLon: 37.0, MaxLon: 38.0,
	}
	moscowQuests, err := suite.questRepo.FindByBoundingBox(ctx, moscowBBox)

	// Assert
	suite.Require().NoError(err)
	suite.Len(moscowQuests, 2) // Should find both Moscow quests, but not St Petersburg

	moscowQuestIDs := make(map[string]bool)
	for _, q := range moscowQuests {
		moscowQuestIDs[q.ID().String()] = true
	}
	suite.True(moscowQuestIDs[quest1.ID().String()])
	suite.True(moscowQuestIDs[quest2.ID().String()])
}

// Helper methods
func (suite *QuestRepositoryTestSuite) createTestQuest(title, difficulty string) quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	q, err := quest.NewQuest(
		title,
		"Test description for "+title,
		difficulty,
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"test-equipment"},
		[]string{"test-skill"},
	)
	suite.Require().NoError(err)
	return q
}

func (suite *QuestRepositoryTestSuite) createTestQuestAtLocation(title string, location kernel.GeoCoordinate) quest.Quest {
	q, err := quest.NewQuest(
		title,
		"Test description for "+title,
		"medium",
		3,
		60,
		location,
		location, // Same location for simplicity
		"test-creator",
		[]string{"test-equipment"},
		[]string{"test-skill"},
	)
	suite.Require().NoError(err)
	return q
}

func (suite *QuestRepositoryTestSuite) assertQuestEquals(expected, actual quest.Quest) {
	suite.Equal(expected.ID(), actual.ID())
	suite.Equal(expected.Title, actual.Title)
	suite.Equal(expected.Description, actual.Description)
	suite.Equal(expected.Difficulty, actual.Difficulty)
	suite.Equal(expected.Reward, actual.Reward)
	suite.Equal(expected.DurationMinutes, actual.DurationMinutes)
	suite.Equal(expected.Status, actual.Status)
	suite.Equal(expected.Creator, actual.Creator)
	suite.Equal(expected.TargetLocation, actual.TargetLocation)
	suite.Equal(expected.ExecutionLocation, actual.ExecutionLocation)

	if expected.Assignee == nil {
		suite.Nil(actual.Assignee)
	} else {
		suite.NotNil(actual.Assignee)
		suite.Equal(*expected.Assignee, *actual.Assignee)
	}
}
