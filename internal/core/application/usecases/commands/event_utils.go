package commands

import (
	"context"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
)

// PublishDomainEventsAsync публикует доменные события асинхронно и очищает их
func PublishDomainEventsAsync(ctx context.Context, eventPublisher ports.EventPublisher, aggregates ...ddd.AggregateRoot) {
	if eventPublisher == nil {
		return
	}

	var allEvents []ddd.DomainEvent

	// Собираем события от всех агрегатов
	for _, aggregate := range aggregates {
		allEvents = append(allEvents, aggregate.GetDomainEvents()...)
	}

	// Публикуем события асинхронно
	if len(allEvents) > 0 {
		eventPublisher.PublishAsync(ctx, allEvents...)
	}

	// Очищаем события после постановки в очередь на публикацию
	for _, aggregate := range aggregates {
		aggregate.ClearDomainEvents()
	}
}
