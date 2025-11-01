package tests

import (
	"net/http"

	"quest-manager/cmd"
	"quest-manager/internal/adapters/out/postgres"
	"quest-manager/internal/adapters/out/postgres/eventrepo"
	"quest-manager/internal/adapters/out/postgres/locationrepo"
	"quest-manager/internal/adapters/out/postgres/questrepo"
	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
	"quest-manager/tests/integration/core/storage"
	testmock "quest-manager/tests/integration/mock"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// TestContainer holds all dependencies for integration tests
type TestContainer struct {
	suite.Suite

	// Database
	DB *gorm.DB

	// Transaction Manager
	TransactionManager ports.TransactionManager

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

	// HTTP Router
	HTTPRouter http.Handler

	// Mock Auth Client
	MockAuthClient *testmock.AlwaysSuccessAuthClient

	// Test utilities
	EventStorage *storage.EventStorage
}

// NewTestContainer creates a new test container with all dependencies
func NewTestContainer(db *gorm.DB) *TestContainer {
	suiteContainer := &TestContainer{}
	suiteContainer.DB = db

	// Create transaction manager
	transactionManager := postgres.NewTransactionManager(db)

	// Create repositories
	questRepo := questrepo.NewRepository(db)
	locationRepo := locationrepo.NewRepository(db)
	eventRepo := eventrepo.NewRepository(db)

	// Create query handlers with bare repositories
	listQuestsHandler := queries.NewListQuestsQueryHandler(questRepo)
	getQuestByIDHandler := queries.NewGetQuestByIDQueryHandler(questRepo)
	searchQuestsByRadiusHandler := queries.NewSearchQuestsByRadiusQueryHandler(questRepo)
	listAssignedQuestsHandler := queries.NewListAssignedQuestsQueryHandler(questRepo)

	// Create command handlers with transaction manager
	createQuestHandler := commands.NewCreateQuestCommandHandler(transactionManager)
	assignQuestHandler := commands.NewAssignQuestCommandHandler(transactionManager)
	changeQuestStatusHandler := commands.NewChangeQuestStatusCommandHandler(transactionManager)

	// Create EventStorage for testing
	eventStorage := storage.NewEventStorage(db)

	suiteContainer.TransactionManager = transactionManager
	suiteContainer.QuestRepository = questRepo
	suiteContainer.LocationRepository = locationRepo
	suiteContainer.EventPublisher = eventRepo

	suiteContainer.CreateQuestHandler = createQuestHandler
	suiteContainer.AssignQuestHandler = assignQuestHandler
	suiteContainer.ChangeQuestStatusHandler = changeQuestStatusHandler

	suiteContainer.ListQuestsHandler = listQuestsHandler
	suiteContainer.GetQuestByIDHandler = getQuestByIDHandler
	suiteContainer.SearchQuestsByRadiusHandler = searchQuestsByRadiusHandler
	suiteContainer.ListAssignedQuestsHandler = listAssignedQuestsHandler

	suiteContainer.EventStorage = eventStorage

	// Create mock auth client
	mockAuthClient := testmock.NewAlwaysSuccessAuthClient()
	suiteContainer.MockAuthClient = mockAuthClient

	// Create HTTP router with mock auth
	configs := cmd.Config{
		Middleware: cmd.MiddlewareConfig{
			DevAuth: cmd.DevAuthConfig{
				Enabled: false, // Use real auth client (mock)
			},
		},
	}
	testContainer, _ := cmd.NewContainer(configs, db)
	testContainer.SetAuthClient(mockAuthClient)
	suiteContainer.HTTPRouter = cmd.NewRouter(testContainer)

	return suiteContainer
}

// NewHTTPRouterWithAuthClient creates a new HTTP router with a custom auth client
func (tc *TestContainer) NewHTTPRouterWithAuthClient(authClient ports.AuthClient) http.Handler {
	configs := cmd.Config{
		Middleware: cmd.MiddlewareConfig{
			DevAuth: cmd.DevAuthConfig{
				Enabled: false, // Use real auth client
			},
		},
	}
	testContainer, _ := cmd.NewContainer(configs, tc.DB)
	testContainer.SetAuthClient(authClient)
	return cmd.NewRouter(testContainer)
}

// CleanupAll clears all test data
func (tc *TestContainer) CleanupAll() {
	// Clear all tables
	tc.DB.Exec("DELETE FROM events")
	tc.DB.Exec("DELETE FROM quests")
	tc.DB.Exec("DELETE FROM locations")
}

// WaitForEventProcessing waits for async events to be processed
func (tc *TestContainer) WaitForEventProcessing(expectedCount int64) {
	// For integration tests, we can check the events table directly
	// This is a simple implementation - in real tests you might want to wait with timeout
	var count int64
	tc.DB.Table("events").Count(&count)
	// Note: This is a simplified implementation. In practice, you might want to add proper waiting logic
	_ = expectedCount
	_ = count
}
