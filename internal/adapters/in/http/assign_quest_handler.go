package http

import (
	"context"

	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/generated/servers"
)

// AssignQuest implements POST /api/v1/quests/{quest_id}/assign from OpenAPI.
func (a *ApiHandler) AssignQuest(ctx context.Context, request servers.AssignQuestRequestObject) (servers.AssignQuestResponseObject, error) {
	// Validate request
	validatedData, validationErr := validations.ValidateAssignQuestRequest(request.Body, request.QuestId)
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
	apiResult := servers.AssignQuestResult{
		Id:       result.ID.String(),
		Assignee: result.Assignee,
		Status:   servers.QuestStatus(result.Status),
	}
	return servers.AssignQuest200JSONResponse(apiResult), nil
}
