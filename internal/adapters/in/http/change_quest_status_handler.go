package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
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

	// Выполняем изменение статуса напрямую
	updatedQuest, err := a.changeQuestStatusHandler.Handle(ctx, validatedData.QuestID, quest.Status(validatedData.Status))
	if err != nil {
		// Передаем ошибку в middleware для правильной обработки (400, 404, 500)
		return nil, err
	}

	// Возвращаем обновленный квест
	apiQuest := QuestToAPI(updatedQuest)
	return servers.ChangeQuestStatus200JSONResponse(apiQuest), nil
}
