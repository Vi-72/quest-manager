package assertions

import (
	"encoding/json"
	"net/http"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/generated/servers"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

type QuestHTTPAssertions struct {
	assert *assert.Assertions
}

func NewQuestHTTPAssertions(a *assert.Assertions) *QuestHTTPAssertions {
	return &QuestHTTPAssertions{assert: a}
}

// QuestHTTPCreatedSuccessfully verifies HTTP quest creation success and parses response
// This is the most useful helper as it eliminates HTTP boilerplate code
func (a *QuestHTTPAssertions) QuestHTTPCreatedSuccessfully(createResp *casesteps.HTTPResponse, err error) servers.Quest {
	a.assert.NoError(err, "HTTP request should not fail")
	a.assert.Equal(http.StatusCreated, createResp.StatusCode, "Should return 201 Created")

	var createdQuest servers.Quest
	parseErr := json.Unmarshal([]byte(createResp.Body), &createdQuest)
	a.assert.NoError(parseErr, "Response should be valid JSON")
	a.assert.NotEmpty(createdQuest.Id, "Quest should have ID")
	a.assert.False(createdQuest.CreatedAt.IsZero(), "CreatedAt should be set")
	a.assert.False(createdQuest.UpdatedAt.IsZero(), "UpdatedAt should be set")

	return createdQuest
}

// QuestHTTPValidationError verifies HTTP validation error response
// Useful for API validation tests
func (a *QuestHTTPAssertions) QuestHTTPValidationError(createResp *casesteps.HTTPResponse, err error, expectedField string) {
	a.assert.NoError(err, "HTTP request should not fail")
	a.assert.Equal(http.StatusBadRequest, createResp.StatusCode, "Should return 400 for validation error")
	a.assert.Contains(createResp.Body, "validation failed", "Error should mention validation failure")
	if expectedField != "" {
		a.assert.Contains(createResp.Body, expectedField, "Error should mention the specific field")
	}
}

// QuestArraysNotNull verifies equipment and skills arrays are not null (HTTP level)
// Ensures [] instead of null serialization
func (a *QuestHTTPAssertions) QuestArraysNotNull(httpQuest servers.Quest) {
	a.assert.NotNil(httpQuest.Equipment, "Equipment should not be null")
	a.assert.NotNil(httpQuest.Skills, "Skills should not be null")
}

// QuestHTTPGetSuccessfully verifies HTTP GET quest success and parses response
// Eliminates boilerplate for single quest retrieval
func (a *QuestHTTPAssertions) QuestHTTPGetSuccessfully(getResp *casesteps.HTTPResponse, err error) servers.Quest {
	a.assert.NoError(err, "HTTP request should not fail")
	a.assert.Equal(http.StatusOK, getResp.StatusCode, "Should return 200 OK")

	var foundQuest servers.Quest
	parseErr := json.Unmarshal([]byte(getResp.Body), &foundQuest)
	a.assert.NoError(parseErr, "Response should be valid JSON")

	return foundQuest
}

// QuestHTTPListSuccessfully verifies HTTP LIST quests success and parses response
// Eliminates boilerplate for quest list retrieval
func (a *QuestHTTPAssertions) QuestHTTPListSuccessfully(listResp *casesteps.HTTPResponse, err error) []servers.Quest {
	a.assert.NoError(err, "HTTP request should not fail")
	a.assert.Equal(http.StatusOK, listResp.StatusCode, "Should return 200 OK")

	var quests []servers.Quest
	parseErr := json.Unmarshal([]byte(listResp.Body), &quests)
	a.assert.NoError(parseErr, "Response should be valid JSON")

	return quests
}

// QuestHTTPErrorResponse verifies HTTP error response
// Eliminates boilerplate for error checking
func (a *QuestHTTPAssertions) QuestHTTPErrorResponse(resp *casesteps.HTTPResponse, err error, expectedStatus int, expectedMessage string) {
	a.assert.NoError(err, "HTTP request should not fail")
	a.assert.Equal(expectedStatus, resp.StatusCode, "Should return expected status code")
	if expectedMessage != "" {
		a.assert.Contains(resp.Body, expectedMessage, "Error should contain expected message")
	}
}
