package testdatagenerators

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

var questRand *rand.Rand

func init() {
	seed := int64(1)
	if s, ok := os.LookupEnv("QUEST_GENERATOR_SEED"); ok {
		if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
			seed = parsed
		}
	}
	questRand = rand.New(rand.NewSource(seed))
}

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

// EmptyArraysQuestData возвращает данные для квеста с пустыми массивами Equipment и Skills
func EmptyArraysQuestData() QuestTestData {
	return QuestTestData{
		Title:           "Empty Arrays Test Quest",
		Description:     "A test quest with empty equipment and skills arrays",
		Difficulty:      "easy",
		Reward:          2,
		DurationMinutes: 30,
		Creator:         "test-creator",
		TargetLocation: kernel.GeoCoordinate{
			Lat: 25.7558,
			Lon: 37.6176,
		},
		ExecutionLocation: kernel.GeoCoordinate{
			Lat: 55.7539,
			Lon: 77.4802,
		},
		Equipment: []string{}, // Empty array
		Skills:    []string{}, // Empty array
	}
}

// QuestDataWithLocations возвращает данные для квеста с заданными локациями и дефолтными названиями
func QuestDataWithLocations(targetLoc, execLoc kernel.GeoCoordinate) QuestTestData {
	data := DefaultQuestData()
	data.TargetLocation = targetLoc
	data.ExecutionLocation = execLoc
	return data
}

// SimpleQuestData возвращает данные для простого квеста с заданными параметрами
func SimpleQuestData(title, description, difficulty string, reward, duration int, targetLoc, execLoc kernel.GeoCoordinate) QuestTestData {
	data := DefaultQuestData()
	data.Title = title
	data.Description = description
	data.Difficulty = difficulty
	data.Reward = reward
	data.DurationMinutes = duration
	data.TargetLocation = targetLoc
	data.ExecutionLocation = execLoc
	return data
}

// RandomQuestData генерирует случайные данные для квеста
func RandomQuestData() *servers.CreateQuestRequest {
	r := questRand

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

// InvalidCoordinatesQuestData creates quest data with invalid coordinates for negative testing
func InvalidCoordinatesQuestData() servers.CreateQuestRequest {
	return servers.CreateQuestRequest{
		Title:           "Invalid Quest",
		Description:     "Quest with invalid coordinates",
		Difficulty:      servers.CreateQuestRequestDifficultyEasy,
		Reward:          1,
		DurationMinutes: 30,
		TargetLocation: servers.Coordinate{
			Latitude:  95.0, // Invalid: > 90
			Longitude: 37.6176,
		},
		ExecutionLocation: servers.Coordinate{
			Latitude:  55.7520,
			Longitude: 37.6175,
		},
	}
}

// ValidHTTPQuestData returns valid quest data for HTTP tests
func ValidHTTPQuestData() map[string]interface{} {
	return map[string]interface{}{
		"title":              "Valid Quest",
		"description":        "Valid description",
		"difficulty":         "easy",
		"reward":             3,
		"duration_minutes":   60,
		"target_location":    map[string]interface{}{"latitude": 55.7558, "longitude": 37.6176},
		"execution_location": map[string]interface{}{"latitude": 55.7560, "longitude": 37.6178},
	}
}

// HTTPQuestDataWithField creates HTTP test data with a specific field set to custom value
func HTTPQuestDataWithField(field string, value interface{}) map[string]interface{} {
	data := ValidHTTPQuestData()
	data[field] = value
	return data
}
