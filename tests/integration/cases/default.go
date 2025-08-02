package cases

import (
	"quest-manager/cmd"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// DefaultSuite basic test suite for integration tests
type DefaultSuite struct {
	SuiteDIContainer
	TestDIContainer
	tx *gorm.DB
}

// NewDefault creates new DefaultSuite
func NewDefault(s suite.TestingSuite) DefaultSuite {
	return DefaultSuite{
		SuiteDIContainer: NewSuite(s),
	}
}

// SetupSuite initializes resources before running all tests in the suite
func (s *DefaultSuite) SetupSuite() {
	// Создаем тестовую базу данных если ее нет
	cmd.CreateDbIfNotExists("localhost", "5432", "postgres", "password", "quest_test", "disable")

	// Подключение к тестовой БД
	databaseURL := "postgres://postgres:password@localhost:5432/quest_test?sslmode=disable"

	db, sqlDB, err := cmd.MustConnectDB(databaseURL)
	s.Require().NoError(err, "Failed to connect to test database")

	// Создание Unit of Work
	unitOfWork, err := postgres.NewUnitOfWork(db)
	s.Require().NoError(err, "Failed to create unit of work")

	s.TestDIContainer = NewTestDIContainer(
		s.SuiteDIContainer,
		db,
		unitOfWork,
		func() {
			_ = sqlDB.Close()
		},
	)

	// Run migrations
	cmd.MustAutoMigrate(s.TestDIContainer.DB)

	err = s.TestDIContainer.CleanupDatabase()
	s.Require().NoError(err, "Failed to cleanup database")
}

// TearDownSuite cleans up resources after completing all tests in the suite
func (s *DefaultSuite) TearDownSuite() {
	s.TestDIContainer.TearDownTest()
}

// SetupTest prepares state before each test
func (s *DefaultSuite) SetupTest() {
	tx := s.TestDIContainer.DB.Begin()
	s.Require().NoError(tx.Error)

	unitOfWork, err := postgres.NewUnitOfWork(tx)
	s.Require().NoError(err)

	eventRepo, err := eventrepo.NewRepository(unitOfWork.(ports.Tracker), 5)
	s.Require().NoError(err)

	s.TestDIContainer.UnitOfWork = unitOfWork
	s.TestDIContainer.QuestRepository = unitOfWork.QuestRepository()
	s.TestDIContainer.LocationRepository = unitOfWork.LocationRepository()
	s.TestDIContainer.EventPublisher = eventRepo

	s.TestDIContainer.CreateQuestHandler = commands.NewCreateQuestCommandHandler(unitOfWork, eventRepo)
	s.TestDIContainer.AssignQuestHandler = commands.NewAssignQuestCommandHandler(unitOfWork, eventRepo)
	s.TestDIContainer.ChangeQuestStatusHandler = commands.NewChangeQuestStatusCommandHandler(unitOfWork, eventRepo)

	s.TestDIContainer.ListQuestsHandler = queries.NewListQuestsQueryHandler(s.TestDIContainer.QuestRepository)
	s.TestDIContainer.GetQuestByIDHandler = queries.NewGetQuestByIDQueryHandler(s.TestDIContainer.QuestRepository)
	s.TestDIContainer.SearchQuestsByRadiusHandler = queries.NewSearchQuestsByRadiusQueryHandler(s.TestDIContainer.QuestRepository)
	s.TestDIContainer.ListAssignedQuestsHandler = queries.NewListAssignedQuestsQueryHandler(s.TestDIContainer.QuestRepository)

	s.tx = tx
}

// TearDownTest cleans state after each test
func (s *DefaultSuite) TearDownTest() {
	// Wait for event processing completion
	s.TestDIContainer.WaitForEventProcessing(0)
	if s.tx != nil {
		_ = s.tx.Rollback()
		s.tx = nil
	}
}
