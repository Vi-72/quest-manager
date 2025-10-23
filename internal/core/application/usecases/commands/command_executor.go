package commands

import (
	"context"

	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
	"quest-manager/internal/pkg/errs"
)

// CommandExecutor wraps command execution in UnitOfWork transaction with automatic rollback
type CommandExecutor struct {
	uowFactory     ports.UnitOfWorkFactory
	eventPublisher ports.EventPublisher
}

func NewCommandExecutor(uowFactory ports.UnitOfWorkFactory, eventPublisher ports.EventPublisher) *CommandExecutor {
	return &CommandExecutor{
		uowFactory:     uowFactory,
		eventPublisher: eventPublisher,
	}
}

// Execute runs business logic within a transaction, handling Begin/Commit/Rollback automatically
// Returns aggregates for event publishing after successful commit
func (e *CommandExecutor) Execute(ctx context.Context, fn func(ctx context.Context, uow ports.UnitOfWork) ([]ddd.AggregateRoot, error)) error {
	uow, err := e.uowFactory.CreateUnitOfWork()
	if err != nil {
		return errs.WrapInfrastructureError("failed to create unit of work", err)
	}

	if err := uow.Begin(ctx); err != nil {
		return errs.WrapInfrastructureError("failed to begin transaction", err)
	}

	// Propagate UoW in context so downstream (e.g., EventPublisher) can reuse current transaction
	ctxWithUow := ports.WithUnitOfWork(ctx, uow)

	aggregates, err := fn(ctxWithUow, uow)
	if err != nil {
		_ = uow.Rollback()
		return err
	}

	if err := uow.Commit(ctx); err != nil {
		return errs.WrapInfrastructureError("failed to commit transaction", err)
	}

	// Publish events after successful commit
	if e.eventPublisher != nil && len(aggregates) > 0 {
		PublishDomainEventsAsync(context.Background(), e.eventPublisher, aggregates...)
	}

	return nil
}
