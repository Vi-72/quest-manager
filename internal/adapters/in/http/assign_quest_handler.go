package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request v1.AssignQuestRequestObject) (v1.AssignQuestResponseObject, error) {
	// Validate request
	validatedData, validationErr := validations.ValidateAssignQuestRequest(request.Body, request.QuestId.String())
	if validationErr != nil {
		// Return validation error, middleware will automatically handle it and return 400 response
		return nil, validationErr
	}

	// Execute assignment command
	cmd := commands.AssignQuestCommand{
		ID:     validatedData.QuestID,
		UserID: validatedData.UserID,
	}

	result, err := a.assignQuestHandler.Handle(ctx, cmd)
	if err != nil {
		// Pass error to middleware for proper handling (400 for validation, 404 for not found, 500 for infrastructure)
		return nil, err
	}

	// Form response from operation result
	apiResult := v1.AssignQuestResult{
		Id:       result.ID.String(),
		Assignee: result.Assignee,
		Status:   v1.QuestStatus(result.Status),
	}
	return v1.AssignQuest200JSONResponse(apiResult), nil
}
