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

func (s *Suite) TestAllEndpointsRejectInvalidToken() {
	ctx := context.Background()

	// Pre-condition - create quest for tests that need it
	createdQuest, _ := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)

	// Setup - create router with invalid token auth client
	invalidTokenAuthClient := mock.NewInvalidTokenAuthClient()
	routerWithInvalidAuth := s.TestDIContainer.NewHTTPRouterWithAuthClient(invalidTokenAuthClient)

	// Test cases covering all endpoints with invalid token
	testCases := []struct {
		name        string
		prepareReq  func() casesteps.HTTPRequest
		checkError  bool
		errorString string
	}{
		{
			name: "POST /quests - create quest",
			prepareReq: func() casesteps.HTTPRequest {
				return casesteps.CreateQuestHTTPRequest(testdatagenerators.RandomCreateQuestRequest())
			},
			checkError:  true,
			errorString: "Authentication Failed",
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
			// Act - send request with invalid token
			req := tc.prepareReq()
			resp, err := casesteps.ExecuteHTTPRequest(ctx, routerWithInvalidAuth, req)

			// Assert - should return 401 Unauthorized
			s.Require().NoError(err, "HTTP request should not fail")
			s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode,
				"Endpoint should return 401 Unauthorized for invalid token")

			// Additional error message check if specified
			if tc.checkError && tc.errorString != "" {
				s.Assert().Contains(resp.Body, tc.errorString,
					"Error should mention: %s", tc.errorString)
			}
		})
	}
}
