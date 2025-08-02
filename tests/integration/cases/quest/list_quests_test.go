package quest

import (
	"context"

	casesteps "quest-manager/tests/integration/core/case_steps"
)

func (s *Suite) TestListQuests() {
	ctx := context.Background()

	// Pre-condition - создаем несколько квестов
	expectedCount := 2
	createdQuests, err := casesteps.CreateMultipleRandomQuests(ctx, s.TestDIContainer.CreateQuestHandler, expectedCount)
	s.Require().NoError(err)
	createdQuestIDs := make(map[string]bool)
	for _, quest := range createdQuests {
		createdQuestIDs[quest.ID().String()] = true
	}

	// Act - получаем список квестов (все квесты, без фильтра по статусу)
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, nil)

	// Assert
	s.Require().NoError(err)
	s.Assert().GreaterOrEqual(len(quests), expectedCount, "Should have at least the created quests")

	// Проверяем, что все созданные квесты присутствуют в полученном списке
	retrievedQuestIDs := make(map[string]bool)
	for _, quest := range quests {
		retrievedQuestIDs[quest.ID().String()] = true
	}

	for _, createdQuest := range createdQuests {
		s.Assert().Contains(retrievedQuestIDs, createdQuest.ID().String(), "Created quest should be in the retrieved list")
	}
}

func (s *Suite) TestListQuestsEmpty() {
	ctx := context.Background()

	// Act - получаем список квестов из пустой базы
	quests, err := casesteps.ListQuestsStep(ctx, s.TestDIContainer.ListQuestsHandler, nil)

	// Assert
	s.Require().NoError(err)
	s.Assert().Len(quests, 0)
}
