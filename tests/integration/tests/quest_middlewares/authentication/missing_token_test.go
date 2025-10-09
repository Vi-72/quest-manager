package authentication

import (
	"context"
	"net/http"

	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

// AUTHENTICATION MIDDLEWARE TESTS - MISSING TOKEN
// Tests for requests without authentication token

func (s *Suite) TestAllEndpointsRequireAuthentication() {
	ctx := context.Background()

	// Pre-condition - create quest for tests that need it
	createdQuest, _ := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)

	// Test cases covering all endpoints without authentication
	testCases := []struct {
		name        string
		method      string
		url         string
		body        interface{}
		contentType string
	}{
		{
			name:        "POST /quests - create quest",
			method:      "POST",
			url:         "/api/v1/quests",
			body:        testdatagenerators.RandomCreateQuestRequest(),
			contentType: "application/json",
		},
		{
			name:        "POST /quests/{id}/assign - assign quest",
			method:      "POST",
			url:         "/api/v1/quests/" + createdQuest.ID().String() + "/assign",
			contentType: "application/json",
		},
		{
			name:   "GET /quests/assigned - list assigned quests",
			method: "GET",
			url:    "/api/v1/quests/assigned",
		},
		{
			name:   "GET /quests/{id} - get quest by id",
			method: "GET",
			url:    "/api/v1/quests/" + createdQuest.ID().String(),
		},
		{
			name:   "GET /quests - list all quests",
			method: "GET",
			url:    "/api/v1/quests",
		},
		{
			name:        "PATCH /quests/{id}/status - change quest status",
			method:      "PATCH",
			url:         "/api/v1/quests/" + createdQuest.ID().String() + "/status",
			body:        map[string]string{"status": "posted"},
			contentType: "application/json",
		},
		{
			name:   "GET /quests/search-radius - search quests by radius",
			method: "GET",
			url:    "/api/v1/quests/search-radius?lat=55.7558&lon=37.6173&radius_km=10",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request WITHOUT Authorization header
			req := casesteps.HTTPRequest{
				Method:      tc.method,
				URL:         tc.url,
				Body:        tc.body,
				ContentType: tc.contentType,
				SkipAuth:    true, // Don't add default Bearer token
			}
			resp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, req)

			// Assert - should return 401 Unauthorized
			s.Require().NoError(err, "HTTP request should not fail")
			s.Assert().Equal(http.StatusUnauthorized, resp.StatusCode,
				"Endpoint %s %s should return 401 Unauthorized without token", tc.method, tc.url)
		})
	}
}
