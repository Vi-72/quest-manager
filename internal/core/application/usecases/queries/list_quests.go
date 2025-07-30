package queries

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// ListQuestsQuery represents the input parameters for listing quests.
// If Status is nil, all quests are returned.
type ListQuestsQuery struct {
	Status *quest.Status
}

// ListQuestsResult represents the result of listing quests.
type ListQuestsResult struct {
	Quests []quest.Quest
}

// ListQuestsQueryHandler defines the interface for handling ListQuestsQuery.
type ListQuestsQueryHandler interface {
	Handle(ctx context.Context, query ListQuestsQuery) (ListQuestsResult, error)
}

type listQuestsHandler struct {
	repo ports.QuestRepository
}

// NewListQuestsQueryHandler creates a new ListQuestsQueryHandler instance.
func NewListQuestsQueryHandler(repo ports.QuestRepository) ListQuestsQueryHandler {
	return &listQuestsHandler{repo: repo}
}

// Handle retrieves quests from the repository, optionally filtered by status.
func (h *listQuestsHandler) Handle(ctx context.Context, query ListQuestsQuery) (ListQuestsResult, error) {
	var quests []quest.Quest
	var err error

	quests, err = h.repo.FindAll(ctx)

	if err != nil {
		return ListQuestsResult{}, err
	}

	return ListQuestsResult{Quests: quests}, nil
}
