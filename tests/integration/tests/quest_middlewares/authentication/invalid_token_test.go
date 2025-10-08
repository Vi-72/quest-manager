package authentication

import (
	"context"
	"net/http"

	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
	"quest-manager/tests/integration/mock"
)

// AUTHENTICATION MIDDLEWARE TESTS - INVALID TOKEN
// Tests for requests with invalid/malformed authentication tokens

func (s *Suite) TestCreateQuestWithInvalidToken() {
	ctx := context.Background()

	// Prepare invalid token auth client
	invalidTokenAuthClient := mock.NewInvalidTokenAuthClient()
	routerWithInvalidAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(invalidTokenAuthClient)

	// Prepare valid quest data
	questData := testdatagenerators.SimpleQuestData(
		"Test Quest",
		"Quest with invalid token",
		"easy",
		2,
		60,
		testdatagenerators.DefaultTestCoordinate(),
		testdatagenerators.DefaultTestCoordinate(),
	)

	// Act - send request with invalid token
	createReq := casesteps.CreateQuestHTTPRequest(questData)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithInvalidAuth, createReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, createResp.StatusCode, "Should return 401 Unauthorized for invalid token")
	s.Assert().Contains(createResp.Body, "Authentication Failed", "Error should mention authentication failure")
}

func (s *Suite) TestAssignQuestWithInvalidToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with valid auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Prepare invalid token auth client
	invalidTokenAuthClient := mock.NewInvalidTokenAuthClient()
	routerWithInvalidAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(invalidTokenAuthClient)

	// Act - try to assign with invalid token
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithInvalidAuth, assignReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, assignResp.StatusCode, "Should return 401 Unauthorized for invalid token")
}

func (s *Suite) TestListAssignedQuestsWithInvalidToken() {
	ctx := context.Background()

	// Prepare invalid token auth client
	invalidTokenAuthClient := mock.NewInvalidTokenAuthClient()
	routerWithInvalidAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(invalidTokenAuthClient)

	// Act - try to list assigned quests with invalid token
	listReq := casesteps.ListAssignedQuestsHTTPRequest()
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithInvalidAuth, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 Unauthorized for invalid token")
}
