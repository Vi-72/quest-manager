package domain

// DOMAIN LAYER UNIT TESTS
// Tests for domain model business rules and validation logic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

func TestIsValidStatus(t *testing.T) {
	validStatuses := []string{"created", "posted", "assigned", "in_progress", "declined", "completed"}
	invalidStatuses := []string{"", "invalid", "CREATED", "unknown", "pending"}

	for _, status := range validStatuses {
		t.Run("valid_"+status, func(t *testing.T) {
			assert.True(t, quest.IsValidStatus(status), "Status %s should be valid", status)
		})
	}

	for _, status := range invalidStatuses {
		t.Run("invalid_"+status, func(t *testing.T) {
			assert.False(t, quest.IsValidStatus(status), "Status %s should be invalid", status)
		})
	}
}

func TestNewQuest_ValidInput(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	q, err := quest.NewQuest(
		"Test Quest",
		"Test description",
		"medium",
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"equipment1", "equipment2"},
		[]string{"skill1", "skill2"},
	)

	assert.NoError(t, err)
	assert.Equal(t, "Test Quest", q.Title)
	assert.Equal(t, "Test description", q.Description)
	assert.Equal(t, quest.DifficultyMedium, q.Difficulty)
	assert.Equal(t, 3, q.Reward)
	assert.Equal(t, 60, q.DurationMinutes)
	assert.Equal(t, targetLocation, q.TargetLocation)
	assert.Equal(t, executionLocation, q.ExecutionLocation)
	assert.Equal(t, "test-creator", q.Creator)
	assert.Equal(t, []string{"equipment1", "equipment2"}, q.Equipment)
	assert.Equal(t, []string{"skill1", "skill2"}, q.Skills)
	assert.Equal(t, quest.StatusCreated, q.Status)
	assert.Nil(t, q.Assignee)
	assert.False(t, q.CreatedAt.IsZero())
	assert.False(t, q.UpdatedAt.IsZero())
	assert.NotNil(t, q.ID())
}

func TestNewQuest_InvalidDifficulty(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	invalidDifficulties := []string{"", "invalid", "EASY", "very_hard", "1"}

	for _, difficulty := range invalidDifficulties {
		t.Run("difficulty_"+difficulty, func(t *testing.T) {
			_, err := quest.NewQuest(
				"Test Quest",
				"Test description",
				difficulty,
				3,
				60,
				targetLocation,
				executionLocation,
				"test-creator",
				[]string{},
				[]string{},
			)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid difficulty")
		})
	}
}

func TestNewQuest_InvalidReward(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	invalidRewards := []int{0, -1, 6, 10, 100}

	for _, reward := range invalidRewards {
		t.Run("reward_"+string(rune(reward+'0')), func(t *testing.T) {
			_, err := quest.NewQuest(
				"Test Quest",
				"Test description",
				"medium",
				reward,
				60,
				targetLocation,
				executionLocation,
				"test-creator",
				[]string{},
				[]string{},
			)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "reward must be between 1 and 5")
		})
	}
}

