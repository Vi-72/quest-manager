package http

import (
	"context"

	v1 "quest-manager/api/http/quests/v1"
)

// GetQuestById implements GET /api/v1/quests/{quest_id} from OpenAPI.
func (a *ApiHandler) GetQuestById(ctx context.Context, request v1.GetQuestByIdRequestObject) (v1.GetQuestByIdResponseObject, error) {
	// QuestId is already UUID type from OpenAPI, just convert to uuid.UUID
	questID := request.QuestId

	// Get quest directly
	quest, err := a.getQuestByIDHandler.Handle(ctx, questID)
	if err != nil {
		// Pass error to middleware for proper handling (404 for NotFoundError, 500 for others)
		return nil, err
	}

	apiQuest := QuestToAPI(quest)

	return v1.GetQuestById200JSONResponse(apiQuest), nil
}
