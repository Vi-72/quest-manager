package contracts

import (
	"context"
	"errors"
	"testing"
	"time"

	"quest-manager/internal/core/application/usecases/commands"
	"quest-manager/internal/core/domain/model/kernel"
	"quest-manager/internal/core/domain/model/quest"
	"quest-manager/internal/core/ports"
	"quest-manager/internal/pkg/errs"
	"quest-manager/tests/contracts/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// CreateQuestCommandHandlerContractSuite defines contract tests for CreateQuestCommandHandler
type CreateQuestCommandHandlerContractSuite struct {
	suite.Suite
	container      *mocks.ContractDIContainer
	handler        commands.CreateQuestCommandHandler
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	ctx            context.Context
}

// AssignQuestCommandHandlerContractSuite defines contract tests for AssignQuestCommandHandler
type AssignQuestCommandHandlerContractSuite struct {
	suite.Suite
	container      *mocks.ContractDIContainer
	handler        commands.AssignQuestCommandHandler
	createHandler  commands.CreateQuestCommandHandler
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	ctx            context.Context
}

// ChangeQuestStatusCommandHandlerContractSuite defines contract tests for ChangeQuestStatusCommandHandler
type ChangeQuestStatusCommandHandlerContractSuite struct {
	suite.Suite
	container      *mocks.ContractDIContainer
	handler        commands.ChangeQuestStatusCommandHandler
	createHandler  commands.CreateQuestCommandHandler
	unitOfWork     ports.UnitOfWork
	eventPublisher ports.EventPublisher
	ctx            context.Context
}

func (s *CreateQuestCommandHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.handler = s.container.CreateQuestHandler
	s.unitOfWork = s.container.UnitOfWork
	s.eventPublisher = s.container.EventPublisher
	s.ctx = context.Background()
}

func (s *CreateQuestCommandHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func (s *AssignQuestCommandHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.handler = s.container.AssignQuestHandler
	s.createHandler = s.container.CreateQuestHandler
	s.unitOfWork = s.container.UnitOfWork
	s.eventPublisher = s.container.EventPublisher
	s.ctx = context.Background()
}

func (s *AssignQuestCommandHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func (s *ChangeQuestStatusCommandHandlerContractSuite) SetupSuite() {
	s.container = mocks.NewContractDIContainer()
	s.handler = s.container.ChangeQuestStatusHandler
	s.createHandler = s.container.CreateQuestHandler
	s.unitOfWork = s.container.UnitOfWork
	s.eventPublisher = s.container.EventPublisher
	s.ctx = context.Background()
}

func (s *ChangeQuestStatusCommandHandlerContractSuite) SetupTest() {
	// Clear all mock repositories before each test
	s.container.CleanupAll()
}

func TestCreateQuestCommandHandlerContract(t *testing.T) {
	suite.Run(t, new(CreateQuestCommandHandlerContractSuite))
}

func TestAssignQuestCommandHandlerContract(t *testing.T) {
	suite.Run(t, new(AssignQuestCommandHandlerContractSuite))
}

func TestChangeQuestStatusCommandHandlerContract(t *testing.T) {
	suite.Run(t, new(ChangeQuestStatusCommandHandlerContractSuite))
}

// CreateQuestCommandHandler contract tests

func (s *CreateQuestCommandHandlerContractSuite) TestHandleValidCommand() {
	// Contract: Handler should successfully create a quest with valid data
	targetAddr := "Target Address"
	execAddr := "Execution Address"
	cmd := commands.CreateQuestCommand{
		Title:             "Contract Test Quest",
		Description:       "Testing command handler contract",
		Difficulty:        "easy",
		Reward:            5,
		DurationMinutes:   60, // 60 minutes
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	// Contract: Handle should return a valid quest without error
	createdQuest, err := s.handler.Handle(s.ctx, cmd)
	s.Require().NoError(err, "Handle should succeed with valid command")

	// Contract: Returned quest should match command data
	s.Assert().Equal(cmd.Title, createdQuest.Title, "Quest title should match command")
	s.Assert().Equal(cmd.Description, createdQuest.Description, "Quest description should match command")
	s.Assert().Equal(quest.Difficulty(cmd.Difficulty), createdQuest.Difficulty, "Quest difficulty should match command")
	s.Assert().Equal(cmd.Reward, createdQuest.Reward, "Quest reward should match command")
	s.Assert().Equal(time.Duration(cmd.DurationMinutes)*time.Minute, time.Duration(createdQuest.DurationMinutes)*time.Minute, "Quest duration should match command")
	s.Assert().Equal(cmd.Creator, createdQuest.Creator, "Quest creator should match command")
	s.Assert().Equal(quest.StatusCreated, createdQuest.Status, "New quest should have 'created' status")
	s.Assert().NotEqual(uuid.Nil, createdQuest.ID(), "Quest should have a valid ID")
}

func (s *CreateQuestCommandHandlerContractSuite) TestHandleInvalidDifficulty() {
	// Contract: Handler should return domain validation error for invalid difficulty
	targetAddr := "Target"
	execAddr := "Execution"
	cmd := commands.CreateQuestCommand{
		Title:             "Test Quest",
		Description:       "Description",
		Difficulty:        "invalid-difficulty", // Invalid
		Reward:            5,
		DurationMinutes:   60,
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	// Contract: Handle should return validation error
	_, err := s.handler.Handle(s.ctx, cmd)
	s.Require().Error(err, "Handle should return error for invalid difficulty")
	var domainErr *errs.DomainValidationError
	s.Assert().True(errors.As(err, &domainErr), "Should return domain validation error")
}

// Note: Coordinate validation is handled at the API layer before hitting the command handler
// Command handlers receive already validated GeoCoordinate structures

// AssignQuestCommandHandler contract tests

func (s *AssignQuestCommandHandlerContractSuite) TestHandleValidAssignment() {
	// Create a quest first
	targetAddr := "Target"
	execAddr := "Execution"
	createCmd := commands.CreateQuestCommand{
		Title:             "Assignment Test Quest",
		Description:       "Quest for assignment testing",
		Difficulty:        "easy",
		Reward:            3,
		DurationMinutes:   45,
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	createdQuest, err := s.createHandler.Handle(s.ctx, createCmd)
	s.Require().NoError(err)

	// Contract: Handler should successfully assign quest to user
	assignCmd := commands.AssignQuestCommand{
		ID:     createdQuest.ID(),
		UserID: uuid.New(),
	}

	// Contract: Handle should return assignment result without error
	result, err := s.handler.Handle(s.ctx, assignCmd)
	s.Require().NoError(err, "Handle should succeed with valid assignment command")

	// Contract: Result should reflect the assignment
	s.Assert().Equal(assignCmd.ID, result.ID, "Result ID should match quest ID")
	s.Assert().Equal(assignCmd.UserID, result.Assignee, "Result assignee should match user ID")
	s.Assert().Equal(string(quest.StatusAssigned), result.Status, "Quest should be assigned status")
}

func (s *AssignQuestCommandHandlerContractSuite) TestHandleNonExistentQuest() {
	// Contract: Handler should return not found error for non-existent quest
	nonExistentID := uuid.New()
	assignCmd := commands.AssignQuestCommand{
		ID:     nonExistentID,
		UserID: uuid.New(),
	}

	// Contract: Handle should return not found error
	_, err := s.handler.Handle(s.ctx, assignCmd)
	s.Require().Error(err, "Handle should return error for non-existent quest")
	var notFoundErr *errs.NotFoundError
	s.Assert().True(errors.As(err, &notFoundErr), "Should return not found error")
}

func (s *AssignQuestCommandHandlerContractSuite) TestHandleAlreadyAssignedQuest() {
	// Create and assign a quest
	targetAddr := "Target"
	execAddr := "Execution"
	createCmd := commands.CreateQuestCommand{
		Title:             "Double Assignment Test",
		Description:       "Quest for double assignment testing",
		Difficulty:        "easy",
		Reward:            3,
		DurationMinutes:   45,
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	createdQuest, err := s.createHandler.Handle(s.ctx, createCmd)
	s.Require().NoError(err)

	// First assignment
	firstAssignCmd := commands.AssignQuestCommand{
		ID:     createdQuest.ID(),
		UserID: uuid.New(),
	}
	_, err = s.handler.Handle(s.ctx, firstAssignCmd)
	s.Require().NoError(err)

	// Contract: Handler should return domain validation error for already assigned quest
	secondAssignCmd := commands.AssignQuestCommand{
		ID:     createdQuest.ID(),
		UserID: uuid.New(),
	}

	// Contract: Handle should return validation error
	_, err = s.handler.Handle(s.ctx, secondAssignCmd)
	s.Require().Error(err, "Handle should return error for already assigned quest")
	var domainErr *errs.DomainValidationError
	s.Assert().True(errors.As(err, &domainErr), "Should return domain validation error")
}

// ChangeQuestStatusCommandHandler contract tests

func (s *ChangeQuestStatusCommandHandlerContractSuite) TestHandleValidStatusChange() {
	// Create a quest first
	targetAddr := "Target"
	execAddr := "Execution"
	createCmd := commands.CreateQuestCommand{
		Title:             "Status Change Test Quest",
		Description:       "Quest for status change testing",
		Difficulty:        "easy",
		Reward:            3,
		DurationMinutes:   45,
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	createdQuest, err := s.createHandler.Handle(s.ctx, createCmd)
	s.Require().NoError(err)

	// Contract: Handler should successfully change quest status
	statusCmd := commands.ChangeQuestStatusCommand{
		QuestID: createdQuest.ID(),
		Status:  quest.StatusPosted,
	}

	// Contract: Handle should return status change result without error
	result, err := s.handler.Handle(s.ctx, statusCmd)
	s.Require().NoError(err, "Handle should succeed with valid status change command")

	// Contract: Result should reflect the status change
	s.Assert().Equal(statusCmd.QuestID, result.ID, "Result ID should match quest ID")
	s.Assert().Equal(string(statusCmd.Status), result.Status, "Result status should match new status")
}

func (s *ChangeQuestStatusCommandHandlerContractSuite) TestHandleInvalidStatus() {
	// Create a quest first
	targetAddr := "Target"
	execAddr := "Execution"
	createCmd := commands.CreateQuestCommand{
		Title:             "Invalid Status Test Quest",
		Description:       "Quest for invalid status testing",
		Difficulty:        "easy",
		Reward:            3,
		DurationMinutes:   45,
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	createdQuest, err := s.createHandler.Handle(s.ctx, createCmd)
	s.Require().NoError(err)

	// Contract: Handler should return domain validation error for invalid status
	statusCmd := commands.ChangeQuestStatusCommand{
		QuestID: createdQuest.ID(),
		Status:  quest.Status("invalid-status"),
	}

	// Contract: Handle should return validation error
	_, err = s.handler.Handle(s.ctx, statusCmd)
	s.Require().Error(err, "Handle should return error for invalid status")
	var domainErr *errs.DomainValidationError
	s.Assert().True(errors.As(err, &domainErr), "Should return domain validation error")
}

func (s *ChangeQuestStatusCommandHandlerContractSuite) TestHandleInvalidStatusTransition() {
	// Create a quest first
	targetAddr := "Target"
	execAddr := "Execution"
	createCmd := commands.CreateQuestCommand{
		Title:             "Invalid Transition Test Quest",
		Description:       "Quest for invalid transition testing",
		Difficulty:        "easy",
		Reward:            3,
		DurationMinutes:   45,
		Creator:           "test-creator",
		TargetLocation:    kernel.GeoCoordinate{Lat: 50.0, Lon: 10.0},
		TargetAddress:     &targetAddr,
		ExecutionLocation: kernel.GeoCoordinate{Lat: 51.0, Lon: 11.0},
		ExecutionAddress:  &execAddr,
		Equipment:         []string{},
		Skills:            []string{},
	}

	createdQuest, err := s.createHandler.Handle(s.ctx, createCmd)
	s.Require().NoError(err)

	// Contract: Handler should return domain validation error for invalid status transition
	statusCmd := commands.ChangeQuestStatusCommand{
		QuestID: createdQuest.ID(),
		Status:  quest.StatusInProgress, // Can't go directly from created to in_progress
	}

	// Contract: Handle should return validation error
	_, err = s.handler.Handle(s.ctx, statusCmd)
	s.Require().Error(err, "Handle should return error for invalid status transition")
	var domainErr *errs.DomainValidationError
	s.Assert().True(errors.As(err, &domainErr), "Should return domain validation error")
}

func (s *ChangeQuestStatusCommandHandlerContractSuite) TestHandleNonExistentQuest() {
	// Contract: Handler should return not found error for non-existent quest
	nonExistentID := uuid.New()
	statusCmd := commands.ChangeQuestStatusCommand{
		QuestID: nonExistentID,
		Status:  quest.StatusPosted,
	}

	// Contract: Handle should return not found error
	_, err := s.handler.Handle(s.ctx, statusCmd)
	s.Require().Error(err, "Handle should return error for non-existent quest")
	var notFoundErr *errs.NotFoundError
	s.Assert().True(errors.As(err, &notFoundErr), "Should return not found error")
}
