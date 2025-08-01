package ddd

import (
	"time"

	"github.com/google/uuid"
)

// BaseEvent предоставляет общую структуру для всех доменных событий
type BaseEvent struct {
	ID          uuid.UUID `json:"id"`           // ID события
	AggregateID uuid.UUID `json:"aggregate_id"` // ID агрегата
	EventType   string    `json:"event_type"`   // тип события
	Timestamp   time.Time `json:"timestamp"`    // время события
}

// GetID возвращает ID события
func (e BaseEvent) GetID() uuid.UUID {
	return e.ID
}

// GetName возвращает тип события
func (e BaseEvent) GetName() string {
	return e.EventType
}

// NewBaseEvent создает новое базовое событие
func NewBaseEvent(aggregateID uuid.UUID, eventType string) BaseEvent {
	return BaseEvent{
		ID:          uuid.New(),
		AggregateID: aggregateID,
		EventType:   eventType,
		Timestamp:   time.Now(),
	}
}
