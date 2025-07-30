package quest

import (
	"time"

	"github.com/google/uuid"
	"quest-manager/internal/core/domain/model/kernel"
)

// Status represents the current state of a quest.
type Status string

const (
	StatusCreated    Status = "created"
	StatusPosted     Status = "posted"
	StatusAssigned   Status = "assigned"
	StatusInProgress Status = "in_progress"
	StatusDeclined   Status = "declined"
	StatusCompleted  Status = "completed"
)

// Difficulty represents the difficulty level of a quest.
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// Quest is the main domain model representing a quest entity.
type Quest struct {
	ID                uuid.UUID
	Title             string
	Description       string
	Difficulty        Difficulty
	Reward            string
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
	Status            Status
	Creator           string
	Assignee          *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewQuest creates a new quest instance with "created" status.
func NewQuest(
	title, description string,
	difficulty Difficulty,
	reward string,
	targetLocation, executionLocation kernel.GeoCoordinate,
	creator string,
	equipment, skills []string,
) Quest {
	now := time.Now()
	return Quest{
		ID:                uuid.New(),
		Title:             title,
		Description:       description,
		Difficulty:        difficulty,
		Reward:            reward,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Equipment:         equipment,
		Skills:            skills,
		Status:            StatusCreated,
		Creator:           creator,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// AssignTo sets the assignee for the quest and changes its status to "assigned".
func (q *Quest) AssignTo(userID string) {
	q.Assignee = &userID
	q.Status = StatusAssigned
	q.UpdatedAt = time.Now()
}

// MarkPosted updates the quest status to "posted".
func (q *Quest) MarkPosted() {
	q.Status = StatusPosted
	q.UpdatedAt = time.Now()
}

// MarkInProgress updates the quest status to "in_progress".
func (q *Quest) MarkInProgress() {
	q.Status = StatusInProgress
	q.UpdatedAt = time.Now()
}

// MarkDeclined updates the quest status to "declined".
func (q *Quest) MarkDeclined() {
	q.Status = StatusDeclined
	q.UpdatedAt = time.Now()
}

// MarkCompleted updates the quest status to "completed".
func (q *Quest) MarkCompleted() {
	q.Status = StatusCompleted
	q.UpdatedAt = time.Now()
}
