package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/quest"
)

// ChangeQuestStatus implements PATCH /api/v1/quests/{quest_id}/status from OpenAPI.
func (a *ApiHandler) ChangeQuestStatus(ctx context.Context, request v1.ChangeQuestStatusRequestObject) (v1.ChangeQuestStatusResponseObject, error) {
	// Validate request
	validatedData, validationErr := validations.ValidateChangeQuestStatusRequest(request.Body, request.QuestId.String())
	if validationErr != nil {
		// Return validation error, middleware will automatically handle it and return 400 response
		return nil, validationErr
	}

	// Execute status change through command
	cmd := commands.ChangeQuestStatusCommand{
		QuestID: validatedData.QuestID,
		Status:  quest.Status(validatedData.Status),
	}
	result, err := a.changeQuestStatusHandler.Handle(ctx, cmd)
	if err != nil {
		// Pass error to middleware for proper handling (400, 404, 500)
		return nil, err
	}

	// Form response from operation result
	apiResult := v1.ChangeQuestStatusResult{
		Id:       result.ID.String(),
		Assignee: result.Assignee,
		Status:   v1.QuestStatus(result.Status),
	}
	return v1.ChangeQuestStatus200JSONResponse(apiResult), nil
}
