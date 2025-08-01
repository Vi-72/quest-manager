package testdatagenerators

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
)

// QuestTestData содержит данные для создания тестового квеста
type QuestTestData struct {
	Title             string
	Description       string
	Difficulty        string
	Reward            int
	DurationMinutes   int
	Creator           string
	TargetLocation    kernel.GeoCoordinate
	ExecutionLocation kernel.GeoCoordinate
	Equipment         []string
	Skills            []string
}

// DefaultQuestData возвращает стандартные данные для квеста
func DefaultQuestData() QuestTestData {
	return QuestTestData{
		Title:           "Test Quest",
		Description:     "A test quest for integration testing",
		Difficulty:      "medium",
		Reward:          3,
		DurationMinutes: 60,
		Creator:         "test-creator",
		TargetLocation: kernel.GeoCoordinate{
			Lat:     55.7558,
			Lon:     37.6176,
			Address: stringPtr("Red Square, Moscow"),
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat:     55.7539,
			Lon:     37.6208,
			Address: stringPtr("GUM, Moscow"),
		},
		Equipment: []string{"map", "compass"},
		Skills:    []string{"navigation", "observation"},
	}
}

// RandomQuestData генерирует случайные данные для квеста
func RandomQuestData() QuestTestData {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	difficulties := []string{"easy", "medium", "hard"}
	titles := []string{
		"Treasure Hunt",
		"City Explorer",
		"Mystery Solver",
		"Adventure Seeker",
		"Photo Quest",
		"Cultural Journey",
		"Food Discovery",
	}

	equipment := [][]string{
		{"map", "compass"},
		{"camera", "notebook"},
		{"smartphone", "headphones"},
		{"binoculars", "magnifying glass"},
		{"flashlight", "rope"},
	}

	skills := [][]string{
		{"navigation", "observation"},
		{"photography", "research"},
		{"problem-solving", "teamwork"},
		{"attention to detail", "patience"},
		{"physical fitness", "communication"},
	}

	return QuestTestData{
		Title:           titles[r.Intn(len(titles))] + fmt.Sprintf(" #%d", r.Intn(1000)),
		Description:     fmt.Sprintf("Generated test quest %d for integration testing", r.Intn(10000)),
		Difficulty:      difficulties[r.Intn(len(difficulties))],
		Reward:          r.Intn(5) + 1,        // 1-5
		DurationMinutes: (r.Intn(6) + 1) * 30, // 30, 60, 90, 120, 150, 180 minutes
		Creator:         fmt.Sprintf("test-creator-%d", r.Intn(100)),
		TargetLocation: kernel.GeoCoordinate{
			Lat:     55.7 + (r.Float64()-0.5)*0.1, // Moscow area ±0.05°
			Lon:     37.6 + (r.Float64()-0.5)*0.1, // Moscow area ±0.05°
			Address: stringPtr(fmt.Sprintf("Test Target Location %d", r.Intn(1000))),
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat:     55.7 + (r.Float64()-0.5)*0.1, // Moscow area ±0.05°
			Lon:     37.6 + (r.Float64()-0.5)*0.1, // Moscow area ±0.05°
			Address: stringPtr(fmt.Sprintf("Test Execution Location %d", r.Intn(1000))),
		},
		Equipment: equipment[r.Intn(len(equipment))],
		Skills:    skills[r.Intn(len(skills))],
	}
}

// QuestWithSpecificDifficulty создает квест с определенной сложностью
func QuestWithSpecificDifficulty(difficulty string) QuestTestData {
	data := DefaultQuestData()
	data.Difficulty = difficulty
	data.Title = fmt.Sprintf("%s Quest", difficulty)

	// Настраиваем параметры в зависимости от сложности
	switch difficulty {
	case "easy":
		data.Reward = 1
		data.DurationMinutes = 30
	case "medium":
		data.Reward = 3
		data.DurationMinutes = 60
	case "hard":
		data.Reward = 5
		data.DurationMinutes = 120
	}

	return data
}

// QuestWithSpecificCreator создает квест для определенного создателя
func QuestWithSpecificCreator(creator string) QuestTestData {
	data := DefaultQuestData()
	data.Creator = creator
	data.Title = fmt.Sprintf("Quest by %s", creator)
	return data
}

// ValidUserID генерирует валидный UUID для пользователя
func ValidUserID() string {
	return uuid.New().String()
}

// InvalidUserID возвращает невалидный ID пользователя
func InvalidUserID() string {
	return "invalid-user-id"
}

// QuestStatuses возвращает все возможные статусы квестов
func QuestStatuses() []quest.Status {
	return []quest.Status{
		quest.StatusCreated,
		quest.StatusPosted,
		quest.StatusAssigned,
		quest.StatusInProgress,
		quest.StatusDeclined,
		quest.StatusCompleted,
	}
}

// RandomQuestStatus возвращает случайный статус квеста
func RandomQuestStatus() quest.Status {
	statuses := QuestStatuses()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return statuses[r.Intn(len(statuses))]
}

// stringPtr возвращает указатель на строку
func stringPtr(s string) *string {
	return &s
}
