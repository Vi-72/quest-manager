package authentication

import (
	"context"
	"net/http"

	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

// AUTHENTICATION MIDDLEWARE TESTS - MISSING TOKEN
// Tests for requests without authentication token

func (s *Suite) TestCreateQuestWithoutToken() {
	ctx := context.Background()

	// Prepare valid quest data
	questData := testdatagenerators.SimpleQuestData(
		"Test Quest",
		"Quest without auth token",
		"easy",
		2,
		60,
		testdatagenerators.DefaultTestCoordinate(),
		testdatagenerators.DefaultTestCoordinate(),
	)

	// Act - send request WITHOUT Authorization header
	createReq := casesteps.HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests",
		Body:        questData,
		Headers:     nil, // Explicitly nil to avoid auto-adding auth
		ContentType: "application/json",
	}
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, createResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestAssignQuestWithoutToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - try to assign WITHOUT Authorization header
	assignReq := casesteps.HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + createdQuest.ID().String() + "/assign",
		Body:        nil,
		Headers:     nil, // Explicitly nil to avoid auto-adding auth
		ContentType: "application/json",
	}
	assignResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, assignReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, assignResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestListAssignedQuestsWithoutToken() {
	ctx := context.Background()

	// Act - try to list assigned quests WITHOUT Authorization header
	listReq := casesteps.HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests/assigned",
		Headers: nil, // Explicitly nil to avoid auto-adding auth
	}
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestGetQuestByIdWithoutToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - try to get quest WITHOUT Authorization header
	getReq := casesteps.HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests/" + createdQuest.ID().String(),
		Headers: nil, // Explicitly nil to avoid auto-adding auth
	}
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, getResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestListQuestsWithoutToken() {
	ctx := context.Background()

	// Act - try to list quests WITHOUT Authorization header
	listReq := casesteps.HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests",
		Headers: nil, // Explicitly nil to avoid auto-adding auth
	}
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestChangeQuestStatusWithoutToken() {
	ctx := context.Background()

	// Pre-condition - create quest using handler (with auth)
	createdQuest, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
	s.Require().NoError(err)

	// Act - try to change status WITHOUT Authorization header
	statusReq := casesteps.HTTPRequest{
		Method:      "PATCH",
		URL:         "/api/v1/quests/" + createdQuest.ID().String() + "/status",
		Body:        map[string]string{"status": "posted"},
		Headers:     nil, // Explicitly nil to avoid auto-adding auth
		ContentType: "application/json",
	}
	statusResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, statusReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, statusResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestSearchQuestsByRadiusWithoutToken() {
	ctx := context.Background()

	// Act - try to search quests WITHOUT Authorization header
	searchReq := casesteps.HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests/search-radius?lat=55.7558&lon=37.6173&radius_km=10",
		Headers: nil, // Explicitly nil to avoid auto-adding auth
	}
	searchResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, searchReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, searchResp.StatusCode, "Should return 401 Unauthorized without token")
}

func (s *Suite) TestWithEmptyAuthorizationHeader() {
	ctx := context.Background()

	// Act - send request with empty Authorization header value
	listReq := casesteps.HTTPRequest{
		Method: "GET",
		URL:    "/api/v1/quests",
		Headers: map[string]string{
			"Authorization": "", // Empty header value
		},
	}
	listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

	// Assert - should return 401 Unauthorized
	s.Require().NoError(err, "HTTP request should not fail")
	s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode, "Should return 401 Unauthorized with empty auth header")
}

func (s *Suite) TestWithInvalidAuthorizationHeaderFormat() {
	ctx := context.Background()

	// Test cases with invalid Authorization header formats
	testCases := []struct {
		name   string
		header string
	}{
		{
			name:   "no Bearer prefix",
			header: "just-a-token",
		},
		{
			name:   "wrong prefix",
			header: "Basic some-token",
		},
		{
			name:   "Bearer without token",
			header: "Bearer ",
		},
		{
			name:   "Bearer with only spaces",
			header: "Bearer    ",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request with invalid Authorization header format
			listReq := casesteps.HTTPRequest{
				Method: "GET",
				URL:    "/api/v1/quests",
				Headers: map[string]string{
					"Authorization": tc.header,
				},
			}
			listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)

			// Assert - should return 401 Unauthorized
			s.Require().NoError(err, "HTTP request should not fail")
			s.Assert().Equal(http.StatusUnauthorized, listResp.StatusCode,
				"Should return 401 Unauthorized with invalid auth header format: %s", tc.name)
		})
	}
}
