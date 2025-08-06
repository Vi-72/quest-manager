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
	UnitOfWork         ports.UnitOfWork

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
	// Create mocked repositories
	questRepo := NewMockQuestRepository()
	locationRepo := NewMockLocationRepository()
	eventPublisher := &MockEventPublisher{}
	unitOfWork := NewMockUnitOfWork()

	// Create command handlers with mocked dependencies
	createQuestHandler := commands.NewCreateQuestCommandHandler(unitOfWork, eventPublisher)
	assignQuestHandler := commands.NewAssignQuestCommandHandler(unitOfWork, eventPublisher)
	changeQuestStatusHandler := commands.NewChangeQuestStatusCommandHandler(unitOfWork, eventPublisher)

	// Create query handlers with mocked dependencies
	listQuestsHandler := queries.NewListQuestsQueryHandler(questRepo)
	getQuestByIDHandler := queries.NewGetQuestByIDQueryHandler(questRepo)
	searchQuestsByRadiusHandler := queries.NewSearchQuestsByRadiusQueryHandler(questRepo)
	listAssignedQuestsHandler := queries.NewListAssignedQuestsQueryHandler(questRepo)

	return &ContractDIContainer{
		QuestRepository:    questRepo,
		LocationRepository: locationRepo,
		EventPublisher:     eventPublisher,
		UnitOfWork:         unitOfWork,

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
	if mockUnitOfWork, ok := c.UnitOfWork.(*MockUnitOfWork); ok {
		mockUnitOfWork.ClearRepositories()
		mockUnitOfWork.SetShouldFail(false)
	}
}

// WaitForEventProcessing is a no-op for mocked implementation
func (c *ContractDIContainer) WaitForEventProcessing(expectedCount int64) {
	// No-op for mocks
}

// MockEventPublisher for testing (moved from event_publisher_contracts_test.go)
type MockEventPublisher struct {
	PublishedEvents    []ddd.DomainEvent
	PublishError       error
	PublishAsyncEvents []ddd.DomainEvent
}

func (m *MockEventPublisher) Publish(ctx context.Context, events ...ddd.DomainEvent) error {
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, events...)
	return nil
}

func (m *MockEventPublisher) PublishAsync(ctx context.Context, events ...ddd.DomainEvent) {
	m.PublishAsyncEvents = append(m.PublishAsyncEvents, events...)
}
