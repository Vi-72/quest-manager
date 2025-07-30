package queries

import (
	"context"
	"fmt"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
)

// ListQuestsQueryHandler defines the interface for handling quest listing.
// If status is nil, all quests are returned. Otherwise, filters by status.
type ListQuestsQueryHandler interface {
	Handle(ctx context.Context, status *quest.Status) ([]quest.Quest, error)
}

type listQuestsHandler struct {
	repo ports.QuestRepository
}

// NewListQuestsQueryHandler creates a new ListQuestsQueryHandler instance.
func NewListQuestsQueryHandler(repo ports.QuestRepository) ListQuestsQueryHandler {
	return &listQuestsHandler{repo: repo}
}

// Handle retrieves quests from the repository, optionally filtered by status.
func (h *listQuestsHandler) Handle(ctx context.Context, status *quest.Status) ([]quest.Quest, error) {
	if status != nil {
		// Валидируем статус используя доменную логику
		if !quest.IsValidStatus(string(*status)) {
			return nil, fmt.Errorf("invalid quest status '%s': must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'", *status)
		}

		// Фильтруем по статусу
		return h.repo.FindByStatus(ctx, *status)
	}
	// Возвращаем все квесты
	return h.repo.FindAll(ctx)
}
