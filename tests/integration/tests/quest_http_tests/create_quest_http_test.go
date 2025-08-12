package quest_http_tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"quest-manager/internal/generated/servers"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateCreateQuestRequest function

// HTTPRequest represents HTTP request for testing (duplicated here for convenience)
type HTTPRequest struct {
	Method      string
	URL         string
	Body        interface{}
	Headers     map[string]string
	ContentType string
}

func (s *Suite) TestCreateQuestHTTP() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	fieldAssertions := assertions.NewQuestFieldAssertions(s.Assert())

	// Pre-condition - prepare quest data
	questRequest := testdatagenerators.RandomCreateQuestRequest()

	// Act - create quest via HTTP API
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)
	httpAssertions.QuestArraysNotNull(createdQuest)

	// Verify quest data matches request using field assertions pattern
	fieldAssertions.VerifyHTTPResponseMatchesRequest(&createdQuest, questRequest)
}

func (s *Suite) TestCreateQuestHTTPWithEmptyArrays() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - prepare quest data with empty arrays
	questRequest := &servers.CreateQuestRequest{
		Title:           "Empty Arrays Quest",
		Description:     "Quest with empty equipment and skills",
		Difficulty:      servers.CreateQuestRequestDifficultyEasy,
		Reward:          2,
		DurationMinutes: 30,
		TargetLocation: servers.Coordinate{
			Latitude:  55.7558,
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Latitude:  55.7560,
			Longitude: 37.6178,
		},
		Equipment: &[]string{}, // Empty array
		Skills:    &[]string{}, // Empty array
	}

	// Act - create quest via HTTP API
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)
	httpAssertions.QuestArraysNotNull(createdQuest)

	// Specifically verify empty arrays are returned as [] not null
	s.Assert().Len(*createdQuest.Equipment, 0, "Equipment should be empty array")
	s.Assert().Len(*createdQuest.Skills, 0, "Skills should be empty array")

	// Verify other data
	s.Assert().Equal(questRequest.Title, createdQuest.Title)
	s.Assert().Equal(questRequest.Description, createdQuest.Description)
	s.Assert().Equal(string(questRequest.Difficulty), string(createdQuest.Difficulty))
}

// API LAYER VALIDATION TESTS
// Only tests that correspond to ValidateCreateQuestRequest function

func (s *Suite) TestCreateQuestHTTPMissingRequiredFields() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Act - send request with empty JSON body to test ValidateBody function
	emptyBodyRequest := map[string]interface{}{} // Empty object

	createReq := casesteps.CreateQuestHTTPRequest(emptyBodyRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert - API layer should reject incomplete body
	httpAssertions.QuestHTTPValidationError(createResp, err, "")
}

func (s *Suite) TestCreateQuestHTTPEmptyStringFields() {
	ctx := context.Background()

	// Test cases with empty string fields (TrimAndValidateString function tests)
	testCases := []struct {
		name    string
		request map[string]interface{}
		field   string
	}{
		{
			name:    "empty title - TrimAndValidateString",
			request: testdatagenerators.HTTPQuestDataWithField("title", ""),
			field:   "title",
		},
		{
			name:    "empty description - TrimAndValidateString",
			request: testdatagenerators.HTTPQuestDataWithField("description", ""),
			field:   "description",
		},
		{
			name:    "empty difficulty - ValidateNotEmpty",
			request: testdatagenerators.HTTPQuestDataWithField("difficulty", ""),
			field:   "difficulty",
		},
	}

	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request with empty field
			createReq := casesteps.CreateQuestHTTPRequest(tc.request)
			createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

			// Assert - API layer should reject empty required fields
			httpAssertions.QuestHTTPValidationError(createResp, err, tc.field)
		})
	}
}

func (s *Suite) TestCreateQuestHTTPNegativeNumbers() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Test cases with negative numbers (API layer technical validation)
	testCases := []struct {
		name    string
		request map[string]interface{}
		field   string
	}{
		{
			name:    "negative reward",
			request: testdatagenerators.HTTPQuestDataWithField("reward", -1),
			field:   "reward",
		},
		{
			name:    "negative duration",
			request: testdatagenerators.HTTPQuestDataWithField("duration_minutes", -30),
			field:   "duration_minutes",
		},
		{
			name:    "zero duration",
			request: testdatagenerators.HTTPQuestDataWithField("duration_minutes", 0),
			field:   "duration_minutes",
		},
		{
			name:    "zero reward",
			request: testdatagenerators.HTTPQuestDataWithField("reward", 0),
			field:   "reward",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act - send request with negative numbers
			createReq := casesteps.CreateQuestHTTPRequest(tc.request)
			createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

			// Assert - API layer should reject negative numbers
			httpAssertions.QuestHTTPValidationError(createResp, err, tc.field)
		})
	}
}

