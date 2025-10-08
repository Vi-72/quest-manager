package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/adapters/in/http/middleware"
	"quest-manager/internal/core/application/usecases/commands"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request v1.AssignQuestRequestObject) (v1.AssignQuestResponseObject, error) {
	// Get authenticated user ID from context (set by auth middleware)
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.NewBadRequest("user ID not found in context")
	}

	// Use user ID from JWT token instead of request body
	cmd := commands.AssignQuestCommand{
		ID:     request.QuestId,
		UserID: userID,
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
