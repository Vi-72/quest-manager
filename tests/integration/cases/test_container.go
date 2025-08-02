package cases

import (
	"net/http"
	"time"

	"quest-manager/cmd"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"

	"gorm.io/gorm"
)

// TestDIContainer содержит все зависимости для интеграционных тестов
type TestDIContainer struct {
	SuiteDIContainer
	DB         *gorm.DB
	CloseDB    func()
	UnitOfWork ports.UnitOfWork

	// Repositories
	QuestRepository    ports.QuestRepository
	LocationRepository ports.LocationRepository
	EventPublisher     ports.EventPublisher

	// Command Handlers
	CreateQuestHandler       commands.CreateQuestCommandHandler
	AssignQuestHandler       commands.AssignQuestCommandHandler
	ChangeQuestStatusHandler commands.ChangeQuestStatusCommandHandler

	// Query Handlers
	ListQuestsHandler           queries.ListQuestsQueryHandler
	GetQuestByIDHandler         queries.GetQuestByIDQueryHandler
	SearchQuestsByRadiusHandler queries.SearchQuestsByRadiusQueryHandler
	ListAssignedQuestsHandler   queries.ListAssignedQuestsQueryHandler

	// HTTP Router for API testing
	HTTPRouter http.Handler
}

// NewTestDIContainer создает новый TestDIContainer для тестов
func NewTestDIContainer(suiteContainer SuiteDIContainer) TestDIContainer {
	// Создаем тестовую базу данных если ее нет
	cmd.CreateDbIfNotExists("localhost", "5432", "postgres", "password", "quest_test", "disable")

	// Подключение к тестовой БД
	databaseURL := "postgres://postgres:password@localhost:5432/quest_test?sslmode=disable"

	db, sqlDB, err := cmd.MustConnectDB(databaseURL)
	suiteContainer.Require().NoError(err, "Failed to connect to test database")

	// Создание Unit of Work (он сам создает внутри себя quest и location репозитории)
	unitOfWork, err := postgres.NewUnitOfWork(db)
	suiteContainer.Require().NoError(err, "Failed to create unit of work")

	// Создание event репозитория отдельно
	eventRepo, err := eventrepo.NewRepository(unitOfWork.(ports.Tracker), 5) // лимит горутин = 5
	suiteContainer.Require().NoError(err, "Failed to create event repository")

	// Получаем репозитории из UnitOfWork
	questRepo := unitOfWork.QuestRepository()
	locationRepo := unitOfWork.LocationRepository()

	// Создание обработчиков команд
	createQuestHandler := commands.NewCreateQuestCommandHandler(
		unitOfWork,
		eventRepo,
	)
	assignQuestHandler := commands.NewAssignQuestCommandHandler(
		unitOfWork,
		eventRepo,
	)
	changeQuestStatusHandler := commands.NewChangeQuestStatusCommandHandler(
		unitOfWork,
		eventRepo,
	)

	// Создание обработчиков запросов
	listQuestsHandler := queries.NewListQuestsQueryHandler(questRepo)
	getQuestByIDHandler := queries.NewGetQuestByIDQueryHandler(questRepo)
	searchQuestsByRadiusHandler := queries.NewSearchQuestsByRadiusQueryHandler(questRepo)
	listAssignedQuestsHandler := queries.NewListAssignedQuestsQueryHandler(questRepo)

	// Create HTTP Router for API testing
	testConfig := cmd.Config{
		EventGoroutineLimit: 5,
	}
	compositionRoot := cmd.NewCompositionRoot(testConfig, db)
	httpRouter := cmd.NewRouter(compositionRoot)

	return TestDIContainer{
		SuiteDIContainer: suiteContainer,
		DB:               db,
		CloseDB:          func() { sqlDB.Close() },
		UnitOfWork:       unitOfWork,

		QuestRepository:    questRepo,
		LocationRepository: locationRepo,
		EventPublisher:     eventRepo,

		CreateQuestHandler:       createQuestHandler,
		AssignQuestHandler:       assignQuestHandler,
		ChangeQuestStatusHandler: changeQuestStatusHandler,

		ListQuestsHandler:           listQuestsHandler,
		GetQuestByIDHandler:         getQuestByIDHandler,
		SearchQuestsByRadiusHandler: searchQuestsByRadiusHandler,
		ListAssignedQuestsHandler:   listAssignedQuestsHandler,

		HTTPRouter: httpRouter,
	}
}

// TearDownTest очищает ресурсы после теста
func (c *TestDIContainer) TearDownTest() {
	if c.CloseDB != nil {
		c.CloseDB()
	}
}

// CleanupDatabase очищает тестовую базу данных
func (c *TestDIContainer) CleanupDatabase() error {
	// Очищаем таблицы в правильном порядке (учитывая внешние ключи)
	if err := c.DB.Exec("TRUNCATE TABLE events CASCADE").Error; err != nil {
		return err
	}
	if err := c.DB.Exec("TRUNCATE TABLE quests CASCADE").Error; err != nil {
		return err
	}
	if err := c.DB.Exec("TRUNCATE TABLE locations CASCADE").Error; err != nil {
		return err
	}
	return nil
}

// WaitForEventProcessing ждет завершения обработки событий
func (c *TestDIContainer) WaitForEventProcessing() {
	// Ждем немного для завершения обработки асинхронных событий
	time.Sleep(100 * time.Millisecond)
}
