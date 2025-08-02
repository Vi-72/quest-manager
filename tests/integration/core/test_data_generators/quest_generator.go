package testdatagenerators

import (
	"fmt"
	"math/rand"
	"time"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
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

// DefaultTestCoordinate returns default coordinates for testing (Moscow center)
func DefaultTestCoordinate() kernel.GeoCoordinate {
	return kernel.GeoCoordinate{
		Lat: 55.7558, // Moscow latitude
		Lon: 37.6176, // Moscow longitude
	}
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
			Lat: 25.7558,
			Lon: 37.6176,
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat: 55.7539,
			Lon: 77.4802,
		},
		Equipment: []string{"map", "compass"},
		Skills:    []string{"navigation", "observation"},
	}
}

// RandomQuestData генерирует случайные данные для квеста
func RandomQuestData() *servers.CreateQuestRequest {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	difficulties := []servers.CreateQuestRequestDifficulty{
		servers.CreateQuestRequestDifficultyEasy,
		servers.CreateQuestRequestDifficultyMedium,
		servers.CreateQuestRequestDifficultyHard,
	}

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

	selectedEquipment := equipment[r.Intn(len(equipment))]
	selectedSkills := skills[r.Intn(len(skills))]

	// Generate addresses for locations
	targetAddress := "Test Target Address: Moscow Center"
	executionAddress := "Test Execution Address: Moscow Suburbs"

	return &servers.CreateQuestRequest{
		Title:           titles[r.Intn(len(titles))] + fmt.Sprintf(" #%d", r.Intn(1000)),
		Description:     fmt.Sprintf("Generated test quest %d for integration testing", r.Intn(10000)),
		Difficulty:      difficulties[r.Intn(len(difficulties))],
		Reward:          r.Intn(5) + 1,        // 1-5
		DurationMinutes: (r.Intn(6) + 1) * 30, // 30, 60, 90, 120, 150, 180 minutes
		TargetLocation: servers.Coordinate{
			Address:   &targetAddress,
			Latitude:  float32(55.7 + (r.Float64()-0.5)*0.1), // Moscow area ±0.05°
			Longitude: float32(37.6 + (r.Float64()-0.5)*0.1), // Moscow area ±0.05°
		},
		ExecutionLocation: servers.Coordinate{
			Address:   &executionAddress,
			Latitude:  float32(65.7 + (r.Float64()-0.5)*0.1), // Different area ±0.05°
			Longitude: float32(47.6 + (r.Float64()-0.5)*0.1), // Different area ±0.05°
		},
		Equipment: &selectedEquipment, // Pointer to slice
		Skills:    &selectedSkills,    // Pointer to slice
	}
}
