package http

import (
	"context"
	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/generated/servers"
)

// GetQuestById implements GET /api/v1/quests/{quest_id} from OpenAPI.
func (a *ApiHandler) GetQuestById(ctx context.Context, request servers.GetQuestByIdRequestObject) (servers.GetQuestByIdResponseObject, error) {
	// Validate UUID
	questID, validationErr := validations.ValidateUUID(request.QuestId, "quest_id")
	if validationErr != nil {
		// Return validation error, middleware will automatically handle it and return 400 response
		return nil, validationErr
	}

	// Get quest directly
	quest, err := a.getQuestByIDHandler.Handle(ctx, questID)
	if err != nil {
		// Pass error to middleware for proper handling (404 for NotFoundError, 500 for others)
		return nil, err
	}

	apiQuest := QuestToAPI(quest)

	return servers.GetQuestById200JSONResponse(apiQuest), nil
}
