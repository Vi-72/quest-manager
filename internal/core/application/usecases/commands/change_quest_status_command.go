package commands

import (
	"quest-manager/internal/core/domain/model/quest"

	"github.com/google/uuid"
)

type ChangeQuestStatusCommand struct {
	QuestID uuid.UUID
	Status  quest.Status
}

// ChangeQuestStatusResult represents the output after status change.
type ChangeQuestStatusResult struct {
	ID       uuid.UUID
	Assignee *string // can be nil if quest is not assigned
	Status   string
}
