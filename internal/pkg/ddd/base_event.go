package ddd

import (
	"time"

	"github.com/google/uuid"
)

// BaseEvent provides common structure for all domain events
type BaseEvent struct {
	ID          uuid.UUID `json:"id"`           // Event ID
	AggregateID uuid.UUID `json:"aggregate_id"` // Aggregate ID
	EventType   string    `json:"event_type"`   // Event type
	Timestamp   time.Time `json:"timestamp"`    // Event time
}

// GetID returns event ID
func (e BaseEvent) GetID() uuid.UUID {
	return e.ID
}

// GetName returns event type
func (e BaseEvent) GetName() string {
	return e.EventType
}

// GetAggregateID returns aggregate ID
func (e BaseEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// NewBaseEvent creates new base event
func NewBaseEvent(aggregateID uuid.UUID, eventType string) BaseEvent {
	return BaseEvent{
		ID:          uuid.New(),
		AggregateID: aggregateID,
		EventType:   eventType,
		Timestamp:   time.Now(),
	}
}
