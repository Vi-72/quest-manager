package quest

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListQuests() {
	ctx := context.Background()

	// Arrange - create multiple quests
	expectedCount := 2
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Act - get list of quests (all quests, without status filter)
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, nil)

	// Assert
	s.Require().NoError(err)

	// Create assertions for list verification
	listAssertions := assertions.NewQuestListAssertions(s.Assert())
	listAssertions.QuestsHaveMinimumCount(quests, expectedCount)
	listAssertions.QuestsContainAllCreated(createdQuests, quests)
}

func (s *Suite) TestListQuestsEmpty() {
	ctx := context.Background()

	// Act - get list of quests from empty database
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, nil)

	// Assert
	s.Require().NoError(err)
	s.Assert().Len(quests, 0)
}

func (s *Suite) TestListQuestsWithValidStatus() {
	ctx := context.Background()

	// Arrange - create multiple quests
	expectedCount := 3
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Change quest status for filtering test
	targetStatus := quest.StatusPosted
	_, err = casesteps.ChangeQuestStatusStep(ctx, s.TestDIContainer.ChangeQuestStatusHandler,
		s.TestDIContainer.QuestRepository, createdQuests[0].ID(), targetStatus)
	s.Require().NoError(err)

	// Leave other quests with default StatusCreated status

	// Act - get list of quests filtered by StatusPosted
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, &targetStatus)

	// Assert
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(len(quests), 1, "Should have at least one quest with StatusPosted")

	listAssertions := assertions.NewQuestListAssertions(s.Assert())
	listAssertions.QuestsAllHaveStatus(quests, targetStatus)
	listAssertions.QuestWithIDExists(quests, createdQuests[0].ID().String())
}

func (s *Suite) TestListQuestsWithEmptyStatus() {
	ctx := context.Background()

	// Arrange - create multiple quests
	expectedCount := 2
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Act - get list of quests without status filter (nil)
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, nil)

	// Assert
	s.Require().NoError(err)

	// Create assertions for list verification
	listAssertions := assertions.NewQuestListAssertions(s.Assert())
	listAssertions.QuestsHaveMinimumCount(quests, expectedCount)
	listAssertions.QuestsContainAllCreated(createdQuests, quests)
}

func (s *Suite) TestListQuestsWithInvalidStatus() {
	ctx := context.Background()

	// Arrange - create quest to ensure database is not empty
	_, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - try to get list of quests with invalid status
	invalidStatus := quest.Status("invalid_status_that_does_not_exist")
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, &invalidStatus)

	// Assert - should return validation error
	s.Require().Error(err, "Should return validation error for invalid status")
	s.Assert().Contains(err.Error(), "must be one of", "Error should mention valid status values")
	s.Assert().Nil(quests, "Quests should be nil when there's a validation error")
}
