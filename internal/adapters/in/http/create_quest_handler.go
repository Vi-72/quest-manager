package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
	httpValidations "quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/core/application/usecases/commands"
)

// CreateQuest implements POST /api/v1/quests from OpenAPI.
func (a *ApiHandler) CreateQuest(ctx context.Context, request v1.CreateQuestRequestObject) (v1.CreateQuestResponseObject, error) {
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
		TargetAddress:     validatedData.TargetAddress,
		ExecutionLocation: validatedData.ExecutionLocation,
		ExecutionAddress:  validatedData.ExecutionAddress,
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

	return v1.CreateQuest201JSONResponse(response), nil
}