// Note: Coordinate validation tests removed from HTTP layer
// Domain coordinate validation is tested in tests/domain/kernel_coordinates_test.go
// HTTP layer focuses on API-specific validations (format, required fields, etc.)

func (s *Suite) TestCreateQuestHTTPMalformedJSON() {
	ctx := context.Background()

	// Create malformed JSON request
	malformedRequest := HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests",
		Body:        `{"title": "Invalid JSON", "description": }`, // Malformed JSON
		ContentType: "application/json",
	}

	// Act - send malformed JSON
	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, malformedRequest.Method, malformedRequest.URL, strings.NewReader(malformedRequest.Body.(string)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", malformedRequest.ContentType)

	s.TestDIContainer.HTTPRouter.ServeHTTP(recorder, req)

	// Assert - API layer should reject malformed JSON
	s.Require().Equal(http.StatusBadRequest, recorder.Code, "Should return 400 for malformed JSON")
}

// Note: Content-Type validation is handled at framework level, not application level
// The OpenAPI generated server accepts JSON regardless of Content-Type header

func (s *Suite) TestCreateQuestHTTPPersistence() {
	ctx := context.Background()
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())

	// Pre-condition - prepare quest data
	questRequest := testdatagenerators.RandomCreateQuestRequest()

	// Act - create quest via HTTP API
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert creation
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Verify quest is persisted by retrieving it via HTTP API
	getReq := casesteps.GetQuestHTTPRequest(createdQuest.Id)
	getResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, getReq)

	// Assert retrieval
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, getResp.StatusCode)

	// Parse retrieved quest
	var retrievedQuest servers.Quest
	err = json.Unmarshal([]byte(getResp.Body), &retrievedQuest)
	s.Require().NoError(err)

	// Verify quests match
	s.Assert().Equal(createdQuest.Id, retrievedQuest.Id)
	s.Assert().Equal(createdQuest.Title, retrievedQuest.Title)
	s.Assert().Equal(createdQuest.Description, retrievedQuest.Description)
	s.Assert().Equal(createdQuest.Equipment, retrievedQuest.Equipment)
	s.Assert().Equal(createdQuest.Skills, retrievedQuest.Skills)
}

func (s *Suite) TestCreateQuestHTTPWithLocationAddresses() {
	ctx := context.Background()

	// Pre-condition - prepare quest data with addresses
	targetAddress := "Moscow, Red Square"
	executionAddress := "Moscow, Gorky Park"

	questRequest := &servers.CreateQuestRequest{
		Title:           "Quest with Addresses",
		Description:     "Quest that has location addresses",
		Difficulty:      servers.CreateQuestRequestDifficultyMedium,
		Reward:          3,
		DurationMinutes: 60,
		TargetLocation: servers.Coordinate{
			Address:   &targetAddress,
			Latitude:  55.7558,
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &executionAddress,
			Latitude:  55.7560,
			Longitude: 37.6178,
		},
		Equipment: &[]string{"map", "camera"},
		Skills:    &[]string{"navigation", "photography"},
	}

	// Act - create quest via HTTP API
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)
	httpAssertions.QuestArraysNotNull(createdQuest)

	// Verify HTTP response structure matches domain expectations
	singleAssertions := assertions.NewQuestSingleAssertions(s.Assert())
	singleAssertions.QuestHTTPHasValidLocationData(createdQuest)
	s.Assert().Len(*createdQuest.Equipment, 2, "Equipment should have 2 items")
	s.Assert().Len(*createdQuest.Skills, 2, "Skills should have 2 items")
	s.Assert().Contains(*createdQuest.Equipment, "map")
	s.Assert().Contains(*createdQuest.Skills, "navigation")
}

// LOCATION-SPECIFIC HTTP TESTS

