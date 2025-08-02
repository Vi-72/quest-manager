package casesteps

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

// CreateQuestStep creates quest with specific parameters
func CreateQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
	title, description, difficulty string,
	reward, durationMinutes int,
	creator string,
	targetLocation, executionLocation kernel.GeoCoordinate,
	equipment, skills []string,
) (quest.Quest, error) {
	// Generate test addresses for locations
	targetAddress := "Target Address: Test Street 123, Test City"
	executionAddress := "Execution Address: Test Avenue 456, Test City"

	// Make sure execution location is slightly different from target to get different addresses
	if targetLocation.Equals(executionLocation) {
		executionLocation = kernel.GeoCoordinate{
			Lat: executionLocation.Latitude() + 0.001,
			Lon: executionLocation.Longitude() + 0.001,
		}
		// Also make execution address different if coordinates are the same
		executionAddress = "Execution Address: Test Avenue 789, Different City"
	}

	cmd := commands.CreateQuestCommand{
		Title:             title,
		Description:       description,
		Difficulty:        difficulty,
		Reward:            reward,
		DurationMinutes:   durationMinutes,
		TargetLocation:    targetLocation,
		TargetAddress:     &targetAddress,
		ExecutionLocation: executionLocation,
		ExecutionAddress:  &executionAddress,
		Creator:           creator,
		Equipment:         equipment,
		Skills:            skills,
	}

	return handler.Handle(ctx, cmd)
}
