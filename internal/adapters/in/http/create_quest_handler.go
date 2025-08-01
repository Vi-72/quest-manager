package http

import (
	"context"
	httpValidations "quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/generated/servers"
)

// CreateQuest implements POST /api/v1/quests from OpenAPI.
func (a *ApiHandler) CreateQuest(ctx context.Context, request servers.CreateQuestRequestObject) (servers.CreateQuestResponseObject, error) {
	// Validate request and get processed data
	validatedData, validationErr := httpValidations.ValidateCreateQuestRequest(request.Body)
	if validationErr != nil {
		// Return detailed error through middleware handler
		return nil, validationErr
	}

	// Extract creator from context or set default (in real app this should be taken from auth token)
	creator := "system" // TODO: get from user token

	cmd := commands.CreateQuestCommand{
		Title:             validatedData.Title,
		Description:       validatedData.Description,
		Difficulty:        validatedData.Difficulty, // Pass string directly
		Reward:            validatedData.Reward,
		DurationMinutes:   validatedData.DurationMinutes,
		TargetLocation:    validatedData.TargetLocation,
		ExecutionLocation: validatedData.ExecutionLocation,
		Equipment:         validatedData.Equipment,
		Skills:            validatedData.Skills,
		Creator:           creator,
	}

	result, err := a.createQuestHandler.Handle(ctx, cmd)
	if err != nil {
		// Pass error to middleware for proper handling (400 for validation, 500 for infrastructure)
		return nil, err
	}

	// Return full response using mapper
	response := QuestToAPI(result)

	return servers.CreateQuest201JSONResponse(response), nil
}
