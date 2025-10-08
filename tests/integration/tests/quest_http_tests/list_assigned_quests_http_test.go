package quest_http_tests

// API LAYER VALIDATION TESTS
// Focused on OpenAPI-driven request validation behavior

import (
	"context"

	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	"quest-manager/tests/integration/mock"
)

func (s *Suite) TestListAssignedQuestsHTTP() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	listAssertions := assertions.NewQuestListAssertions(s.Assert())

	// Pre-condition - create quests via handler and assign them to authenticated user
	// User ID comes from mock auth client (DefaultUserID)
	testUserUUID := mock.NewAlwaysSuccessAuthClient().DefaultUserID
	expectedCount := 2

	// Create quests via handler
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)

	// Assign all created quests to the test user
	for _, quest := range createdQuests {
		_, err := casesteps.AssignQuestStep(ctx, s.TestDIContainer.AssignQuestHandler, quest.ID(), testUserUUID)
		s.Require().NoError(err)
	}

	// Act - get list of assigned quests via HTTP API (user ID comes from JWT token)
	listReq := casesteps.ListAssignedQuestsHTTPRequest()
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert using helpers to eliminate boilerplate
	assignedQuests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)
	listAssertions.QuestListHTTPAllAssignedToUser(assignedQuests, testUserUUID, expectedCount)
	listAssertions.QuestListHTTPContainsAllCreated(assignedQuests, createdQuests)
}

func (s *Suite) TestListAssignedQuestsHTTPEmpty() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - no quests assigned to the authenticated user (DefaultUserID)
	// Since we're using a fresh test environment, there should be no assigned quests

	// Act - get list of assigned quests for authenticated user via HTTP API
	listReq := casesteps.ListAssignedQuestsHTTPRequest()
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert using helper
	assignedQuests := httpAssertions.QuestHTTPListSuccessfully(listResp, err)
	s.Assert().Len(assignedQuests, 0, "Should return empty list for user with no assigned quests")
}
