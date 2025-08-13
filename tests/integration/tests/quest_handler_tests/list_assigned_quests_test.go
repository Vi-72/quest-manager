package quest_handler_tests

// HANDLER LAYER INTEGRATION TESTS
// Tests for listAssignedQuestsHandler.Handle orchestration logic

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListAssignedQuests() {
	ctx := context.Background()
	listAssertions := assertions.NewQuestListAssertions(s.Assert())

	// Pre-condition - create multiple quests and assign them to a specific user
	testUserID := "test-user-123"
	expectedCount := 2
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Assign all created quests to the test user
	for _, q := range createdQuests {
		_, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, q.ID(), testUserID)
		s.Require().NoError(err)
	}

	// Act - get list of quests assigned to the user
	quests, err := casesteps.ListAssignedQuestsStep(ctx, s.TestDIContainer.ListAssignedQuestsHandler, testUserID)

	// Assert
	s.Require().NoError(err)
	listAssertions.QuestsHaveMinimumCount(quests, expectedCount)
	listAssertions.QuestsContainAllCreated(createdQuests, quests)

	// Verify all returned quests are assigned to the correct user
	for _, q := range quests {
		s.Assert().Equal(testUserID, *q.Assignee, "Quest should be assigned to the test user")
		s.Assert().Equal(quest.StatusAssigned, q.Status, "Quest should have 'assigned' status")
	}
}

func (s *Suite) TestListAssignedQuestsEmpty() {
	ctx := context.Background()

	// Pre-condition - use a user ID that has no assigned quests
	nonExistentUserID := "user-with-no-quests"

	// Act - get list of assigned quests for user with no assignments
	quests, err := casesteps.ListAssignedQuestsStep(ctx, s.TestDIContainer.ListAssignedQuestsHandler, nonExistentUserID)

	// Assert
	s.Require().NoError(err)
	s.Assert().Len(quests, 0, "Should return empty list for user with no assigned quests")
}
