package assertions

import (
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"

	"github.com/stretchr/testify/assert"
)

// QuestFieldAssertions provides reusable assertions for quest field comparisons
type QuestFieldAssertions struct {
	assert *assert.Assertions
}

// NewQuestFieldAssertions creates a new instance of quest field assertions
func NewQuestFieldAssertions(assert *assert.Assertions) *QuestFieldAssertions {
	return &QuestFieldAssertions{
		assert: assert,
	}
}

// VerifyQuestMatchesHTTPRequest verifies that a quest domain object matches an HTTP request
func (a *QuestFieldAssertions) VerifyQuestMatchesHTTPRequest(quest quest.Quest, request *servers.CreateQuestRequest) {
	a.assert.Equal(request.Title, quest.Title, "Quest title should match HTTP request")
	a.assert.Equal(request.Description, quest.Description, "Quest description should match HTTP request")
	a.assert.Equal(string(request.Difficulty), string(quest.Difficulty), "Quest difficulty should match HTTP request")
	a.assert.Equal(request.Reward, quest.Reward, "Quest reward should match HTTP request")
	a.assert.Equal(request.DurationMinutes, quest.DurationMinutes, "Quest duration should match HTTP request")

	// Handle optional arrays
	if request.Equipment != nil {
		a.assert.Equal(*request.Equipment, quest.Equipment, "Quest equipment should match HTTP request")
	}
	if request.Skills != nil {
		a.assert.Equal(*request.Skills, quest.Skills, "Quest skills should match HTTP request")
	}
}

// VerifyHTTPResponseMatchesRequest verifies that an HTTP response matches the original request
func (a *QuestFieldAssertions) VerifyHTTPResponseMatchesRequest(response *servers.Quest, request *servers.CreateQuestRequest) {
	a.assert.Equal(request.Title, response.Title, "HTTP response title should match request")
	a.assert.Equal(request.Description, response.Description, "HTTP response description should match request")
	a.assert.Equal(string(request.Difficulty), string(response.Difficulty), "HTTP response difficulty should match request")
	a.assert.Equal(request.Reward, response.Reward, "HTTP response reward should match request")
	a.assert.Equal(request.DurationMinutes, response.DurationMinutes, "HTTP response duration should match request")
}
