package quest_operations

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
)

func (s *Suite) TestListQuests() {
	ctx := context.Background()

	// Arrange - создаем несколько квестов
	quest1Cmd := commands.CreateQuestCommand{
		Title:           "First Quest",
		Description:     "First Description",
		Difficulty:      "easy",
		Reward:          2,
		DurationMinutes: 20,
		Creator:         "550e8400-e29b-41d4-a716-446655440001",
		TargetLocation: kernel.GeoCoordinate{
			Lat: 55.7558,
			Lon: 37.6176,
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat: 55.7558,
			Lon: 37.6176,
		},
	}

	quest2Cmd := commands.CreateQuestCommand{
		Title:           "Second Quest",
		Description:     "Second Description",
		Difficulty:      "hard",
		Reward:          5,
		DurationMinutes: 60,
		Creator:         "550e8400-e29b-41d4-a716-446655440002",
		TargetLocation: kernel.GeoCoordinate{
			Lat: 56.8431,
			Lon: 60.6454,
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat: 56.8431,
			Lon: 60.6454,
		},
	}

	// Создаем квесты
	_, err := s.TestDIContainer.CreateQuestHandler.Handle(ctx, quest1Cmd)
	s.Require().NoError(err)

	_, err = s.TestDIContainer.CreateQuestHandler.Handle(ctx, quest2Cmd)
	s.Require().NoError(err)

	// Act - получаем список квестов (все квесты, без фильтра по статусу)
	quests, err := s.TestDIContainer.ListQuestsHandler.Handle(ctx, nil)

	// Assert
	s.Require().NoError(err)
	s.Assert().Len(quests, 2)

	// Проверяем, что квесты содержат ожидаемые данные
	questTitles := make([]string, len(quests))
	for i, quest := range quests {
		questTitles[i] = quest.Title
	}

	s.Assert().Contains(questTitles, "First Quest")
	s.Assert().Contains(questTitles, "Second Quest")
}

func (s *Suite) TestListQuestsEmpty() {
	ctx := context.Background()

	// Act - получаем список квестов из пустой базы
	quests, err := s.TestDIContainer.ListQuestsHandler.Handle(ctx, nil)

	// Assert
	s.Require().NoError(err)
	s.Assert().Len(quests, 0)
}
