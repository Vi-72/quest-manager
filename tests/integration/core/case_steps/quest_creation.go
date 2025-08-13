package casesteps

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

// CreateRandomQuestStep создает квест с рандомными данными
func CreateRandomQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
) (quest.Quest, error) {
	randomData := testdatagenerators.RandomQuestData()
	return CreateQuestStep(ctx, handler, randomData)
}

// CreateMultipleRandomQuests создает несколько квестов и проверяет ошибки
func CreateMultipleRandomQuests(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
	count int,
) ([]quest.Quest, error) {
	quests := make([]quest.Quest, 0, count)

	for i := 0; i < count; i++ {
		// Создаем квест
		createdQuest, err := CreateRandomQuestStep(ctx, handler)
		if err != nil {
			return nil, err
		}
		quests = append(quests, createdQuest)
	}

	return quests, nil
}

// CreateQuestStep creates quest using QuestTestData structure
func CreateQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
	questData testdatagenerators.QuestTestData,
) (quest.Quest, error) {
	// Convert to command using factory method
	cmd := questData.ToCreateCommand()

	// Generate test addresses for locations (handler layer needs addresses)
	targetAddress := "Target Address: Test Street 123, Test City"
	executionAddress := "Execution Address: Test Avenue 456, Test City"

	// Make sure execution location is slightly different from target to get different addresses
	if questData.TargetLocation.Equals(questData.ExecutionLocation) {
		cmd.ExecutionLocation = kernel.GeoCoordinate{
			Lat: questData.ExecutionLocation.Latitude() + 0.001,
			Lon: questData.ExecutionLocation.Longitude() + 0.001,
		}
		// Also make execution address different if coordinates are the same
		executionAddress = "Execution Address: Test Avenue 789, Different City"
	}

	// Set addresses (required by handler layer)
	cmd.TargetAddress = &targetAddress
	cmd.ExecutionAddress = &executionAddress

	return handler.Handle(ctx, cmd)
}
