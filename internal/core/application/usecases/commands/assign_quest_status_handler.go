package commands

import (
	"context"

	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
)

// AssignQuestCommandHandler defines the interface for handling AssignQuestCommand.
type AssignQuestCommandHandler interface {
	Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error)
}

var _ AssignQuestCommandHandler = &assignQuestHandler{}

// assignQuestHandler implements AssignQuestCommandHandler.
type assignQuestHandler struct {
	repo           ports.QuestRepository
	eventPublisher ports.EventPublisher
}

// NewAssignQuestCommandHandler creates a new instance of AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(repo ports.QuestRepository, eventPublisher ports.EventPublisher) AssignQuestCommandHandler {
	return &assignQuestHandler{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

// Handle assigns a quest to a user using domain business rules.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	// Получаем квест - если не найден → 404
	q, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return AssignQuestResult{}, errs.NewNotFoundErrorWithCause("quest", cmd.ID.String(), err)
	}

	// Используем доменную логику - ошибки бизнес-правил → 400
	if err := q.AssignTo(cmd.UserID); err != nil {
		return AssignQuestResult{}, errs.NewDomainValidationErrorWithCause("assignment", "failed to assign quest", err)
	}

	// Сохраняем квест - infrastructure ошибка → 500
	if err := h.repo.Save(ctx, q); err != nil {
		return AssignQuestResult{}, errs.WrapInfrastructureError("failed to save quest", err)
	}

	// Публикуем доменные события (включая QuestAssigned и QuestStatusChanged)
	if h.eventPublisher != nil {
		err = h.eventPublisher.Publish(ctx, q.GetDomainEvents()...)
		if err != nil {
			return AssignQuestResult{}, errs.WrapInfrastructureError("failed to publish quest events", err)
		}
	}

	// Очищаем события после успешной публикации
	q.ClearDomainEvents()

	return AssignQuestResult{
		ID:       q.ID(),
		Assignee: cmd.UserID,
		Status:   string(q.Status),
	}, nil
}
