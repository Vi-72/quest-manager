package quest_handler

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
