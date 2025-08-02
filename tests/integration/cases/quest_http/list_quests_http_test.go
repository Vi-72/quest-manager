package quest_http

import (
	"context"
	"encoding/json"
	"net/http"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListQuestsHTTP() {
	ctx := context.Background()

	testCases := []struct {
		desc         string
		prepare      func() ([]quest.Quest, error)
		statusQuery  string
		expectedCode int
		assert       func(created []quest.Quest, resp *casesteps.HTTPResponse)
	}{
		{
			desc:         "list quests HTTP",
			statusQuery:  "",
			expectedCode: http.StatusOK,
			prepare: func() ([]quest.Quest, error) {
				expectedCount := 2
				return casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
			},
			assert: func(created []quest.Quest, resp *casesteps.HTTPResponse) {
				var quests []servers.Quest
				err := json.Unmarshal([]byte(resp.Body), &quests)
				s.Require().NoError(err)
				s.Assert().GreaterOrEqual(len(quests), len(created), "Should return at least %d quests", len(created))
				returnedQuestIDs := make(map[string]bool)
				for _, q := range quests {
					returnedQuestIDs[q.Id] = true
				}
				for _, createdQuest := range created {
					questID := createdQuest.ID().String()
					s.Assert().True(returnedQuestIDs[questID], "Created quest %s should be in quests list", questID)
				}
			},
		},
		{
			desc:         "list quests HTTP empty",
			statusQuery:  "",
			expectedCode: http.StatusOK,
			prepare: func() ([]quest.Quest, error) {
				return []quest.Quest{}, nil
			},
			assert: func(created []quest.Quest, resp *casesteps.HTTPResponse) {
				var quests []servers.Quest
				err := json.Unmarshal([]byte(resp.Body), &quests)
				s.Require().NoError(err)
				s.Assert().Len(quests, 0, "Should return empty list when no quests exist")
			},
		},
		{
			desc:         "list quests HTTP with valid status",
			statusQuery:  string(quest.StatusPosted),
			expectedCode: http.StatusOK,
			prepare: func() ([]quest.Quest, error) {
				expectedCount := 3
				created, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
				if err != nil {
					return nil, err
				}
				targetStatus := quest.StatusPosted
				_, err = casesteps.ChangeQuestStatusStep(ctx, s.TestDIContainer.ChangeQuestStatusHandler,
					s.TestDIContainer.QuestRepository, created[0].ID(), targetStatus)
				return created, err
			},
			assert: func(created []quest.Quest, resp *casesteps.HTTPResponse) {
				var quests []servers.Quest
				err := json.Unmarshal([]byte(resp.Body), &quests)
				s.Require().NoError(err)
				s.Assert().GreaterOrEqual(len(quests), 1, "Should have at least one quest with StatusPosted")
				for _, q := range quests {
					s.Assert().Equal(string(quest.StatusPosted), string(q.Status), "All quests should have StatusPosted")
				}
				foundTargetQuest := false
				targetQuestID := created[0].ID().String()
				for _, q := range quests {
					if q.Id == targetQuestID {
						foundTargetQuest = true
						break
					}
				}
				s.Assert().True(foundTargetQuest, "Quest with StatusPosted should be in filtered list")
			},
		},
		{
			desc:         "list quests HTTP with empty status",
			statusQuery:  "",
			expectedCode: http.StatusOK,
			prepare: func() ([]quest.Quest, error) {
				expectedCount := 2
				return casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
			},
			assert: func(created []quest.Quest, resp *casesteps.HTTPResponse) {
				var quests []servers.Quest
				err := json.Unmarshal([]byte(resp.Body), &quests)
				s.Require().NoError(err)
				s.Assert().GreaterOrEqual(len(quests), len(created), "Should return at least %d quests", len(created))
				returnedQuestIDs := make(map[string]bool)
				for _, q := range quests {
					returnedQuestIDs[q.Id] = true
				}
				for _, createdQuest := range created {
					questID := createdQuest.ID().String()
					s.Assert().True(returnedQuestIDs[questID], "Created quest %s should be in quests list", questID)
				}
			},
		},
		{
			desc:         "list quests HTTP with invalid status",
			statusQuery:  "invalid_status_that_does_not_exist",
			expectedCode: http.StatusBadRequest,
			prepare: func() ([]quest.Quest, error) {
				_, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
				return nil, err
			},
			assert: func(created []quest.Quest, resp *casesteps.HTTPResponse) {
				s.Assert().Contains(resp.Body, "validation failed", "Error message should contain validation failure details")
				s.Assert().Contains(resp.Body, "must be one of", "Error message should mention valid status values")
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.desc, func() {
			created, err := tc.prepare()
			s.Require().NoError(err)

			listReq := casesteps.ListQuestsHTTPRequest(tc.statusQuery)
			listResp, err := casesteps.ExecuteHTTPRequest(ctx, s.TestDIContainer.HTTPRouter, listReq)
			s.Require().NoError(err)
			s.Require().Equal(tc.expectedCode, listResp.StatusCode)
			tc.assert(created, listResp)
		})
	}
}
