package domain

// DOMAIN LAYER UNIT TESTS
// Tests for domain model business rules and validation logic

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

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
		t.Run("reward_"+strconv.Itoa(reward), func(t *testing.T) {
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
		t.Run("reward_"+strconv.Itoa(reward), func(t *testing.T) {
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
