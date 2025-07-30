package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// ChangeQuestStatus implements PATCH /api/v1/quests/{quest_id}/status from OpenAPI.
func (a *ApiHandler) ChangeQuestStatus(ctx context.Context, request servers.ChangeQuestStatusRequestObject) (servers.ChangeQuestStatusResponseObject, error) {
	// Валидация запроса
	validatedData, validationErr := validations.ValidateChangeQuestStatusRequest(request.Body, request.QuestId)
	if validationErr != nil {
		// Возвращаем ошибку валидации, middleware автоматически обработает её и вернет 400 ответ
		return nil, validationErr
	}

	// Выполняем команду изменения статуса
	cmd := commands.ChangeQuestStatusCommand{
		ID:     validatedData.QuestID,
		Status: quest.Status(validatedData.Status),
	}

	result, err := a.changeQuestStatusHandler.Handle(ctx, cmd)
	if err != nil {
		return servers.ChangeQuestStatus500Response{}, nil
	}

	// Получаем обновленный квест для возврата
	quest, err := a.getQuestByIDHandler.Handle(ctx, result.ID)
	if err != nil {
		return servers.ChangeQuestStatus500Response{}, nil
	}

	apiQuest := QuestToAPI(quest)
	return servers.ChangeQuestStatus200JSONResponse(apiQuest), nil
}
