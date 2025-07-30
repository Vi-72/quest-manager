package http

import (
	"context"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/generated/servers"
)

// CreateQuest implements POST /api/v1/quests from OpenAPI.
func (a *ApiHandler) CreateQuest(ctx context.Context, request servers.CreateQuestRequestObject) (servers.CreateQuestResponseObject, error) {
	if request.Body == nil {
		return servers.CreateQuest400Response{}, nil
	}

	// Extract creator from context or set default (в реальном приложении это должно браться из токена аутентификации)
	creator := "system" // TODO: получать из токена пользователя

	targetLocation, err := APICoordinateToKernel(request.Body.TargetLocation)
	if err != nil {
		return servers.CreateQuest400Response{}, nil
	}

	executionLocation, err := APICoordinateToKernel(request.Body.ExecutionLocation)
	if err != nil {
		return servers.CreateQuest400Response{}, nil
	}

	equipment := []string{}
	if request.Body.Equipment != nil {
		equipment = *request.Body.Equipment
	}

	skills := []string{}
	if request.Body.Skills != nil {
		skills = *request.Body.Skills
	}

	reward := ""
	if request.Body.Reward != nil {
		reward = *request.Body.Reward
	}

	cmd := commands.CreateQuestCommand{
		Title:             request.Body.Title,
		Description:       request.Body.Description,
		Difficulty:        quest.Difficulty(request.Body.Difficulty),
		Reward:            reward,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Equipment:         equipment,
		Skills:            skills,
		Creator:           creator,
	}

	result, err := a.createQuestHandler.Handle(ctx, cmd)
	if err != nil {
		return servers.CreateQuest500Response{}, nil
	}

	// Get the created quest from repository to return full object
	createdQuest, err := a.getQuestByIDHandler.Handle(ctx, queries.GetQuestByIDQuery{ID: result.ID})
	if err != nil {
		return servers.CreateQuest500Response{}, nil
	}

	return servers.CreateQuest201JSONResponse(QuestToAPI(createdQuest.Quest)), nil
}