func TestNewQuest_InvalidDuration(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	tests := []struct {
		name     string
		duration int
		errorMsg string
	}{
		{"zero duration", 0, "duration must be greater than 0 minutes"},
		{"negative duration", -10, "duration must be greater than 0 minutes"},
		{"too long duration", 525601, "duration too long, maximum is 1 year"},
		{"way too long", 1000000, "duration too long, maximum is 1 year"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := quest.NewQuest(
				"Test Quest",
				"Test description",
				"medium",
				3,
				tt.duration,
				targetLocation,
				executionLocation,
				"test-creator",
				[]string{},
				[]string{},
			)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestNewQuest_ValidDifficultyBoundaries(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	validDifficulties := []struct {
		input    string
		expected quest.Difficulty
	}{
		{"easy", quest.DifficultyEasy},
		{"medium", quest.DifficultyMedium},
		{"hard", quest.DifficultyHard},
	}

	for _, tc := range validDifficulties {
		t.Run("difficulty_"+tc.input, func(t *testing.T) {
			q, err := quest.NewQuest(
				"Test Quest",
				"Test description",
				tc.input,
				3,
				60,
				targetLocation,
				executionLocation,
				"test-creator",
				[]string{},
				[]string{},
			)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, q.Difficulty)
		})
	}
}

func TestNewQuest_ValidRewardBoundaries(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	validRewards := []int{1, 2, 3, 4, 5}

	for _, reward := range validRewards {
		t.Run("reward_"+string(rune(reward+'0')), func(t *testing.T) {
			q, err := quest.NewQuest(
				"Test Quest",
				"Test description",
				"medium",
				reward,
				60,
				targetLocation,
				executionLocation,
				"test-creator",
				[]string{},
				[]string{},
			)

			assert.NoError(t, err)
			assert.Equal(t, reward, q.Reward)
		})
	}
}

func TestNewQuest_ValidDurationBoundaries(t *testing.T) {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	validDurations := []struct {
		name     string
		duration int
	}{
		{"minimum duration", 1},
		{"short duration", 30},
		{"normal duration", 60},
		{"long duration", 480},        // 8 hours
		{"very long duration", 10080}, // 1 week
		{"maximum duration", 525600},  // 1 year exactly
	}

	for _, tc := range validDurations {
		t.Run(tc.name, func(t *testing.T) {
			q, err := quest.NewQuest(
				"Test Quest",
				"Test description",
				"medium",
				3,
				tc.duration,
				targetLocation,
				executionLocation,
				"test-creator",
				[]string{},
				[]string{},
			)

			assert.NoError(t, err)
			assert.Equal(t, tc.duration, q.DurationMinutes)
		})
	}
}

func TestQuest_AssignTo_Success(t *testing.T) {
	q := createValidQuest(t)
	userID := "test-user-123"

	err := q.AssignTo(userID)

	assert.NoError(t, err)
	assert.Equal(t, quest.StatusAssigned, q.Status)
	assert.NotNil(t, q.Assignee)
	assert.Equal(t, userID, *q.Assignee)
	assert.True(t, q.UpdatedAt.After(q.CreatedAt))
}

func TestQuest_AssignTo_FromPostedStatus(t *testing.T) {
	q := createValidQuest(t)

	// Change status to posted first
	err := q.ChangeStatus(quest.StatusPosted)
	assert.NoError(t, err)

	userID := "test-user-123"
	err = q.AssignTo(userID)

	assert.NoError(t, err)
	assert.Equal(t, quest.StatusAssigned, q.Status)
	assert.Equal(t, userID, *q.Assignee)
}

func TestQuest_AssignTo_InvalidStatus(t *testing.T) {
	q := createValidQuest(t)

	// Change to in_progress status
	q.Status = quest.StatusInProgress

	err := q.AssignTo("test-user")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quest can only be assigned if status is 'created' or 'posted'")
	assert.Nil(t, q.Assignee)
}

func TestQuest_AssignTo_AlreadyAssigned(t *testing.T) {
	q := createValidQuest(t)
	firstUser := "first-user"
	secondUser := "second-user"

	// Assign to first user
	err := q.AssignTo(firstUser)
	assert.NoError(t, err)
	assert.Equal(t, quest.StatusAssigned, q.Status)

	// Try to assign to second user - should fail because status is now "assigned"
	// and quest can only be assigned from "created" or "posted" status
	err = q.AssignTo(secondUser)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quest can only be assigned if status is 'created' or 'posted'")
	assert.Equal(t, firstUser, *q.Assignee) // Should still be assigned to first user
}

func TestQuest_AssignTo_AlreadyAssignedToSameUser(t *testing.T) {
	q := createValidQuest(t)

	// Set status to posted and manually assign (simulating a quest that was assigned then changed back to posted)
	q.Status = quest.StatusPosted
	userID := "test-user"
	q.Assignee = &userID

	// Try to assign to same user - should still fail because assignee is already set
	err := q.AssignTo(userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quest is already assigned to another user")
}

func TestQuest_ChangeStatus_ValidTransitions(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus quest.Status
		toStatus   quest.Status
	}{
		{"created to posted", quest.StatusCreated, quest.StatusPosted},
		{"created to assigned", quest.StatusCreated, quest.StatusAssigned},
		{"posted to assigned", quest.StatusPosted, quest.StatusAssigned},
		{"posted to created", quest.StatusPosted, quest.StatusCreated},
		{"assigned to in_progress", quest.StatusAssigned, quest.StatusInProgress},
		{"assigned to declined", quest.StatusAssigned, quest.StatusDeclined},
		{"assigned to posted", quest.StatusAssigned, quest.StatusPosted},
		{"in_progress to completed", quest.StatusInProgress, quest.StatusCompleted},
		{"in_progress to declined", quest.StatusInProgress, quest.StatusDeclined},
		{"declined to posted", quest.StatusDeclined, quest.StatusPosted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := createValidQuest(t)
			q.Status = tt.fromStatus
			originalUpdatedAt := q.UpdatedAt

			// Small delay to ensure UpdatedAt changes
			time.Sleep(1 * time.Millisecond)

			err := q.ChangeStatus(tt.toStatus)

			assert.NoError(t, err)
			assert.Equal(t, tt.toStatus, q.Status)
			assert.True(t, q.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")
		})
	}
}

func TestQuest_ChangeStatus_InvalidTransitions(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus quest.Status
		toStatus   quest.Status
	}{
		{"created to in_progress", quest.StatusCreated, quest.StatusInProgress},
		{"created to completed", quest.StatusCreated, quest.StatusCompleted},
		{"created to declined", quest.StatusCreated, quest.StatusDeclined},
		{"posted to in_progress", quest.StatusPosted, quest.StatusInProgress},
		{"posted to completed", quest.StatusPosted, quest.StatusCompleted},
		{"posted to declined", quest.StatusPosted, quest.StatusDeclined},
		{"assigned to created", quest.StatusAssigned, quest.StatusCreated},
		{"assigned to completed", quest.StatusAssigned, quest.StatusCompleted},
		{"in_progress to created", quest.StatusInProgress, quest.StatusCreated},
		{"in_progress to posted", quest.StatusInProgress, quest.StatusPosted},
		{"in_progress to assigned", quest.StatusInProgress, quest.StatusAssigned},
		{"completed to any", quest.StatusCompleted, quest.StatusCreated},
		{"declined to assigned", quest.StatusDeclined, quest.StatusAssigned},
		{"declined to in_progress", quest.StatusDeclined, quest.StatusInProgress},
		{"declined to completed", quest.StatusDeclined, quest.StatusCompleted},
		{"declined to declined", quest.StatusDeclined, quest.StatusDeclined},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := createValidQuest(t)
			q.Status = tt.fromStatus
			originalStatus := q.Status

			err := q.ChangeStatus(tt.toStatus)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid status transition")
			assert.Equal(t, originalStatus, q.Status, "Status should not change on invalid transition")
		})
	}
}

