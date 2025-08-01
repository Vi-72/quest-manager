package cases

import (
	"quest-manager/cmd"

	"github.com/stretchr/testify/suite"
)

// DefaultSuite базовый тестовый набор для интеграционных тестов
type DefaultSuite struct {
	SuiteDIContainer
	TestDIContainer
}

// NewDefault создает новый DefaultSuite
func NewDefault(s suite.TestingSuite) DefaultSuite {
	return DefaultSuite{
		SuiteDIContainer: NewSuite(s),
	}
}

// SetupSuite инициализирует ресурсы перед запуском всех тестов в наборе
func (s *DefaultSuite) SetupSuite() {
	s.TestDIContainer = NewTestDIContainer(s.SuiteDIContainer)

	// Выполняем миграции
	cmd.MustAutoMigrate(s.TestDIContainer.DB)
}

// TearDownSuite очищает ресурсы после завершения всех тестов в наборе
func (s *DefaultSuite) TearDownSuite() {
	s.TestDIContainer.TearDownTest()
}

// SetupTest подготавливает состояние перед каждым тестом
func (s *DefaultSuite) SetupTest() {
	// Очищаем базу данных перед каждым тестом
	err := s.TestDIContainer.CleanupDatabase()
	s.Require().NoError(err, "Failed to cleanup database")

	// Пересоздаем TestDIContainer для каждого теста чтобы избежать проблем с транзакциями
	s.TestDIContainer = NewTestDIContainer(s.SuiteDIContainer)
}

// TearDownTest очищает состояние после каждого теста
func (s *DefaultSuite) TearDownTest() {
	// Ждем завершения обработки событий
	s.TestDIContainer.WaitForEventProcessing()
}
