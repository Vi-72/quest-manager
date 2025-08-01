package quest

import (
	"time"

	"github.com/google/uuid"
)

// QuestCreated представляет событие создания квеста
type QuestCreated struct {
	ID        uuid.UUID `json:"id"`       // ID события
	QuestID   uuid.UUID `json:"quest_id"` // ID квеста
	Creator   string    `json:"creator"`
	Timestamp time.Time `json:"timestamp"`
}

func (e QuestCreated) GetID() uuid.UUID {
	return e.ID
}

func (e QuestCreated) GetName() string {
	return "quest.created"
}

// QuestAssigned представляет событие назначения квеста
type QuestAssigned struct {
	ID        uuid.UUID `json:"id"`       // ID события
	QuestID   uuid.UUID `json:"quest_id"` // ID квеста
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (e QuestAssigned) GetID() uuid.UUID {
	return e.ID
}

func (e QuestAssigned) GetName() string {
	return "quest.assigned"
}

// QuestStatusChanged представляет событие изменения статуса квеста
type QuestStatusChanged struct {
	ID        uuid.UUID `json:"id"`       // ID события
	QuestID   uuid.UUID `json:"quest_id"` // ID квеста
	OldStatus Status    `json:"old_status"`
	NewStatus Status    `json:"new_status"`
	Timestamp time.Time `json:"timestamp"`
}

func (e QuestStatusChanged) GetID() uuid.UUID {
	return e.ID
}

func (e QuestStatusChanged) GetName() string {
	return "quest.status_changed"
}
