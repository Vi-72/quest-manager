package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/problems"
	httpValidations "quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/generated/servers"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request servers.AssignQuestRequestObject) (servers.AssignQuestResponseObject, error) {
	// Валидация запроса с помощью централизованной функции
	validatedData, validationErr := httpValidations.ValidateAssignQuestRequest(request.Body, request.QuestId)
	if validationErr != nil {
		// Возвращаем детальную ошибку через middleware обработчик
		return nil, validationErr
	}

	cmd := commands.AssignQuestCommand{
		ID:     validatedData.QuestID,
		UserID: validatedData.UserID,
	}

	_, err := a.assignQuestHandler.Handle(ctx, cmd)
	if err != nil {
		return servers.AssignQuest500Response{}, nil
	}

	// Get updated quest from repository
	updatedQuest, err := a.getQuestByIDHandler.Handle(ctx, queries.GetQuestByIDQuery{ID: validatedData.QuestID})
	if err != nil {
		return nil, problems.NewNotFound("Quest not found")
	}

	apiQuest := QuestToAPI(updatedQuest.Quest)
	return servers.AssignQuest200JSONResponse(apiQuest), nil
}
