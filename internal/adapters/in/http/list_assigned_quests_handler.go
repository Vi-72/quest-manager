package http

import (
	"context"

	"quest-manager/internal/adapters/in/http/validations"
	"quest-manager/internal/generated/servers"
)

// ListAssignedQuests implements GET /api/v1/quests/assigned from OpenAPI.
func (a *ApiHandler) ListAssignedQuests(ctx context.Context, request servers.ListAssignedQuestsRequestObject) (servers.ListAssignedQuestsResponseObject, error) {
	// Validate UUID for user_id
	_, validationErr := validations.ValidateUUID(request.Params.UserId, "user_id")
	if validationErr != nil {
		// Return validation error, middleware will automatically handle it and return 400 response
		return nil, validationErr
	}

	// Get quest list directly
	quests, err := a.listAssignedQuestsHandler.Handle(ctx, request.Params.UserId)
	if err != nil {
		// Pass error to middleware for proper handling
		return nil, err
	}

	apiQuests := make([]servers.Quest, 0)
	for _, q := range quests {
		apiQuests = append(apiQuests, QuestToAPI(q))
	}

	return servers.ListAssignedQuests200JSONResponse(apiQuests), nil
}
