package casesteps

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

// CreateQuestStep выполняет создание квеста
func CreateQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
	title, description, difficulty string,
	reward, durationMinutes int,
	creator string,
	targetLocation, executionLocation kernel.GeoCoordinate,
	equipment, skills []string,
) (quest.Quest, error) {
	cmd := commands.CreateQuestCommand{
		Title:             title,
		Description:       description,
		Difficulty:        difficulty,
		Reward:            reward,
		DurationMinutes:   durationMinutes,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Creator:           creator,
		Equipment:         equipment,
		Skills:            skills,
	}

	return handler.Handle(ctx, cmd)
}
