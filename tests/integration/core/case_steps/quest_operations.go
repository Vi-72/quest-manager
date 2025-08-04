package casesteps

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

// CreateQuestStep creates quest using QuestTestData structure
func CreateQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
	questData testdatagenerators.QuestTestData,
) (quest.Quest, error) {
	// Generate test addresses for locations
	targetAddress := "Target Address: Test Street 123, Test City"
	executionAddress := "Execution Address: Test Avenue 456, Test City"

	// Make sure execution location is slightly different from target to get different addresses
	if questData.TargetLocation.Equals(questData.ExecutionLocation) {
		questData.ExecutionLocation = kernel.GeoCoordinate{
			Lat: questData.ExecutionLocation.Latitude() + 0.001,
			Lon: questData.ExecutionLocation.Longitude() + 0.001,
		}
		// Also make execution address different if coordinates are the same
		executionAddress = "Execution Address: Test Avenue 789, Different City"
	}

	cmd := commands.CreateQuestCommand{
		Title:             questData.Title,
		Description:       questData.Description,
		Difficulty:        questData.Difficulty,
		Reward:            questData.Reward,
		DurationMinutes:   questData.DurationMinutes,
		TargetLocation:    questData.TargetLocation,
		TargetAddress:     &targetAddress,
		ExecutionLocation: questData.ExecutionLocation,
		ExecutionAddress:  &executionAddress,
		Creator:           questData.Creator,
		Equipment:         questData.Equipment,
		Skills:            questData.Skills,
	}

	return handler.Handle(ctx, cmd)
}
