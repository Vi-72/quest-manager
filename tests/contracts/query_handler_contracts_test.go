package contracts

import (
	"context"
	"errors"
	"testing"

	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/pkg/errs"
	"quest-manager/tests/contracts/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// ListQuestsQueryHandlerContractSuite defines contract tests for ListQuestsQueryHandler
type ListQuestsQueryHandlerContractSuite struct {
	suite.Suite
	container *mocks.ContractDIContainer
	ctx       context.Context
	handler   queries.ListQuestsQueryHandler
}

func (s *ListQuestsQueryHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.ctx = context.Background()
	s.handler = s.container.ListQuestsHandler
}

func (s *ListQuestsQueryHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func (s *ListQuestsQueryHandlerContractSuite) TestHandleValidQuery() {
	// Create a test quest directly in the mock repository
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q, err := quest.NewQuest(
		"List Test Quest",
		"Quest for list testing",
		"easy",
		3,
		45,
		coord1,
		coord2,
		"test-creator",
		[]string{},
		[]string{},
	)
	s.Require().NoError(err)

	// Save directly to mock repository
	err = s.container.QuestRepository.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: Handler should return a list of quests without error
	result, err := s.handler.Handle(s.ctx, nil) // nil status means all quests
	s.Require().NoError(err, "Handle should succeed with valid query")

	// Contract: Result should contain the created quest
	s.Assert().NotEmpty(result, "Should return at least one quest")
	found := false
	for _, returnedQuest := range result {
		if returnedQuest.ID() == q.ID() {
			found = true
			break
		}
	}
	s.Assert().True(found, "Should include the created quest")
}

func (s *ListQuestsQueryHandlerContractSuite) TestHandleWithStatusFilter() {
	// Create a test quest directly in the mock repository
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q, err := quest.NewQuest(
		"Status Filter Test Quest",
		"Quest for status filter testing",
		"easy",
		3,
		45,
		coord1,
		coord2,
		"test-creator",
		[]string{},
		[]string{},
	)
	s.Require().NoError(err)

	// Save directly to mock repository
	err = s.container.QuestRepository.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: Handler should return only quests with the specified status
	createdStatus := quest.StatusCreated
	result, err := s.handler.Handle(s.ctx, &createdStatus)
	s.Require().NoError(err, "Handle should succeed with status filter")

	// Contract: All returned quests should have the specified status
	for _, returnedQuest := range result {
		s.Assert().Equal(createdStatus, returnedQuest.Status, "All quests should have 'created' status")
	}

	// Contract: Result should include our created quest
	found := false
	for _, returnedQuest := range result {
		if returnedQuest.ID() == q.ID() {
			found = true
			break
		}
	}
	s.Assert().True(found, "Should include the created quest with matching status")
}

// GetQuestByIDQueryHandlerContractSuite defines contract tests for GetQuestByIDQueryHandler
type GetQuestByIDQueryHandlerContractSuite struct {
	suite.Suite
	container *mocks.ContractDIContainer
	ctx       context.Context
	handler   queries.GetQuestByIDQueryHandler
}

func (s *GetQuestByIDQueryHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.ctx = context.Background()
	s.handler = s.container.GetQuestByIDHandler
}

func (s *GetQuestByIDQueryHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func (s *GetQuestByIDQueryHandlerContractSuite) TestHandleValidQuery() {
	// Create a test quest directly in the mock repository
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q, err := quest.NewQuest(
		"Get By ID Test Quest",
		"Quest for get by ID testing",
		"medium",
		4,
		45,
		coord1,
		coord2,
		"test-creator",
		[]string{},
		[]string{},
	)
	s.Require().NoError(err)

	// Save directly to mock repository
	err = s.container.QuestRepository.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: Handler should return the specific quest without error
	result, err := s.handler.Handle(s.ctx, q.ID())
	s.Require().NoError(err, "Handle should succeed with valid ID")

	// Contract: Returned quest should match the created quest
	s.Assert().Equal(q.ID(), result.ID(), "Quest ID should match")
	s.Assert().Equal(q.Title, result.Title, "Quest title should match")
	s.Assert().Equal(q.Description, result.Description, "Quest description should match")
}

func (s *GetQuestByIDQueryHandlerContractSuite) TestHandleNonExistentQuest() {
	// Contract: Handler should return not found error for non-existent quest
	nonExistentID := uuid.New()

	_, err := s.handler.Handle(s.ctx, nonExistentID)
	s.Require().Error(err, "Handle should return error for non-existent quest")
	var notFoundErr *errs.NotFoundError
	s.Assert().True(errors.As(err, &notFoundErr), "Should return not found error")
}

// SearchQuestsByRadiusQueryHandlerContractSuite defines contract tests for SearchQuestsByRadiusQueryHandler
type SearchQuestsByRadiusQueryHandlerContractSuite struct {
	suite.Suite
	container *mocks.ContractDIContainer
	ctx       context.Context
	handler   queries.SearchQuestsByRadiusQueryHandler
}

func (s *SearchQuestsByRadiusQueryHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.ctx = context.Background()
	s.handler = s.container.SearchQuestsByRadiusHandler
}

func (s *SearchQuestsByRadiusQueryHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func (s *SearchQuestsByRadiusQueryHandlerContractSuite) TestHandleValidQuery() {
	// Create a test quest at known coordinates
	coord := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}

	q, err := quest.NewQuest(
		"Radius Search Test Quest",
		"Quest for radius search testing",
		"hard",
		5,
		45,
		coord,
		coord,
		"test-creator",
		[]string{},
		[]string{},
	)
	s.Require().NoError(err)

	// Save directly to mock repository
	err = s.container.QuestRepository.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: Handler should return quests within the specified radius
	center := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0} // Same as target location
	radiusKm := 1.0                                      // 1km radius

	result, err := s.handler.Handle(s.ctx, center, radiusKm)
	s.Require().NoError(err, "Handle should succeed with valid query")

	// Contract: Result should include the created quest within radius
	found := false
	for _, returnedQuest := range result {
		if returnedQuest.ID() == q.ID() {
			found = true
			break
		}
	}
	s.Assert().True(found, "Should include the quest within radius")
}

// ListAssignedQuestsQueryHandlerContractSuite defines contract tests for ListAssignedQuestsQueryHandler
type ListAssignedQuestsQueryHandlerContractSuite struct {
	suite.Suite
	container *mocks.ContractDIContainer
	ctx       context.Context
	handler   queries.ListAssignedQuestsQueryHandler
}

func (s *ListAssignedQuestsQueryHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.ctx = context.Background()
	s.handler = s.container.ListAssignedQuestsHandler
}

func (s *ListAssignedQuestsQueryHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func (s *ListAssignedQuestsQueryHandlerContractSuite) TestHandleValidQuery() {
	// Create and assign a test quest
	coord1 := kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0}
	coord2 := kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0}

	q, err := quest.NewQuest(
		"Assigned Quest Test",
		"Quest for assigned quest testing",
		"easy",
		3,
		45,
		coord1,
		coord2,
		"test-creator",
		[]string{},
		[]string{},
	)
	s.Require().NoError(err)

	// Assign the quest
	userID := uuid.New()
	err = q.AssignTo(userID)
	s.Require().NoError(err)

	// Save directly to mock repository
	err = s.container.QuestRepository.Save(s.ctx, q)
	s.Require().NoError(err)

	// Contract: Handler should return assigned quests for the user
	result, err := s.handler.Handle(s.ctx, userID)
	s.Require().NoError(err, "Handle should succeed with valid query")

	// Contract: Result should include the assigned quest
	found := false
	for _, returnedQuest := range result {
		if returnedQuest.ID() == q.ID() {
			found = true
			s.Assert().Equal(userID, *returnedQuest.Assignee, "Quest should be assigned to the correct user")
			break
		}
	}
	s.Assert().True(found, "Should include the assigned quest")
}

func TestQueryHandlerContracts(t *testing.T) {
	suite.Run(t, new(ListQuestsQueryHandlerContractSuite))
	suite.Run(t, new(GetQuestByIDQueryHandlerContractSuite))
	suite.Run(t, new(SearchQuestsByRadiusQueryHandlerContractSuite))
	suite.Run(t, new(ListAssignedQuestsQueryHandlerContractSuite))
}
