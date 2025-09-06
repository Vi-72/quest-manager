package testdatagenerators

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/generated/servers"
)

// ============================
// RNG (сид из ENV или дефолт)
// ============================

var defaultRng *rand.Rand

func init() {
	seed := time.Now().UnixNano()
	if s, ok := os.LookupEnv("QUEST_GENERATOR_SEED"); ok {
		if parsed, err := strconv.ParseInt(s, 10, 64); err == nil {
			seed = parsed
		}
	}
	defaultRng = rand.New(rand.NewSource(seed))
}

// ============================
// Базовые константы/хелперы
// ============================

var (
	MoscowCenter = kernel.GeoCoordinate{Lat: 55.7558, Lon: 37.6176}
	MoscowNear   = kernel.GeoCoordinate{Lat: 55.7539, Lon: 37.6202}
)

func clampFloat64(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func pick[T any](r *rand.Rand, xs []T) T {
	return xs[r.Intn(len(xs))]
}

func ptr[T any](v T) *T { return &v }

// ============================
// Модель тестовых данных
// ============================

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

// ============================
// Конвертеры
// ============================

func (data QuestTestData) ToCreateCommand() commands.CreateQuestCommand {
	return commands.CreateQuestCommand{
		Title:             data.Title,
		Description:       data.Description,
		Difficulty:        data.Difficulty,
		Reward:            data.Reward,
		DurationMinutes:   data.DurationMinutes,
		TargetLocation:    data.TargetLocation,
		ExecutionLocation: data.ExecutionLocation,
		Creator:           data.Creator,
		Equipment:         data.Equipment,
		Skills:            data.Skills,
	}
}

func (data QuestTestData) ToHTTPRequest() map[string]interface{} {
	return map[string]interface{}{
		"title":              data.Title,
		"description":        data.Description,
		"difficulty":         data.Difficulty,
		"reward":             data.Reward,
		"duration_minutes":   data.DurationMinutes,
		"creator":            data.Creator,
		"target_location":    map[string]interface{}{"latitude": data.TargetLocation.Lat, "longitude": data.TargetLocation.Lon},
		"execution_location": map[string]interface{}{"latitude": data.ExecutionLocation.Lat, "longitude": data.ExecutionLocation.Lon},
		"equipment":          data.Equipment,
		"skills":             data.Skills,
	}
}

func (data QuestTestData) ToCreateQuestRequest() servers.CreateQuestRequest {
	return servers.CreateQuestRequest{
		Title:           data.Title,
		Description:     data.Description,
		Difficulty:      servers.CreateQuestRequestDifficulty(data.Difficulty),
		Reward:          data.Reward,
		DurationMinutes: data.DurationMinutes,
		TargetLocation: servers.Coordinate{
			Latitude:  float32(data.TargetLocation.Lat),
			Longitude: float32(data.TargetLocation.Lon),
		},
		ExecutionLocation: servers.Coordinate{
			Latitude:  float32(data.ExecutionLocation.Lat),
			Longitude: float32(data.ExecutionLocation.Lon),
		},
		Equipment: ptr(data.Equipment),
		Skills:    ptr(data.Skills),
	}
}

// ============================
// Option-паттерн
// ============================

type Option func(*QuestTestData, *rand.Rand)

func NewQuest(opts ...Option) QuestTestData {
	// База по умолчанию
	data := QuestTestData{
		Title:             "Test Quest",
		Description:       "A test quest for integration testing",
		Difficulty:        "medium",
		Reward:            3,
		DurationMinutes:   60,
		Creator:           "test-creator",
		TargetLocation:    MoscowCenter,
		ExecutionLocation: MoscowNear,
		Equipment:         []string{"map", "compass"},
		Skills:            []string{"navigation", "observation"},
	}
	// Прогоняем опции на дефолтном RNG
	for _, opt := range opts {
		opt(&data, defaultRng)
	}
	return data
}

// ---- Готовые опции (композиция, без дублирования):

func WithRand(rng *rand.Rand) Option {
	return func(q *QuestTestData, _ *rand.Rand) { // подменим RNG через замыкание
		// Ничего не меняем в q — эта опция полезна как "обёртка" для других WithRandom* опций
		// Используйте вместе: NewQuest(WithRand(myRng), WithRandom(...))
		defaultRng = rng
	}
}

func WithTitle(title string) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Title = title }
}

func WithDescription(desc string) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Description = desc }
}

func WithDifficulty(diff string) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Difficulty = diff }
}

func WithReward(reward int) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Reward = reward }
}

func WithDuration(minutes int) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.DurationMinutes = minutes }
}

func WithCreator(creator string) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Creator = creator }
}

func WithLocations(target, exec kernel.GeoCoordinate) Option {
	return func(q *QuestTestData, _ *rand.Rand) {
		q.TargetLocation = target
		q.ExecutionLocation = exec
	}
}

