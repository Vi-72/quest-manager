package tests

import (
	"context"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"

	"quest-manager/cmd"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
	teststorage "quest-manager/tests/integration/core/storage"
	integrationmock "quest-manager/tests/integration/mock"
)

// getTestConfig возвращает конфигурацию для тестов, используя те же env переменные что и приложение
func getTestConfig() cmd.Config {
	return cmd.Config{
		DbHost:     getTestEnv("DB_HOST", "localhost"),
		DbPort:     getTestEnv("DB_PORT", "5432"),
		DbUser:     getTestEnv("DB_USER", "postgres"),
		DbPassword: getTestEnv("DB_PASSWORD", "password"),
		DbName:     getTestEnv("DB_NAME", "quest_test"),
		DbSslMode:  getTestEnv("DB_SSLMODE", "disable"),
	}
}

// getTestEnv получает environment переменную или возвращает default значение
func getTestEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestDIContainer содержит все зависимости для интеграционных тестов
type TestDIContainer struct {
	SuiteDIContainer
	DB         *gorm.DB
	CloseDB    func()
	UnitOfWork ports.UnitOfWork

	// Auth (mock for tests)
	MockAuthClient *integrationmock.AlwaysSuccessAuthClient

	// Repositories
	QuestRepository    ports.QuestRepository
	LocationRepository ports.LocationRepository
	EventPublisher     ports.EventPublisher
	EventStorage       *teststorage.EventStorage

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
	// Используем тот же подход что и в основном приложении
	testConfig := getTestConfig()

	// Создаем базу данных если ее нет (для локальной разработки и CI)
	cmd.CreateDbIfNotExists(
		testConfig.DbHost,
		testConfig.DbPort,
		testConfig.DbUser,
		testConfig.DbPassword,
		testConfig.DbName,
		testConfig.DbSslMode,
	)

	// Формируем connection string используя ту же функцию что и в приложении
	databaseURL, err := cmd.MakeConnectionString(
		testConfig.DbHost,
		testConfig.DbPort,
		testConfig.DbUser,
		testConfig.DbPassword,
		testConfig.DbName,
		testConfig.DbSslMode,
	)
	suiteContainer.Require().NoError(err, "Failed to create database connection string")

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

	// Создание EventStorage для тестирования
	eventStorage := teststorage.NewEventStorage(db)

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

	// Create Mock Auth Client for tests (always returns successful authentication)
	mockAuthClient := integrationmock.NewAlwaysSuccessAuthClient()

	// Create HTTP Router for API testing with mock auth client
	appConfig := cmd.Config{
		EventGoroutineLimit: 5,
		AuthClient:          mockAuthClient, // Inject mock for tests
	}
	compositionRoot := cmd.NewCompositionRoot(appConfig, db)
	httpRouter := cmd.NewRouter(compositionRoot)

	return TestDIContainer{
		SuiteDIContainer: suiteContainer,
		DB:               db,
		CloseDB: func() {
			err := sqlDB.Close()
			if err != nil {
				return
			}
		},
		UnitOfWork: unitOfWork,

		MockAuthClient: mockAuthClient,

		QuestRepository:    questRepo,
		LocationRepository: locationRepo,
		EventPublisher:     eventRepo,
		EventStorage:       eventStorage,

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

// WaitForEventProcessing actively waits until the expected number of events is stored.
// If expectedCount is 0, the method waits until the number of events stops changing.
// Waiting is canceled by a context with timeout to avoid hanging.
func (c *TestDIContainer) WaitForEventProcessing(expectedCount int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	storage := teststorage.NewEventStorage(c.DB)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	var lastCount int64 = -1

	for {
		select {
		case <-ctx.Done():
			c.Require().Fail("timeout waiting for events")
			return
		case <-ticker.C:
			count, err := storage.CountEvents(ctx)
			c.Require().NoError(err)

			if expectedCount > 0 {
				if count >= expectedCount {
					return
				}
			} else {
				if lastCount != -1 && count == lastCount {
					return
				}
				lastCount = count
			}
		}
	}
}
