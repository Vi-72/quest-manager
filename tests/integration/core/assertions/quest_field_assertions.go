package assertions

import (
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

// VerifyHTTPResponseMatchesRequest verifies that an HTTP response matches the original request
func (a *QuestFieldAssertions) VerifyHTTPResponseMatchesRequest(response *servers.Quest, request *servers.CreateQuestRequest) {
	a.assert.Equal(request.Title, response.Title, "HTTP response title should match request")
	a.assert.Equal(request.Description, response.Description, "HTTP response description should match request")
	a.assert.Equal(string(request.Difficulty), string(response.Difficulty), "HTTP response difficulty should match request")
	a.assert.Equal(request.Reward, response.Reward, "HTTP response reward should match request")
	a.assert.Equal(request.DurationMinutes, response.DurationMinutes, "HTTP response duration should match request")
}
