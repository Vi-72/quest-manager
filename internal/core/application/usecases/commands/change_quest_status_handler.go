package commands

import (
	"context"

	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// ChangeQuestStatusCommandHandler defines the interface for handling quest status changes.
type ChangeQuestStatusCommandHandler interface {
	Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error)
}

type changeQuestStatusHandler struct {
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
}

// NewChangeQuestStatusCommandHandler creates a new ChangeQuestStatusCommandHandler instance.
func NewChangeQuestStatusCommandHandler(unitOfWork ports.UnitOfWork, eventPublisher ports.EventPublisher) ChangeQuestStatusCommandHandler {
	return &changeQuestStatusHandler{
		unitOfWork:     unitOfWork,
		eventPublisher: eventPublisher,
	}
}

// Handle updates the quest status with validation and domain business rules.
func (h *changeQuestStatusHandler) Handle(ctx context.Context, cmd ChangeQuestStatusCommand) (ChangeQuestStatusResult, error) {
	// Валидируем статус - это domain validation ошибка → 400
	if !quest.IsValidStatus(string(cmd.Status)) {
		return ChangeQuestStatusResult{}, errs.NewDomainValidationError("status", "must be one of 'created', 'posted', 'assigned', 'in_progress', 'declined', 'completed'")
	}

	// Начинаем транзакцию
	h.unitOfWork.Begin(ctx)

	// Получаем квест - если не найден → 404
	q, err := h.unitOfWork.QuestRepository().GetByID(ctx, cmd.QuestID)
	if err != nil {
		_ = h.unitOfWork.Rollback()
		return ChangeQuestStatusResult{}, errs.NewNotFoundErrorWithCause("quest", cmd.QuestID.String(), err)
	}

	// Используем доменную логику для изменения статуса - domain validation ошибка → 400
	if err := q.ChangeStatus(cmd.Status); err != nil {
		_ = h.unitOfWork.Rollback()
		return ChangeQuestStatusResult{}, errs.NewDomainValidationErrorWithCause("status", "invalid status transition", err)
	}

	// Сохраняем квест - infrastructure ошибка → 500
	if err := h.unitOfWork.QuestRepository().Save(ctx, q); err != nil {
		_ = h.unitOfWork.Rollback()
		return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Публикуем доменные события в рамках той же транзакции
	if h.eventPublisher != nil {
		if err := h.eventPublisher.Publish(ctx, q.GetDomainEvents()...); err != nil {
			_ = h.unitOfWork.Rollback()
			return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to publish events", err)
		}
	}

	// Коммитим транзакцию
	err = h.unitOfWork.Commit(ctx)
	if err != nil {
		return ChangeQuestStatusResult{}, errs.WrapInfrastructureError("failed to commit quest status change transaction", err)
	}

	// Очищаем события после успешного коммита
	q.ClearDomainEvents()

	// Формируем результат из обновленного квеста
	var assignee *string
	if q.Assignee != nil {
		assignee = q.Assignee
	}

	return ChangeQuestStatusResult{
		ID:       q.ID(),
		Assignee: assignee,
		Status:   string(q.Status),
	}, nil
}
