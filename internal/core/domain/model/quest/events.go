package quest

import (
	"quest-manager/internal/pkg/ddd"

	"github.com/google/uuid"
)

// QuestCreated represents quest creation event
type QuestCreated struct {
	ddd.BaseEvent
	Creator string `json:"creator"`
}

func NewQuestCreated(questID uuid.UUID, creator string) QuestCreated {
	return QuestCreated{
		BaseEvent: ddd.NewBaseEvent(questID, "quest.created"),
		Creator:   creator,
	}
}

// QuestAssigned represents quest assignment event
type QuestAssigned struct {
	ddd.BaseEvent
	UserID uuid.UUID `json:"user_id"`
}

func NewQuestAssigned(questID uuid.UUID, userID uuid.UUID) QuestAssigned {
	return QuestAssigned{
		BaseEvent: ddd.NewBaseEvent(questID, "quest.assigned"),
		UserID:    userID,
	}
}

// QuestStatusChanged represents quest status change event
type QuestStatusChanged struct {
	ddd.BaseEvent
	OldStatus Status `json:"old_status"`
	NewStatus Status `json:"new_status"`
}

func NewQuestStatusChanged(questID uuid.UUID, oldStatus, newStatus Status) QuestStatusChanged {
	return QuestStatusChanged{
		BaseEvent: ddd.NewBaseEvent(questID, "quest.status_changed"),
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
}
