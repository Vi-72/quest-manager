package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/problems"
	httpValidations "quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// ChangeQuestStatus implements PATCH /api/v1/quests/{quest_id}/status from OpenAPI.
func (a *ApiHandler) ChangeQuestStatus(ctx context.Context, request servers.ChangeQuestStatusRequestObject) (servers.ChangeQuestStatusResponseObject, error) {
	// Валидация запроса с помощью централизованной функции
	validatedData, validationErr := httpValidations.ValidateChangeQuestStatusRequest(request.Body, request.QuestId)
	if validationErr != nil {
		// Возвращаем детальную ошибку через middleware обработчик
		return nil, validationErr
	}

	cmd := commands.ChangeQuestStatusCommand{
		ID:     validatedData.QuestID,
		Status: quest.Status(validatedData.Status),
	}

	_, err := a.changeQuestStatusHandler.Handle(ctx, cmd)
	if err != nil {
		return servers.ChangeQuestStatus500Response{}, nil
	}

	// Get updated quest from repository
	updatedQuest, err := a.getQuestByIDHandler.Handle(ctx, queries.GetQuestByIDQuery{ID: validatedData.QuestID})
	if err != nil {
		return nil, problems.NewNotFound("Quest not found")
	}

	apiQuest := QuestToAPI(updatedQuest.Quest)
	return servers.ChangeQuestStatus200JSONResponse(apiQuest), nil
}
