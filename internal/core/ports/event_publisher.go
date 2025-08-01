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
	// In production there would be integration with message broker (RabbitMQ, Kafka, etc.)
	// For now just log events
	for _, event := range events {
		// TODO: add structured logging
		_ = event
	}
	return nil
}

func (p *NullEventPublisher) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	// Async version - just call synchronous
	_ = p.Publish(ctx, events...)
}
