package http

import (
	"context"
	"net/http"

	v1 "quest-manager/api/http/quests/v1"
	"quest-manager/internal/adapters/in/http/errors"
	"quest-manager/internal/adapters/in/http/middleware"
	"quest-manager/internal/core/application/usecases/commands"
)

// CreateQuest implements POST /api/v1/quests from OpenAPI.
func (a *ApiHandler) CreateQuest(ctx context.Context, request v1.CreateQuestRequestObject) (v1.CreateQuestResponseObject, error) {
	if request.Body == nil {
		return nil, errors.NewBadRequest("request body is required")
	}

	targetLocation, err := convertAPICoordinateToKernel(request.Body.TargetLocation)
	if err != nil {
		return nil, errors.NewBadRequest("Request validation failed: target_location invalid coordinate values (" + err.Error() + ")")
	}

	executionLocation, err := convertAPICoordinateToKernel(request.Body.ExecutionLocation)
	if err != nil {
		return nil, errors.NewBadRequest("Request validation failed: execution_location invalid coordinate values (" + err.Error() + ")")
	}

	equipment := []string{}
	if request.Body.Equipment != nil {
		equipment = *request.Body.Equipment
	}

	skills := []string{}
	if request.Body.Skills != nil {
		skills = *request.Body.Skills
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, errors.NewProblem(http.StatusUnauthorized, "Unauthorized", "authentication context is missing user information")
	}
	creator := userID.String()

	cmd := commands.CreateQuestCommand{
		Title:             request.Body.Title,
		Description:       request.Body.Description,
		Difficulty:        string(request.Body.Difficulty),
		Reward:            request.Body.Reward,
		DurationMinutes:   request.Body.DurationMinutes,
		TargetLocation:    targetLocation,
		TargetAddress:     request.Body.TargetLocation.Address,
		ExecutionLocation: executionLocation,
		ExecutionAddress:  request.Body.ExecutionLocation.Address,
		Equipment:         equipment,
		Skills:            skills,
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
