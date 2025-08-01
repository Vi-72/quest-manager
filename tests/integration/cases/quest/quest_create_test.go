package quest

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
)

func (s *Suite) TestCreateQuest() {
	ctx := context.Background()

	// Arrange
	cmd := commands.CreateQuestCommand{
		Title:           "Test Quest",
		Description:     "Test Description",
		Difficulty:      "easy",
		Reward:          3,
		DurationMinutes: 30,
		Creator:         "550e8400-e29b-41d4-a716-446655440000",
		TargetLocation: kernel.GeoCoordinate{
			Lat: 55.7558,
			Lon: 37.6176,
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat: 55.7558,
			Lon: 37.6176,
		},
	}

	// Act
	result, err := s.TestDIContainer.CreateQuestHandler.Handle(ctx, cmd)

	// Assert
	s.Require().NoError(err)
	s.Assert().NotEmpty(result.ID())
	s.Assert().Equal(cmd.Title, result.Title)
	s.Assert().Equal(cmd.Description, result.Description)
}

func (s *Suite) TestGetQuestByID() {
	ctx := context.Background()

	// Arrange - сначала создаем квест
	createCmd := commands.CreateQuestCommand{
		Title:           "Test Quest for Get",
		Description:     "Test Description for Get",
		Difficulty:      "medium",
		Reward:          4,
		DurationMinutes: 45,
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

	// Создаем квест
	createdQuest, err := s.TestDIContainer.CreateQuestHandler.Handle(ctx, createCmd)
	s.Require().NoError(err)
	s.Require().NotEmpty(createdQuest.ID())

	// Act - получаем квест по ID
	foundQuest, err := s.TestDIContainer.GetQuestByIDHandler.Handle(ctx, createdQuest.ID())

	// Assert
	s.Require().NoError(err)
	s.Assert().Equal(createdQuest.ID(), foundQuest.ID())
	s.Assert().Equal(createCmd.Title, foundQuest.Title)
	s.Assert().Equal(createCmd.Description, foundQuest.Description)
	s.Assert().Equal(createCmd.Difficulty, string(foundQuest.Difficulty))
	s.Assert().Equal(createCmd.Reward, foundQuest.Reward)
	s.Assert().Equal(createCmd.DurationMinutes, foundQuest.DurationMinutes)
}