func (s *Suite) TestCreateQuestHTTPWithoutAddresses() {
	ctx := context.Background()

	// Pre-condition - prepare quest data without addresses (only coordinates)
	questRequest := &servers.CreateQuestRequest{
		Title:           "Quest Without Address",
		Description:     "Quest with coordinates but no addresses",
		Difficulty:      servers.CreateQuestRequestDifficultyEasy,
		Reward:          2,
		DurationMinutes: 45,
		TargetLocation: servers.Coordinate{
			Latitude:  55.7558,
			Longitude: 37.6176,
			// No Address field
		},
		ExecutionLocation: servers.Coordinate{
			Latitude:  55.7560,
			Longitude: 37.6178,
			// No Address field
		},
		Equipment: &[]string{"notebook"},
		Skills:    &[]string{"observation"},
	}

	// Act - create quest via HTTP API
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Verify quest was created successfully without addresses
	s.Assert().Equal(questRequest.Title, createdQuest.Title)
	s.Assert().Equal(questRequest.TargetLocation.Latitude, createdQuest.TargetLocation.Latitude)
	s.Assert().Equal(questRequest.TargetLocation.Longitude, createdQuest.TargetLocation.Longitude)

	// Verify no addresses are set (should be nil or empty)
	if createdQuest.TargetLocation.Address != nil {
		s.Assert().Empty(*createdQuest.TargetLocation.Address, "Target address should be empty when not provided")
	}
	if createdQuest.ExecutionLocation.Address != nil {
		s.Assert().Empty(*createdQuest.ExecutionLocation.Address, "Execution address should be empty when not provided")
	}
}

func (s *Suite) TestCreateQuestHTTPWithSameLocations() {
	ctx := context.Background()

	// Pre-condition - prepare quest data with identical target and execution locations
	sameAddress := "Moscow, Red Square"
	questRequest := &servers.CreateQuestRequest{
		Title:           "Quest with Same Locations",
		Description:     "Quest where target and execution are the same place",
		Difficulty:      servers.CreateQuestRequestDifficultyMedium,
		Reward:          3,
		DurationMinutes: 60,
		TargetLocation: servers.Coordinate{
			Address:   &sameAddress,
			Latitude:  55.7558,
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &sameAddress, // Same address
			Latitude:  55.7558,      // Same coordinates
			Longitude: 37.6176,
		},
		Equipment: &[]string{"camera"},
		Skills:    &[]string{"photography"},
	}

	// Act - create quest via HTTP API
	createReq := casesteps.CreateQuestHTTPRequest(questRequest)
	createResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq)

	// Assert
	httpAssertions := assertions.NewQuestHTTPAssertions(s.Assert())
	createdQuest := httpAssertions.QuestHTTPCreatedSuccessfully(createResp, err)

	// Verify both locations are identical
	s.Assert().Equal(createdQuest.TargetLocation.Latitude, createdQuest.ExecutionLocation.Latitude)
	s.Assert().Equal(createdQuest.TargetLocation.Longitude, createdQuest.ExecutionLocation.Longitude)

	if createdQuest.TargetLocation.Address != nil && createdQuest.ExecutionLocation.Address != nil {
		s.Assert().Equal(*createdQuest.TargetLocation.Address, *createdQuest.ExecutionLocation.Address)
	}
}

func (s *Suite) TestCreateQuestHTTPWithExistingLocationSameAddress() {
	ctx := context.Background()

	// Pre-condition - create first quest to establish location in database
	sharedAddress := "Moscow, Kremlin"
	firstQuestRequest := &servers.CreateQuestRequest{
		Title:           "First Quest",
		Description:     "First quest at this location",
		Difficulty:      servers.CreateQuestRequestDifficultyEasy,
		Reward:          2,
		DurationMinutes: 30,
		TargetLocation: servers.Coordinate{
			Address:   &sharedAddress,
			Latitude:  55.7520,
			Longitude: 37.6175,
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &sharedAddress,
			Latitude:  55.7520,
			Longitude: 37.6175,
		},
		Equipment: &[]string{},
		Skills:    &[]string{},
	}

	// Act - create first quest
	createReq1 := casesteps.CreateQuestHTTPRequest(firstQuestRequest)
	createResp1, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq1)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, createResp1.StatusCode)

	// Parse first quest
	var firstQuest servers.Quest
	err = json.Unmarshal([]byte(createResp1.Body), &firstQuest)
	s.Require().NoError(err)

	// Act - create second quest with exact same location and address (should reuse location)
	secondQuestRequest := &servers.CreateQuestRequest{
		Title:           "Second Quest",
		Description:     "Second quest at the same location",
		Difficulty:      servers.CreateQuestRequestDifficultyMedium,
		Reward:          3,
		DurationMinutes: 45,
		TargetLocation: servers.Coordinate{
			Address:   &sharedAddress, // Same address
			Latitude:  55.7520,        // Same coordinates
			Longitude: 37.6175,
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &sharedAddress, // Same address
			Latitude:  55.7520,        // Same coordinates
			Longitude: 37.6175,
		},
		Equipment: &[]string{"map"},
		Skills:    &[]string{"navigation"},
	}

	createReq2 := casesteps.CreateQuestHTTPRequest(secondQuestRequest)
	createResp2, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq2)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, createResp2.StatusCode)

	// Parse second quest
	var secondQuest servers.Quest
	err = json.Unmarshal([]byte(createResp2.Body), &secondQuest)
	s.Require().NoError(err)

	// Verify both quests exist but are different
	s.Assert().NotEqual(firstQuest.Id, secondQuest.Id, "Quests should have different IDs")
	s.Assert().NotEqual(firstQuest.Title, secondQuest.Title, "Quests should have different titles")

	// Verify locations are identical (coordinates and addresses)
	s.Assert().Equal(firstQuest.TargetLocation.Latitude, secondQuest.TargetLocation.Latitude)
	s.Assert().Equal(firstQuest.TargetLocation.Longitude, secondQuest.TargetLocation.Longitude)

	if firstQuest.TargetLocation.Address != nil && secondQuest.TargetLocation.Address != nil {
		s.Assert().Equal(*firstQuest.TargetLocation.Address, *secondQuest.TargetLocation.Address)
	}
}

