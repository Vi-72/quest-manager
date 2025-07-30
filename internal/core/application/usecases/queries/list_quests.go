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

// ListQuestsQueryHandler defines the interface for handling ListQuestsQuery.
type ListQuestsQueryHandler interface {
	Handle(ctx context.Context, query ListQuestsQuery) ([]quest.Quest, error)
}

type listQuestsHandler struct {
	repo ports.QuestRepository
}

// NewListQuestsQueryHandler creates a new ListQuestsQueryHandler instance.
func NewListQuestsQueryHandler(repo ports.QuestRepository) ListQuestsQueryHandler {
	return &listQuestsHandler{repo: repo}
}

// Handle retrieves quests from the repository, optionally filtered by status.
func (h *listQuestsHandler) Handle(ctx context.Context, query ListQuestsQuery) ([]quest.Quest, error) {
	return h.repo.FindAll(ctx)
}