func TestQuest_DomainEvents(t *testing.T) {
	// Test that creating quest raises QuestCreated event
	q := createValidQuest(t)

	events := q.GetDomainEvents()
	assert.Len(t, events, 1, "NewQuest should raise one domain event")

	// Clear events
	q.ClearDomainEvents()
	assert.Len(t, q.GetDomainEvents(), 0, "Events should be cleared")

	// Test that assigning quest raises events
	err := q.AssignTo("test-user")
	assert.NoError(t, err)

	events = q.GetDomainEvents()
	assert.Len(t, events, 2, "AssignTo should raise two domain events (assigned + status changed)")

	// Clear and test status change
	q.ClearDomainEvents()
	err = q.ChangeStatus(quest.StatusInProgress)
	assert.NoError(t, err)

	events = q.GetDomainEvents()
	assert.Len(t, events, 1, "ChangeStatus should raise one domain event")
}

// Helper function to create a valid quest for testing
func createValidQuest(t *testing.T) *quest.Quest {
	targetLocation := kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	executionLocation := kernel.GeoCoordinate{Lat: 55.7559, Lon: 37.6177}

	q, err := quest.NewQuest(
		"Test Quest",
		"Test description",
		"medium",
		3,
		60,
		targetLocation,
		executionLocation,
		"test-creator",
		[]string{"equipment"},
		[]string{"skill"},
	)

	assert.NoError(t, err)
	return &q
}
