package commands

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// ChangeQuestStatusCommandHandler defines the interface for handling quest status changes.
type ChangeQuestStatusCommandHandler interface {
	Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (quest.Quest, error)
}

type changeQuestStatusHandler struct {
	repo ports.QuestRepository
}

// NewChangeQuestStatusCommandHandler creates a new ChangeQuestStatusCommandHandler instance.
func NewChangeQuestStatusCommandHandler(repo ports.QuestRepository) ChangeQuestStatusCommandHandler {
	return &changeQuestStatusHandler{repo: repo}
}

// Handle updates the quest status with validation and domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (quest.Quest, error) {
	// Валидируем статус - это domain validation ошибка → 400
	if !quest.IsValidStatus(string(cmd.Status)) {
		return quest.Quest{}, errs.NewDomainValidationError("status", "must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'")
	}

	// Получаем квест - если не найден → 404
	q, err := h.repo.GetByID(ctx, cmd.QuestID)
	if err != nil {
		return quest.Quest{}, errs.NewNotFoundErrorWithCause("quest", cmd.QuestID.String(), err)
	}

	// Используем доменную логику для изменения статуса - domain validation ошибка → 400
	if err := q.ChangeStatus(cmd.Status); err != nil {
		return quest.Quest{}, errs.NewDomainValidationErrorWithCause("status", "invalid status transition", err)
	}

	// Сохраняем квест - infrastructure ошибка → 500
	if err := h.repo.Save(ctx, q); err != nil {
		return quest.Quest{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	return q, nil
}
