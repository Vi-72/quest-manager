package domain

// DOMAIN LAYER UNIT TESTS
// Additional tests for quest.AssignTo domain model business rules and validation logic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/quest"
)

func TestQuest_AssignTo_AllInvalidStatuses(t *testing.T) {
	testCases := []struct {
		name        string
		status      quest.Status
		expectedErr string
	}{
		{
			name:        "assigned status",
			status:      quest.StatusAssigned,
			expectedErr: "quest can only be assigned if status is 'created' or 'posted'",
		},
		{
			name:        "declined status",
			status:      quest.StatusDeclined,
			expectedErr: "quest can only be assigned if status is 'created' or 'posted'",
		},
		{
			name:        "completed status",
			status:      quest.StatusCompleted,
			expectedErr: "quest can only be assigned if status is 'created' or 'posted'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := createValidQuest(t)

			// Change to invalid status for assignment
			q.Status = tc.status

			err := q.AssignTo("test-user")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
			assert.Nil(t, q.Assignee)
		})
	}
}

func TestQuest_AssignTo_AlreadyAssignedToDifferentUser(t *testing.T) {
	q := createValidQuest(t)

	// Set status to posted and manually assign to first user
	q.Status = quest.StatusPosted
	firstUserID := "first-user"
	q.Assignee = &firstUserID

	// Try to assign to different user - should fail because assignee is already set
	secondUserID := "second-user"
	err := q.AssignTo(secondUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quest is already assigned to another user")
	assert.Equal(t, firstUserID, *q.Assignee) // Should still be assigned to first user
}

func TestQuest_AssignTo_ValidStatusBoundaries(t *testing.T) {
	validStatuses := []quest.Status{
		quest.StatusCreated,
		quest.StatusPosted,
	}

	for _, status := range validStatuses {
		t.Run("status_"+string(status), func(t *testing.T) {
			q := createValidQuest(t)
			q.Status = status
			q.Assignee = nil // Ensure no assignee

			userID := "test-user-" + string(status)
			err := q.AssignTo(userID)

			assert.NoError(t, err)
			assert.Equal(t, quest.StatusAssigned, q.Status)
			assert.Equal(t, userID, *q.Assignee)
		})
	}
}

func TestQuest_AssignTo_EmptyUserID(t *testing.T) {
	q := createValidQuest(t)

	// Test empty user ID
	err := q.AssignTo("")

	assert.NoError(t, err, "Domain layer doesn't validate empty user ID - that's API layer responsibility")
	assert.Equal(t, quest.StatusAssigned, q.Status)
	assert.NotNil(t, q.Assignee)
	assert.Equal(t, "", *q.Assignee)
}



func TestQuest_AssignTo_StatusTransition(t *testing.T) {
	q := createValidQuest(t)
	userID := "transition-test-user"

	// Record original status
	originalStatus := q.Status
	assert.Equal(t, quest.StatusCreated, originalStatus)

	// Act - assign quest
	err := q.AssignTo(userID)
	assert.NoError(t, err)

	// Assert - status should transition to assigned
	assert.Equal(t, quest.StatusAssigned, q.Status)
	assert.NotEqual(t, originalStatus, q.Status, "Status should have changed")
}

func TestQuest_AssignTo_TimestampUpdate(t *testing.T) {
	q := createValidQuest(t)
	userID := "timestamp-test-user"

	// Record original timestamps
	originalCreatedAt := q.CreatedAt
	originalUpdatedAt := q.UpdatedAt

	// Act - assign quest
	err := q.AssignTo(userID)
	assert.NoError(t, err)

	// Assert - timestamps should be updated correctly
	assert.Equal(t, originalCreatedAt, q.CreatedAt, "CreatedAt should not change")
	assert.True(t, q.UpdatedAt.After(originalUpdatedAt), "UpdatedAt should be updated")
	assert.True(t, q.UpdatedAt.After(q.CreatedAt), "UpdatedAt should be after CreatedAt")
}

func TestQuest_AssignTo_AssigneeFieldUpdate(t *testing.T) {
	q := createValidQuest(t)
	userID := "assignee-test-user"

	// Verify initial state
	assert.Nil(t, q.Assignee, "Initially quest should have no assignee")

	// Act - assign quest
	err := q.AssignTo(userID)
	assert.NoError(t, err)

	// Assert - assignee should be set
	assert.NotNil(t, q.Assignee, "Quest should have assignee after assignment")
	assert.Equal(t, userID, *q.Assignee, "Assignee should match provided user ID")
}

func TestQuest_AssignTo_MultipleAssignmentAttempts(t *testing.T) {
	q := createValidQuest(t)
	firstUser := "first-user"
	secondUser := "second-user"

	// First assignment should succeed
	err := q.AssignTo(firstUser)
	assert.NoError(t, err)
	assert.Equal(t, quest.StatusAssigned, q.Status)
	assert.Equal(t, firstUser, *q.Assignee)

	// Change status back to posted to test assignee validation
	q.Status = quest.StatusPosted

	// Second assignment should fail due to assignee check
	err = q.AssignTo(secondUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quest is already assigned to another user")
	assert.Equal(t, firstUser, *q.Assignee, "Assignee should not change")
}
