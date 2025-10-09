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

func (s *Suite) TestAllEndpointsRejectExpiredToken() {
	ctx := context.Background()

	// Pre-condition - create quest for tests that need it
	createdQuest, _ := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)

	// Setup - create router with expired token auth client
	expiredTokenAuthClient := mock.NewExpiredTokenAuthClient()
	routerWithExpiredAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(expiredTokenAuthClient)

	// Test cases covering all endpoints with expired token
	testCases := []struct {
		name       string
		prepareReq func() casesteps.HTTPRequest
	}{
		{
			name: "POST /quests - create quest",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.CreateQuestHTTPRequest(testdatagenerators.RandomCreateQuestRequest())
			},
		},
		{
			name: "POST /quests/{id}/assign - assign quest",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.AssignQuestHTTPRequest(createdQuest.ID())
			},
		},
		{
			name: "GET /quests/assigned - list assigned quests",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.ListAssignedQuestsHTTPRequest()
			},
		},
		{
			name: "GET /quests/{id} - get quest by id",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.GetQuestHTTPRequest(createdQuest.ID())
			},
		},
		{
			name: "GET /quests - list all quests",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.ListQuestsHTTPRequest("")
			},
		},
		{
			name: "PATCH /quests/{id}/status - change quest status",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.ChangeQuestStatusHTTPRequest(createdQuest.ID(), map[string]string{"status": "posted"})
			},
		},
		{
			name: "GET /quests/search-radius - search quests",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.SearchQuestsByRadiusHTTPRequest(55.7558, 37.6173, 10)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request with expired token
			req := tc.prepareReq()
			resp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithExpiredAuth, req)

			// Assert - should return 401 Unauthorized with token expired message
			s.Require().NoError(err, "HTTP request should not fail")
			s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode,
				"Endpoint should return 401 Unauthorized for expired token")
			s.Assert().Contains(resp.Body, "expired",
				"Error should mention token expiration")
		})
	}
}
