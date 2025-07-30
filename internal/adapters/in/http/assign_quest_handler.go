package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/generated/servers"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request servers.AssignQuestRequestObject) (servers.AssignQuestResponseObject, error) {
	// Валидация запроса
	validatedData, validationErr := validations.ValidateAssignQuestRequest(request.Body, request.QuestId)
	if validationErr != nil {
		// Возвращаем ошибку валидации, middleware автоматически обработает её и вернет 400 ответ
		return nil, validationErr
	}

	// Выполняем команду назначения
	cmd := commands.AssignQuestCommand{
		ID:     validatedData.QuestID,
		UserID: validatedData.UserID,
	}

	result, err := a.assignQuestHandler.Handle(ctx, cmd)
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки (400 для validation, 404 для not found, 500 для infrastructure)
		return nil, err
	}

	// Получаем обновленный квест для возврата
	quest, err := a.getQuestByIDHandler.Handle(ctx, result.ID)
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки
		return nil, err
	}

	apiQuest := QuestToAPI(quest)
	return servers.AssignQuest200JSONResponse(apiQuest), nil
}
