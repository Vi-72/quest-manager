package eventrepo

import (
	"encoding/json"
	"time"
)

// EventDTO is the database model for all Domain Events.
type EventDTO struct {
	ID          string    `gorm:"primaryKey"`     // event_id
	EventType   string    `gorm:"index;not null"` // event_type: quest.created, location.created, etc.
	AggregateID string    `gorm:"index;not null"` // aggregate_id: Aggregate ID (quest, location, etc.)
	Data        string    `gorm:"type:jsonb"`     // data: JSON event data
	CreatedAt   time.Time `gorm:"index"`          // event creation date
}

func (EventDTO) TableName() string {
	return "events"
}

// MarshalEventData serializes event data to JSON
func MarshalEventData(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
