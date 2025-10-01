package quest_handler_tests

import (
	"context"
	"sync"

	"quest-manager/internal/core/domain/model/quest"
	casesteps "quest-manager/tests/integration/core/case_steps"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"

	"golang.org/x/sync/errgroup"
)

// TestCreateQuestHandlerParallelRequests ensures that each command execution
// uses an independent transactional scope so parallel requests do not
// interfere with each other.
func (s *Suite) TestCreateQuestHandlerParallelRequests() {
	ctx := context.Background()

	questDataA := testdatagenerators.NewQuest(testdatagenerators.WithTitle("Parallel Quest A"))
	questDataB := testdatagenerators.NewQuest(testdatagenerators.WithTitle("Parallel Quest B"))

	quests := make([]quest.Quest, 0, 2)
	var mu sync.Mutex

	g := errgroup.Group{}
	g.Go(func() error {
		created, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questDataA)
		if err != nil {
			return err
		}
		mu.Lock()
		quests = append(quests, created)
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		created, err := casesteps.CreateQuestStep(ctx, s.TestDIContainer.CreateQuestHandler, questDataB)
		if err != nil {
			return err
		}
		mu.Lock()
		quests = append(quests, created)
		mu.Unlock()
		return nil
	})

	s.Require().NoError(g.Wait())
	s.Require().Len(quests, 2)
	s.Assert().NotEqual(quests[0].ID(), quests[1].ID(), "quests created in parallel should have unique IDs")

	for _, createdQuest := range quests {
		savedQuest, err := s.TestDIContainer.QuestRepository.GetByID(ctx, createdQuest.ID())
		s.Require().NoError(err)
		s.Assert().Equal(createdQuest.ID(), savedQuest.ID())
	}
}
