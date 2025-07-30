package commands

import (
	"github.com/google/uuid"
	"quest-manager/internal/core/domain/model/quest"
)

// ChangeQuestStatusCommand represents the input data required to change the quest status.
type ChangeQuestStatusCommand struct {
	ID       uuid.UUID
	Status   quest.Status
	Assignee *string // ID of the user to assign (required for "assigned")
}

// ChangeQuestStatusResult represents the result after the quest status has been updated.
type ChangeQuestStatusResult struct {
	ID     uuid.UUID
	Status quest.Status
}
