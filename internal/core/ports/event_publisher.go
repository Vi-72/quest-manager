package ports

import (
	"context"
	"quest-manager/internal/pkg/ddd"
)

// EventPublisher defines methods for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, events ...ddd.DomainEvent) error
	PublishAsync(ctx context.Context, events ...ddd.DomainEvent)
}

// NullEventPublisher is a no-op implementation for development
type NullEventPublisher struct{}

func (p *NullEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	// В production здесь была бы интеграция с message broker (RabbitMQ, Kafka и т.д.)
	// Пока что просто логируем события
	for _, event := range events {
		// TODO: добавить структурированное логирование
		_ = event
	}
	return nil
}

func (p *NullEventPublisher) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	// Асинхронная версия - просто вызываем синхронную
	_ = p.Publish(ctx, events...)
}
