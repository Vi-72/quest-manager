package commands

import (
	"context"

	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
)

// PublishDomainEventsAsync publishes domain events asynchronously and clears them
func PublishDomainEventsAsync(ctx context.Context, eventPublisher ports.EventPublisher, aggregates ...ddd.AggregateRoot) {
	if eventPublisher == nil {
		return
	}

	var allEvents []ddd.DomainEvent

	// Collect events from all aggregates
	for _, aggregate := range aggregates {
		allEvents = append(allEvents, aggregate.GetDomainEvents()...)
	}

	// Publish events asynchronously
	if len(allEvents) > 0 {
		eventPublisher.PublishAsync(ctx, allEvents...)
	}

	// Clear events after queuing for publication
	for _, aggregate := range aggregates {
		aggregate.ClearDomainEvents()
	}
}
