package authentication

import (
	"context"
	"net/http"

	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
	"quest-manager/tests/integration/mock"
)

// AUTHENTICATION MIDDLEWARE TESTS - EXPIRED TOKEN
// Tests for requests with expired authentication tokens

func (s *Suite) TestCreateQuestWithExpiredToken() {
	ctx := context.Background()

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Prepare valid quest data
	questData := testdatagenerators.SimpleQuestData(
		"Test Quest",
		"Quest with expired token",
		"easy",
		2,
		60,
		testdatagenerators.DefaultTestCoordinate(),
		testdatagenerators.DefaultTestCoordinate(),
	)

	// Act - send request with expired token
	createReq := casesteps.CreateQuestHTTPRequest(questData)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, createReq)

	// Assert - should return 401 Unauthorized with token expired message
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, createResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(createResp.Body, "Token Expired", "Error should mention token expiration")
	s.Assert().Contains(createResp.Body, "expired", "Error should mention expired")
}

func (s *Suite) TestAssignQuestWithExpiredToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with valid auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Act - try to assign with expired token
	assignReq := casesteps.AssignQuestHTTPRequest(createdQuest.ID())
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, assignReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, assignResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(assignResp.Body, "expired", "Error should mention token expiration")
}

func (s *Suite) TestListAssignedQuestsWithExpiredToken() {
	ctx := context.Background()

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Act - try to list assigned quests with expired token
	listReq := casesteps.ListAssignedQuestsHTTPRequest()
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(listResp.Body, "expired", "Error should mention token expiration")
}

func (s *Suite) TestGetQuestWithExpiredToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with valid auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Act - try to get quest with expired token
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.ID())
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, getReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, getResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(getResp.Body, "expired", "Error should mention token expiration")
}

func (s *Suite) TestListQuestsWithExpiredToken() {
	ctx := context.Background()

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Act - try to list quests with expired token
	listReq := casesteps.ListQuestsHTTPRequest("")
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(listResp.Body, "expired", "Error should mention token expiration")
}

func (s *Suite) TestChangeQuestStatusWithExpiredToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with valid auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Act - try to change status with expired token
	statusReq := casesteps.ChangeQuestStatusHTTPRequest(createdQuest.ID(), map[string]string{"status": "posted"})
	statusResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, statusReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, statusResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(statusResp.Body, "expired", "Error should mention token expiration")
}

func (s *Suite) TestSearchQuestsByRadiusWithExpiredToken() {
	ctx := context.Background()

	// Prepare expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Act - try to search quests with expired token
	searchReq := casesteps.SearchQuestsByRadiusHTTPRequest(55.7558, 37.6173, 10)
	searchResp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, searchReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, searchResp.StatusCode, "Should return 401 Unauthorized for expired token")
	s.Assert().Contains(searchResp.Body, "expired", "Error should mention token expiration")
}
