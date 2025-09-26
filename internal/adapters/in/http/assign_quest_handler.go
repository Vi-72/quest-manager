package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/core/application/usecases/commands"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request v1.AssignQuestRequestObject) (v1.AssignQuestResponseObject, error) {
	if request.Body == nil {
		return nil, errors.NewBadRequest("request body is required")
	}

	cmd := commands.AssignQuestCommand{
		ID:     request.QuestId,
		UserID: request.Body.UserId,
	}

	result, err := a.assignQuestHandler.Handle(ctx, cmd)
	if err != nil {
		// Pass error to middleware for proper handling (400 for validation, 404 for not found, 500 for infrastructure)
		return nil, err
	}

	// Form response from operation result
	apiResult := v1.AssignQuestResult{
		Id:       result.ID,
		Assignee: result.Assignee,
		Status:   v1.QuestStatus(result.Status),
	}
	return v1.AssignQuest200JSONResponse(apiResult), nil
}
