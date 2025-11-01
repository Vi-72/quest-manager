package ports

import (
	"context"
	"log/slog"

	"quest-manager/internal/pkg/ddd"
)

// EventPublisher defines methods for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, events ...ddd.DomainEvent) error
}

// NullEventPublisher is a no-op implementation for development
type NullEventPublisher struct{}

func (p *NullEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	// In production, integrate with a message broker (RabbitMQ, Kafka, etc.)
	// For now, just log the events
	for _, event := range events {
		slog.InfoContext(ctx, "publishing domain event",
			slog.String("event_name", event.GetName()),
			slog.String("event_id", event.GetID().String()),
		)
	}
	return nil
}
