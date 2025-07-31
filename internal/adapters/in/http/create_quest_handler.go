package http

import (
	"context"
	httpValidations "quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/generated/servers"
)

// CreateQuest implements POST /api/v1/quests from OpenAPI.
func (a *ApiHandler) CreateQuest(ctx context.Context, request servers.CreateQuestRequestObject) (servers.CreateQuestResponseObject, error) {
	// Валидация запроса и получение обработанных данных
	validatedData, validationErr := httpValidations.ValidateCreateQuestRequest(request.Body)
	if validationErr != nil {
		// Возвращаем детальную ошибку через middleware обработчик
		return nil, validationErr
	}

	// Extract creator from context or set default (в реальном приложении это должно браться из токена аутентификации)
	creator := "system" // TODO: получать из токена пользователя

	cmd := commands.CreateQuestCommand{
		Title:             validatedData.Title,
		Description:       validatedData.Description,
		Difficulty:        validatedData.Difficulty, // Передаем string напрямую
		Reward:            validatedData.Reward,
		DurationMinutes:   validatedData.DurationMinutes,
		TargetLocation:    validatedData.TargetLocation,
		ExecutionLocation: validatedData.ExecutionLocation,
		Equipment:         validatedData.Equipment,
		Skills:            validatedData.Skills,
		Creator:           creator,
	}

	result, err := a.createQuestHandler.Handle(ctx, cmd)
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки (400 для validation, 500 для infrastructure)
		return nil, err
	}

	// Get the created quest from repository to return full object
	createdQuest, err := a.getQuestByIDHandler.Handle(ctx, result.ID)
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки
		return nil, err
	}

	return servers.CreateQuest201JSONResponse(QuestToAPI(createdQuest)), nil
}
