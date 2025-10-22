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

	aggregates, err := fn(ctx, uow)
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

// ExecuteWithResult runs business logic and returns a result along with aggregates
func (e *CommandExecutor) ExecuteWithResult(ctx context.Context, fn func(ctx context.Context, uow ports.UnitOfWork) (interface{}, []ddd.AggregateRoot, error)) (interface{}, error) {
	var result interface{}

	err := e.Execute(ctx, func(ctx context.Context, uow ports.UnitOfWork) ([]ddd.AggregateRoot, error) {
		res, aggs, err := fn(ctx, uow)
		result = res
		return aggs, err
	})

	return result, err
}
