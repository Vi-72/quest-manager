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
// Validates that difficulty is a valid domain value.
func NewQuest(
	title, description string,
	difficulty string, // Принимаем string для валидации
	reward string,
	targetLocation, executionLocation kernel.GeoCoordinate,
	creator string,
	equipment, skills []string,
) (Quest, error) {
	// Валидация difficulty в домене
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

	questID := uuid.New()
	now := time.Now()

	quest := Quest{
		BaseAggregate:     ddd.NewBaseAggregate(questID),
		ID:                questID,
		Title:             title,
		Description:       description,
		Difficulty:        questDifficulty,
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

	// Создаем доменное событие
	quest.RaiseDomainEvent(QuestCreated{
		ID:         uuid.New(),
		EventID:    uuid.New(),
		Title:      title,
		Creator:    creator,
		Difficulty: difficulty,
		Timestamp:  now,
	})

	return quest, nil
}

// AssignTo sets the assignee for the quest and changes its status to "assigned".
// Содержит бизнес-логику назначения квеста.
func (q *Quest) AssignTo(userID string) error {
	// Бизнес-правила назначения
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

	// Создаем доменные события
	q.RaiseDomainEvent(QuestAssigned{
		ID:        uuid.New(),
		EventID:   uuid.New(),
		QuestID:   q.ID,
		UserID:    userID,
		Timestamp: q.UpdatedAt,
	})

	q.RaiseDomainEvent(QuestStatusChanged{
		ID:        uuid.New(),
		EventID:   uuid.New(),
		QuestID:   q.ID,
		OldStatus: oldStatus,
		NewStatus: q.Status,
		Timestamp: q.UpdatedAt,
	})

	return nil
}

// ChangeStatus изменяет статус квеста с проверкой бизнес-правил
func (q *Quest) ChangeStatus(newStatus Status) error {
	// Валидация переходов статуса (business rules)
	if !q.isValidStatusTransition(q.Status, newStatus) {
		return errors.New("invalid status transition from " + string(q.Status) + " to " + string(newStatus))
	}

	oldStatus := q.Status
	q.Status = newStatus
	q.UpdatedAt = time.Now()

	// Создаем доменное событие
	q.RaiseDomainEvent(QuestStatusChanged{
		ID:        uuid.New(),
		EventID:   uuid.New(),
		QuestID:   q.ID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Timestamp: q.UpdatedAt,
	})

	return nil
}

// isValidStatusTransition проверяет допустимость перехода между статусами
func (q *Quest) isValidStatusTransition(from, to Status) bool {
	validTransitions := map[Status][]Status{
		StatusCreated:    {StatusPosted, StatusAssigned},
		StatusPosted:     {StatusAssigned, StatusCreated},
		StatusAssigned:   {StatusInProgress, StatusDeclined, StatusPosted},
		StatusInProgress: {StatusCompleted, StatusDeclined},
		StatusDeclined:   {StatusPosted},
		StatusCompleted:  {}, // Финальный статус
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

// Deprecated: использовать ChangeStatus вместо этих методов
func (q *Quest) MarkPosted() {
	q.ChangeStatus(StatusPosted)
}

func (q *Quest) MarkInProgress() {
	q.ChangeStatus(StatusInProgress)
}

func (q *Quest) MarkDeclined() {
	q.ChangeStatus(StatusDeclined)
}

func (q *Quest) MarkCompleted() {
	q.ChangeStatus(StatusCompleted)
}
