package eventrepo

import (
	"encoding/json"
	"time"
)

// EventDTO is the database model for all Domain Events.
type EventDTO struct {
	ID          string    `gorm:"primaryKey"`     // event_id
	EventType   string    `gorm:"index;not null"` // event_type: quest.created, location.created, etc.
	AggregateID string    `gorm:"index;not null"` // aggregate_id: ID агрегата (квест, локация, etc.)
	Data        string    `gorm:"type:jsonb"`     // data: JSON данные события
	CreatedAt   time.Time `gorm:"index"`          // дата создания события
}

func (EventDTO) TableName() string {
	return "events"
}

// MarshalEventData сериализует данные события в JSON
func MarshalEventData(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// UnmarshalEventData десериализует JSON данные события
func UnmarshalEventData(data string, target interface{}) error {
	return json.Unmarshal([]byte(data), target)
}
