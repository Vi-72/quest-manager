package casesteps

import (
	"context"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
	testdatagenerators "quest-manager/tests/integration/core/test_data_generators"
)

// convertCreateQuestRequestToCommand from CreateQuestRequest to CreateQuestCommand
func convertCreateQuestRequestToCommand(req *servers.CreateQuestRequest) commands.CreateQuestCommand {
	return commands.CreateQuestCommand{
		Title:           req.Title,
		Description:     req.Description,
		Difficulty:      string(req.Difficulty),
		Reward:          req.Reward,
		DurationMinutes: req.DurationMinutes,
		TargetLocation: kernel.GeoCoordinate{
			Lat: float64(req.TargetLocation.Latitude),
			Lon: float64(req.TargetLocation.Longitude),
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat: float64(req.ExecutionLocation.Latitude),
			Lon: float64(req.ExecutionLocation.Longitude),
		},
		Creator:   "test-creator",
		Equipment: getStringSlice(req.Equipment),
		Skills:    getStringSlice(req.Skills),
	}
}

// getStringSlice безопасно извлекает slice из указателя
func getStringSlice(ptr *[]string) []string {
	if ptr == nil {
		return []string{}
	}
	return *ptr
}

// CreateRandomQuestStep создает квест с рандомными данными
func CreateRandomQuestStep(
	ctx context.Context,
	handler commands.CreateQuestCommandHandler,
) (quest.Quest, error) {
	requestData := testdatagenerators.RandomQuestData()
	questCmd := convertCreateQuestRequestToCommand(requestData)

	return handler.Handle(ctx, questCmd)
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