func WithEquipment(eq []string) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Equipment = eq }
}

func WithSkills(sk []string) Option {
	return func(q *QuestTestData, _ *rand.Rand) { q.Skills = sk }
}

func WithEmptyArrays() Option {
	return func(q *QuestTestData, _ *rand.Rand) {
		q.Equipment = []string{}
		q.Skills = []string{}
	}
}

func WithInvalidCoordinates() Option {
	return func(q *QuestTestData, _ *rand.Rand) {
		q.TargetLocation = kernel.GeoCoordinate{Lat: 95.0, Lon: 37.6176} // специально невалидно
	}
}

// Случайная генерация (компонуется с остальными)
func WithRandom() Option {
	difficulties := []string{"easy", "medium", "hard"}
	titles := []string{
		"Treasure Hunt", "City Explorer", "Mystery Solver",
		"Adventure Seeker", "Photo Quest", "Cultural Journey", "Food Discovery",
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

	return func(q *QuestTestData, r *rand.Rand) {
		q.Title = fmt.Sprintf("%s #%d", pick(r, titles), r.Intn(1000))
		q.Description = fmt.Sprintf("Generated test quest %d for integration testing", r.Intn(10000))
		q.Difficulty = pick(r, difficulties)
		q.Reward = r.Intn(5) + 1
		q.DurationMinutes = (r.Intn(6) + 1) * 30
		q.Equipment = pick(r, equipment)
		q.Skills = pick(r, skills)

		// Координаты около Москвы в пределах ~±0.05°
		q.TargetLocation = kernel.GeoCoordinate{
			Lat: 55.7 + (r.Float64()-0.5)*0.1,
			Lon: 37.6 + (r.Float64()-0.5)*0.1,
		}
		q.ExecutionLocation = kernel.GeoCoordinate{
			Lat: 55.7 + (r.Float64()-0.5)*0.1,
			Lon: 37.6 + (r.Float64()-0.5)*0.1,
		}
		// при желании можно clamp-нуть:
		q.TargetLocation.Lat = clampFloat64(q.TargetLocation.Lat, -90, 90)
		q.ExecutionLocation.Lat = clampFloat64(q.ExecutionLocation.Lat, -90, 90)
	}
}

// ============================
// Функции для обратной совместимости
// ============================

// DefaultQuestData возвращает стандартные данные для квеста
func DefaultQuestData() QuestTestData {
	return NewQuest()
}

// EmptyArraysQuestData возвращает данные с пустыми массивами
func EmptyArraysQuestData() QuestTestData {
	return NewQuest(WithEmptyArrays())
}

// RandomQuestData возвращает случайные данные для квеста
func RandomQuestData() QuestTestData {
	return NewQuest(WithRandom())
}

// InvalidCoordinatesQuestData возвращает данные с некорректными координатами
func InvalidCoordinatesQuestData() QuestTestData {
	return NewQuest(WithInvalidCoordinates())
}

// QuestDataWithLocations возвращает данные с заданными локациями
func QuestDataWithLocations(targetLoc, execLoc kernel.GeoCoordinate) QuestTestData {
	return NewQuest(WithLocations(targetLoc, execLoc))
}

// SimpleQuestData возвращает простые данные с заданными параметрами
func SimpleQuestData(title, description, difficulty string, reward, duration int, targetLoc, execLoc kernel.GeoCoordinate) QuestTestData {
	return NewQuest(
		WithTitle(title),
		WithDescription(description),
		WithDifficulty(difficulty),
		WithReward(reward),
		WithDuration(duration),
		WithLocations(targetLoc, execLoc),
	)
}

// ============================
// Хелперы для *servers.CreateQuestRequest
// ============================

// RandomCreateQuestRequest возвращает случайный запрос
func RandomCreateQuestRequest() *servers.CreateQuestRequest {
	req := NewQuest(WithRandom()).ToCreateQuestRequest()
	return &req
}

// InvalidCoordinatesCreateQuestRequest возвращает запрос с некорректными координатами
func InvalidCoordinatesCreateQuestRequest() *servers.CreateQuestRequest {
	req := NewQuest(WithInvalidCoordinates()).ToCreateQuestRequest()
	return &req
}

// ============================
// Хелперы для HTTP и прочего
// ============================

// HTTPQuestDataWithField создает HTTP данные с кастомным полем
func HTTPQuestDataWithField(field string, value interface{}) map[string]interface{} {
	data := NewQuest().ToHTTPRequest()
	data[field] = value
	return data
}

// ============================
// Экспорт констант для обратной совместимости
// ============================

// DefaultTestCoordinate возвращает стандартные координаты для тестов (Москва центр)
func DefaultTestCoordinate() kernel.GeoCoordinate {
	return MoscowCenter
}