func (s *Suite) TestCreateQuestHTTPWithExistingLocationDifferentAddress() {
	ctx := context.Background()

	// Pre-condition - create first quest with one address
	firstAddress := "Moscow, Red Square"
	firstQuestRequest := &servers.CreateQuestRequest{
		Title:           "Red Square Quest",
		Description:     "Quest at Red Square",
		Difficulty:      servers.CreateQuestRequestDifficultyEasy,
		Reward:          2,
		DurationMinutes: 30,
		TargetLocation: servers.Coordinate{
			Address:   &firstAddress,
			Latitude:  55.7558, // Same coordinates
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &firstAddress,
			Latitude:  55.7558,
			Longitude: 37.6176,
		},
		Equipment: &[]string{},
		Skills:    &[]string{},
	}

	// Act - create first quest
	createReq1 := casesteps.CreateQuestHTTPRequest(firstQuestRequest)
	createResp1, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq1)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, createResp1.StatusCode)

	// Parse first quest
	var firstQuest servers.Quest
	err = json.Unmarshal([]byte(createResp1.Body), &firstQuest)
	s.Require().NoError(err)

	// Act - create second quest with same coordinates but different address
	secondAddress := "Moscow, Historical Museum"
	secondQuestRequest := &servers.CreateQuestRequest{
		Title:           "Historical Museum Quest",
		Description:     "Quest at Historical Museum (same coordinates, different address)",
		Difficulty:      servers.CreateQuestRequestDifficultyMedium,
		Reward:          3,
		DurationMinutes: 45,
		TargetLocation: servers.Coordinate{
			Address:   &secondAddress, // Different address
			Latitude:  55.7558,        // Same coordinates
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &secondAddress, // Different address
			Latitude:  55.7558,        // Same coordinates
			Longitude: 37.6176,
		},
		Equipment: &[]string{"guidebook"},
		Skills:    &[]string{"history"},
	}

	createReq2 := casesteps.CreateQuestHTTPRequest(secondQuestRequest)
	createResp2, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, createReq2)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, createResp2.StatusCode)

	// Parse second quest
	var secondQuest servers.Quest
	err = json.Unmarshal([]byte(createResp2.Body), &secondQuest)
	s.Require().NoError(err)

	// Verify both quests exist but are different
	s.Assert().NotEqual(firstQuest.Id, secondQuest.Id, "Quests should have different IDs")
	s.Assert().NotEqual(firstQuest.Title, secondQuest.Title, "Quests should have different titles")

	// Verify coordinates are the same but addresses are different
	s.Assert().Equal(firstQuest.TargetLocation.Latitude, secondQuest.TargetLocation.Latitude)
	s.Assert().Equal(firstQuest.TargetLocation.Longitude, secondQuest.TargetLocation.Longitude)

	if firstQuest.TargetLocation.Address != nil && secondQuest.TargetLocation.Address != nil {
		s.Assert().NotEqual(*firstQuest.TargetLocation.Address, *secondQuest.TargetLocation.Address, "Addresses should be different")
		s.Assert().Equal(firstAddress, *firstQuest.TargetLocation.Address)
		s.Assert().Equal(secondAddress, *secondQuest.TargetLocation.Address)
	}
}
