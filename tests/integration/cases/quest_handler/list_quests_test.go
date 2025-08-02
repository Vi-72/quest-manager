package quest_handler

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/tests/integration/core/assertions"
	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListQuests() {
	ctx := context.Background()

	testCases := []struct {
		desc    string
		prepare func() ([]quest.Quest, error)
		status  *quest.Status
		assert  func(created []quest.Quest, quests []quest.Quest, err error)
	}{
		{
			desc:   "list quests",
			status: nil,
			prepare: func() ([]quest.Quest, error) {
				expectedCount := 2
				return casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
			},
			assert: func(created []quest.Quest, quests []quest.Quest, err error) {
				s.Require().NoError(err)
				listAssertions := assertions.NewQuestListAssertions(s.Assert())
				listAssertions.QuestsHaveMinimumCount(quests, len(created))
				listAssertions.QuestsContainAllCreated(created, quests)
			},
		},
		{
			desc:   "list quests empty database",
			status: nil,
			prepare: func() ([]quest.Quest, error) {
				return []quest.Quest{}, nil
			},
			assert: func(created []quest.Quest, quests []quest.Quest, err error) {
				s.Require().NoError(err)
				s.Assert().Len(quests, 0)
			},
		},
		{
			desc:   "list quests with valid status",
			status: func() *quest.Status { st := quest.StatusPosted; return &st }(),
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
			assert: func(created []quest.Quest, quests []quest.Quest, err error) {
				s.Require().NoError(err)
				s.Assert().GreaterOrEqual(len(quests), 1, "Should have at least one quest with StatusPosted")
				listAssertions := assertions.NewQuestListAssertions(s.Assert())
				listAssertions.QuestsAllHaveStatus(quests, quest.StatusPosted)
				listAssertions.QuestWithIDExists(quests, created[0].ID().String())
			},
		},
		{
			desc:   "list quests with empty status",
			status: nil,
			prepare: func() ([]quest.Quest, error) {
				expectedCount := 2
				return casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
			},
			assert: func(created []quest.Quest, quests []quest.Quest, err error) {
				s.Require().NoError(err)
				listAssertions := assertions.NewQuestListAssertions(s.Assert())
				listAssertions.QuestsHaveMinimumCount(quests, len(created))
				listAssertions.QuestsContainAllCreated(created, quests)
			},
		},
		{
			desc:   "list quests with invalid status",
			status: func() *quest.Status { st := quest.Status("invalid_status_that_does_not_exist"); return &st }(),
			prepare: func() ([]quest.Quest, error) {
				_, err := casesteps.CreateRandomQuestStep(ctx, s.TestDIContainer.CreateQuestHandler)
				return nil, err
			},
			assert: func(created []quest.Quest, quests []quest.Quest, err error) {
				s.Require().Error(err, "Should return validation error for invalid status")
				s.Assert().Contains(err.Error(), "must be one of", "Error should mention valid status values")
				s.Assert().Nil(quests, "Quests should be nil when there's a validation error")
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.desc, func() {
			created, err := tc.prepare()
			s.Require().NoError(err)
			quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, tc.status)
			tc.assert(created, quests, err)
		})
	}
}
