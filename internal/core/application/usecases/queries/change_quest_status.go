package queries

import (
	"context"
	"fmt"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"

	"github.com/google/uuid"
)

// ChangeQuestStatusHandler defines the interface for handling quest status changes.
type ChangeQuestStatusHandler interface {
	Handle(ctx context.Context, questID uuid.UUID, status quest.Status) (quest.Quest, error)
}

type changeQuestStatusHandler struct {
	repo ports.QuestRepository
}

// NewChangeQuestStatusHandler creates a new ChangeQuestStatusHandler instance.
func NewChangeQuestStatusHandler(repo ports.QuestRepository) ChangeQuestStatusHandler {
	return &changeQuestStatusHandler{repo: repo}
}

// Handle updates the quest status with validation and domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, questID uuid.UUID, status quest.Status) (quest.Quest, error) {
	// Валидируем статус
	if !quest.IsValidStatus(string(status)) {
		return quest.Quest{}, fmt.Errorf("invalid quest status '%s': must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'", status)
	}

	// Получаем квест
	q, err := h.repo.GetByID(ctx, questID)
	if err != nil {
		return quest.Quest{}, fmt.Errorf("quest not found: %w", err)
	}

	// Используем доменную логику для изменения статуса
	if err := q.ChangeStatus(status); err != nil {
		return quest.Quest{}, fmt.Errorf("failed to change quest status: %w", err)
	}

	// Сохраняем квест
	if err := h.repo.Save(ctx, q); err != nil {
		return quest.Quest{}, fmt.Errorf("failed to save quest: %w", err)
	}

	return q, nil
}
