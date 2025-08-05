package quest_handler

// HANDLER LAYER INTEGRATION TESTS
// Tests for assignQuestHandler.Handle orchestration logic

import (
	"context"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/quest"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestAssignQuest() {
	ctx := context.Background()

	// Pre-condition - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - assign quest with valid data
	userID := "test-user-123"
	assignResult, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, createdQuest.ID(), userID)

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(createdQuest.ID(), assignResult.ID)
	s.Assert().Equal(userID, assignResult.Assignee)
	s.Assert().Equal(string(quest.StatusAssigned), assignResult.Status)
}

func (s *Suite) TestAssignQuestFromPostedStatus() {
	ctx := context.Background()

	// Pre-condition - create quest and change to posted status
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Change quest status to posted
	postedQuest, err := casesteps.ChangeQuestStatusStep(
		ctx,
		s.TestDIContainer.ChangeQuestStatusHandler,
		s.TestDIContainer.QuestRepository,
		createdQuest.ID(),
		quest.StatusPosted,
	)
	s.Require().NoError(err)
	s.Assert().Equal(quest.StatusPosted, postedQuest.Status)

	// Act - assign quest from posted status
	userID := "test-user-456"
	assignResult, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, postedQuest.ID(), userID)

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(postedQuest.ID(), assignResult.ID)
	s.Assert().Equal(userID, assignResult.Assignee)
	s.Assert().Equal(string(quest.StatusAssigned), assignResult.Status)
}

func (s *Suite) TestAssignQuestNotFound() {
	ctx := context.Background()

	// Act - try to assign non-existent quest (handler should return 404 error)
	nonExistentQuestID := uuid.New() // Generate random UUID
	userID := "test-user-123"
	_, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, nonExistentQuestID, userID)

	// Assert - handler should return not found error
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "quest")
	s.Assert().Contains(err.Error(), "not found")
}

func (s *Suite) TestAssignQuestInvalidStatus() {
	ctx := context.Background()

	// Pre-condition - create quest, assign it, then change to in_progress status
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// First assign quest to a user
	firstUserID := "first-user"
	_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, createdQuest.ID(), firstUserID)
	s.Require().NoError(err)

	// Change quest status to in_progress (invalid for assignment)
	inProgressQuest, err := casesteps.ChangeQuestStatusStep(
		ctx,
		s.TestDIContainer.ChangeQuestStatusHandler,
		s.TestDIContainer.QuestRepository,
		createdQuest.ID(),
		quest.StatusInProgress,
	)
	s.Require().NoError(err)
	s.Assert().Equal(quest.StatusInProgress, inProgressQuest.Status)

	// Act - try to assign quest with invalid status (domain validation error → 400)
	secondUserID := "test-user-789"
	_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, inProgressQuest.ID(), secondUserID)

	// Assert - handler should return domain validation error
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "assignment")
	s.Assert().Contains(err.Error(), "failed to assign quest")
}

func (s *Suite) TestAssignQuestAlreadyAssigned() {
	ctx := context.Background()

	// Pre-condition - create and assign quest to first user
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	firstUserID := "first-user-123"
	_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, createdQuest.ID(), firstUserID)
	s.Require().NoError(err)

	// Act - try to assign already assigned quest to second user (domain validation error → 400)
	secondUserID := "second-user-456"
	_, err = casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, createdQuest.ID(), secondUserID)

	// Assert - handler should return domain validation error
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "assignment")
	s.Assert().Contains(err.Error(), "failed to assign quest")
}

func (s *Suite) TestAssignQuestPersistence() {
	ctx := context.Background()

	// Pre-condition - create quest
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - assign quest
	userID := "persistent-user-123"
	assignResult, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, createdQuest.ID(), userID)
	s.Require().NoError(err)

	// Verify quest is persisted by retrieving it
	foundQuest, err := casesteps.GetQuestByIDStep(ctx, s.TestDIContainer.GetQuestByIDHandler, assignResult.ID)

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(assignResult.ID, foundQuest.ID())
	s.Assert().NotNil(foundQuest.Assignee, "Quest should have assignee")
	s.Assert().Equal(userID, *foundQuest.Assignee)
	s.Assert().Equal(quest.StatusAssigned, foundQuest.Status)

	// Verify timestamps are updated
	s.Assert().False(foundQuest.CreatedAt.IsZero(), "CreatedAt should be set")
	s.Assert().False(foundQuest.UpdatedAt.IsZero(), "UpdatedAt should be set")
	s.Assert().True(foundQuest.UpdatedAt.After(foundQuest.CreatedAt), "UpdatedAt should be after CreatedAt")
}

func (s *Suite) TestAssignQuestWithSameUserMultipleTimes() {
	ctx := context.Background()

	// Pre-condition - create multiple random quests
	firstQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	secondQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - assign both quests to same user
	userID := "multi-quest-user-123"

	firstAssignResult, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, firstQuest.ID(), userID)
	s.Require().NoError(err)

	secondAssignResult, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, secondQuest.ID(), userID)
	s.Require().NoError(err)

	// Assert - both assignments should succeed
	s.Assert().Equal(firstQuest.ID(), firstAssignResult.ID)
	s.Assert().Equal(userID, firstAssignResult.Assignee)
	s.Assert().Equal(string(quest.StatusAssigned), firstAssignResult.Status)

	s.Assert().Equal(secondQuest.ID(), secondAssignResult.ID)
	s.Assert().Equal(userID, secondAssignResult.Assignee)
	s.Assert().Equal(string(quest.StatusAssigned), secondAssignResult.Status)

	// Verify both quests are assigned to the same user
	s.Assert().NotEqual(firstAssignResult.ID, secondAssignResult.ID, "Quest IDs should be different")
}
