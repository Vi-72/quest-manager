package domain

// DOMAIN LAYER UNIT TESTS
// Tests for quest status validation and transitions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

func TestQuest_ChangeStatus_InvalidEnumStatus(t *testing.T) {
	tests := []struct {
		name          string
		invalidStatus quest.Status
	}{
		{"empty status", quest.Status("")},
		{"unknown status", quest.Status("unknown")},
		{"uppercase status", quest.Status("CREATED")},
		{"invalid status", quest.Status("invalid")},
		{"random string", quest.Status("random_status")},
		{"pending status", quest.Status("pending")},
		{"canceled status", quest.Status("canceled")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := createValidQuest(t)
			originalStatus := q.Status

			err := q.ChangeStatus(tt.invalidStatus)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid status:")
			assert.Contains(t, err.Error(), "is not a valid quest status")
			assert.Equal(t, originalStatus, q.Status, "Status should not change on invalid enum status")
		})
	}
}
