package commands

import (
	"github.com/google/uuid"
)

// AssignQuestCommand represents the input for assigning a quest to a user.
type AssignQuestCommand struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

// AssignQuestResult represents the output after assignment.
type AssignQuestResult struct {
	ID       uuid.UUID
	Assignee uuid.UUID
	Status   string
}
