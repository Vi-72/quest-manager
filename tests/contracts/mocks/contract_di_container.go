package mocks

import (
	"context"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/ddd"
)

// ContractDIContainer provides mocked dependencies for contract testing
type ContractDIContainer struct {
	// Repositories
	QuestRepository    ports.QuestRepository
	LocationRepository ports.LocationRepository
	EventPublisher     ports.EventPublisher
	TransactionManager ports.TransactionManager

	// Command Handlers
	CreateQuestHandler       commands.CreateQuestCommandHandler
	AssignQuestHandler       commands.AssignQuestCommandHandler
	ChangeQuestStatusHandler commands.ChangeQuestStatusCommandHandler

	// Query Handlers
	ListQuestsHandler           queries.ListQuestsQueryHandler
	GetQuestByIDHandler         queries.GetQuestByIDQueryHandler
	SearchQuestsByRadiusHandler queries.SearchQuestsByRadiusQueryHandler
	ListAssignedQuestsHandler   queries.ListAssignedQuestsQueryHandler
}

// NewContractDIContainer creates a new DI container with mocked dependencies
func NewContractDIContainer() *ContractDIContainer {
	// Create mock repositories
	questRepo := NewMockQuestRepository()
	locationRepo := NewMockLocationRepository()
	eventPublisher := &MockEventPublisher{}

	// Create transaction manager with mock repositories
	txManager := &MockTransactionManager{
		questRepo:    questRepo,
		locationRepo: locationRepo,
		eventRepo:    eventPublisher,
	}

	// Create command handlers with transaction manager
	createQuestHandler := commands.NewCreateQuestCommandHandler(txManager)
	assignQuestHandler := commands.NewAssignQuestCommandHandler(txManager)
	changeQuestStatusHandler := commands.NewChangeQuestStatusCommandHandler(txManager)

	// Create query handlers with bare repositories
	listQuestsHandler := queries.NewListQuestsQueryHandler(questRepo)
	getQuestByIDHandler := queries.NewGetQuestByIDQueryHandler(questRepo)
	searchQuestsByRadiusHandler := queries.NewSearchQuestsByRadiusQueryHandler(questRepo)
	listAssignedQuestsHandler := queries.NewListAssignedQuestsQueryHandler(questRepo)

	return &ContractDIContainer{
		QuestRepository:    questRepo,
		LocationRepository: locationRepo,
		EventPublisher:     eventPublisher,
		TransactionManager: txManager,

		CreateQuestHandler:       createQuestHandler,
		AssignQuestHandler:       assignQuestHandler,
		ChangeQuestStatusHandler: changeQuestStatusHandler,

		ListQuestsHandler:           listQuestsHandler,
		GetQuestByIDHandler:         getQuestByIDHandler,
		SearchQuestsByRadiusHandler: searchQuestsByRadiusHandler,
		ListAssignedQuestsHandler:   listAssignedQuestsHandler,
	}
}

// CleanupAll clears all mock repositories
func (c *ContractDIContainer) CleanupAll() {
	if mockQuestRepo, ok := c.QuestRepository.(*MockQuestRepository); ok {
		mockQuestRepo.Clear()
	}
	if mockLocationRepo, ok := c.LocationRepository.(*MockLocationRepository); ok {
		mockLocationRepo.Clear()
	}
	if mockEventPublisher, ok := c.EventPublisher.(*MockEventPublisher); ok {
		mockEventPublisher.PublishedEvents = nil
		mockEventPublisher.PublishAsyncEvents = nil
		mockEventPublisher.PublishError = nil
	}
}

// WaitForEventProcessing is a no-op for mocked implementation
func (c *ContractDIContainer) WaitForEventProcessing(expectedCount int64) {
	_ = expectedCount // unused in mock
	// No-op for mocks
}

// MockTransactionManager for testing
type MockTransactionManager struct {
	questRepo    ports.QuestRepository
	locationRepo ports.LocationRepository
	eventRepo    ports.EventPublisher
}

var _ ports.TransactionManager = &MockTransactionManager{}

func (m *MockTransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context, repos ports.Repositories) error) error {
	repos := ports.Repositories{
		Quest:    m.questRepo,
		Location: m.locationRepo,
		Event:    m.eventRepo,
	}
	return fn(ctx, repos)
}

// MockEventPublisher for testing
type MockEventPublisher struct {
	PublishedEvents    []ddd.DomainEvent
	PublishError       error
	PublishAsyncEvents []ddd.DomainEvent
}

func (m *MockEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	_ = ctx // unused in mock
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, events...)
	return nil
}

func (m *MockEventPublisher) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	_ = ctx // unused in mock
	m.PublishAsyncEvents = append(m.PublishAsyncEvents, events...)
}
