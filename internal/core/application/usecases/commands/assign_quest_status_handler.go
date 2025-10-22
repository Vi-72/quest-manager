package commands

import (
	"context"

	_ "quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

// AssignQuestCommandHandler defines the interface for handling AssignQuestCommand.
type AssignQuestCommandHandler interface {
	Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error)
}

var _ AssignQuestCommandHandler = &assignQuestHandler{}

// assignQuestHandler implements AssignQuestCommandHandler.
type assignQuestHandler struct {
	executor *CommandExecutor
}

// NewAssignQuestCommandHandler creates a new instance of AssignQuestCommandHandler.
func NewAssignQuestCommandHandler(uowFactory ports.UnitOfWorkFactory, eventPublisher ports.EventPublisher) AssignQuestCommandHandler {
	return &assignQuestHandler{
		executor: NewCommandExecutor(uowFactory, eventPublisher),
	}
}

// Handle assigns a quest to a user using domain business rules.
func (h *assignQuestHandler) Handle(ctx context.Context, cmd AssignQuestCommand) (AssignQuestResult, error) {
	var result AssignQuestResult

	err := h.executor.Execute(ctx, func(ctx context.Context, uow ports.UnitOfWork) ([]ddd.AggregateRoot, error) {
		// Get quest - if not found → 404
		q, err := uow.QuestRepository().GetByID(ctx, cmd.ID)
		if err != nil {
			return nil, errs.NewNotFoundErrorWithCause("quest", cmd.ID.String(), err)
		}

		// Use domain logic - business rules errors → 400
		if err := q.AssignTo(cmd.UserID); err != nil {
			return nil, errs.NewDomainValidationErrorWithCause("assignment", "failed to assign quest", err)
		}

		// Save quest - infrastructure error → 500
		if err := uow.QuestRepository().Save(ctx, q); err != nil {
			return nil, errs.WrapInfrastructureError("failed to save quest", err)
		}

		result = AssignQuestResult{
			ID:       q.ID(),
			Assignee: cmd.UserID,
			Status:   string(q.Status),
		}

		// Return aggregate for event publishing
		return []ddd.AggregateRoot{&q}, nil
	})

	return result, err
}
