package quest

import (
	"errors"
	"time"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/pkg/ddd"

	"github.com/google/uuid"
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

// IsValidStatus checks if string is a valid quest status
func IsValidStatus(status string) bool {
	switch Status(status) {
	case StatusCreated, StatusPosted, StatusAssigned, StatusInProgress, StatusDeclined, StatusCompleted:
		return true
	default:
		return false
	}
}

// Difficulty represents the difficulty level of a quest.
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// Quest is the main domain aggregate representing a quest entity.
type Quest struct {
	*ddd.BaseAggregate[uuid.UUID]
	Title           string
	Description     string
	Difficulty      Difficulty
	Reward          int
	DurationMinutes int

	// Main coordinates (denormalized for performance)
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate

	// Optional references to location directory
	TargetLocationID    *uuid.UUID
	ExecutionLocationID *uuid.UUID

	// Optional addresses from location directory
	TargetAddress    *string
	ExecutionAddress *string

	Equipment []string
	Skills    []string
	Status    Status
	Creator   string
	Assignee  *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewQuest creates a new quest instance with "created" status.
// Validates that difficulty is a valid domain value.
func NewQuest(
	title, description string,
	difficulty string, // Accept string for validation
	reward int,
	durationMinutes int,
	targetLocation, executionLocation kernel.GeoCoordinate,
	creator string,
	equipment, skills []string,
) (Quest, error) {
	// Validate difficulty in domain
	var questDifficulty Difficulty
	switch difficulty {
	case string(DifficultyEasy):
		questDifficulty = DifficultyEasy
	case string(DifficultyMedium):
		questDifficulty = DifficultyMedium
	case string(DifficultyHard):
		questDifficulty = DifficultyHard
	default:
		return Quest{}, errors.New("invalid difficulty: must be one of 'easy', 'medium', 'hard'")
	}

	// Validate reward
	if reward < 1 || reward > 5 {
		return Quest{}, errors.New("reward must be between 1 and 5")
	}

	// Validate duration
	if durationMinutes <= 0 {
		return Quest{}, errors.New("duration must be greater than 0 minutes")
	}
	if durationMinutes > 525600 { // 1 year in minutes (365 * 24 * 60)
		return Quest{}, errors.New("duration too long, maximum is 1 year (525600 minutes)")
	}

	questID := uuid.New()
	now := time.Now()

	quest := Quest{
		BaseAggregate:     ddd.NewBaseAggregate(questID),
		Title:             title,
		Description:       description,
		Difficulty:        questDifficulty,
		Reward:            reward,
		DurationMinutes:   durationMinutes,
		TargetLocation:    targetLocation,
		ExecutionLocation: executionLocation,
		Equipment:         equipment,
		Skills:            skills,
		Status:            StatusCreated,
		Creator:           creator,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Create domain event
	quest.RaiseDomainEvent(NewQuestCreated(questID, creator))

	return quest, nil
}

// AssignTo sets the assignee for the quest and changes its status to "assigned".
// Contains business logic for quest assignment.
func (q *Quest) AssignTo(userID uuid.UUID) error {
	// Business rules for assignment
	if q.Status != StatusCreated && q.Status != StatusPosted {
		return errors.New("quest can only be assigned if status is 'created' or 'posted'")
	}

	if q.Assignee != nil {
		return errors.New("quest is already assigned to another user")
	}

	oldStatus := q.Status
	q.Assignee = &userID
	q.Status = StatusAssigned
	q.UpdatedAt = time.Now()

	// Create domain events
	q.RaiseDomainEvent(NewQuestAssigned(q.ID(), userID))
	q.RaiseDomainEvent(NewQuestStatusChanged(q.ID(), oldStatus, q.Status))

	return nil
}

// ChangeStatus changes quest status with business rules validation
func (q *Quest) ChangeStatus(newStatus Status) error {
	// Validate that the new status is a valid enum value
	if !IsValidStatus(string(newStatus)) {
		return errors.New("invalid status: " + string(newStatus) + " is not a valid quest status")
	}

	// Validate status transitions (business rules)
	if !q.isValidStatusTransition(q.Status, newStatus) {
		return errors.New("invalid status transition from " + string(q.Status) + " to " + string(newStatus))
	}

	oldStatus := q.Status
	q.Status = newStatus
	q.UpdatedAt = time.Now()

	// Create domain event
	q.RaiseDomainEvent(NewQuestStatusChanged(q.ID(), oldStatus, newStatus))

	return nil
}

// isValidStatusTransition checks validity of transition between statuses
func (q *Quest) isValidStatusTransition(from, to Status) bool {
	validTransitions := map[Status][]Status{
		StatusCreated:    {StatusPosted, StatusAssigned},
		StatusPosted:     {StatusAssigned, StatusCreated},
		StatusAssigned:   {StatusInProgress, StatusDeclined, StatusPosted},
		StatusInProgress: {StatusCompleted, StatusDeclined},
		StatusDeclined:   {StatusPosted},
		StatusCompleted:  {}, // Final status
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == to {
			return true
		}
	}
	return false
}
