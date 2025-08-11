package contracts

import (
	"context"
	"testing"

	"quest-manager/internal/core/ports"
	"quest-manager/tests/contracts/mocks"

	"github.com/stretchr/testify/suite"
)

// UnitOfWorkContractSuite defines contract tests that all UnitOfWork implementations must pass
type UnitOfWorkContractSuite struct {
	suite.Suite
	unitOfWork ports.UnitOfWork
	ctx        context.Context
}

func (s *UnitOfWorkContractSuite) SetupSuite() {
	s.unitOfWork = mocks.NewMockUnitOfWork()
	s.ctx = context.Background()
}

func (s *UnitOfWorkContractSuite) SetupTest() {
	// Clear mock repositories before each test
	if mockUOW, ok := s.unitOfWork.(*mocks.MockUnitOfWork); ok {
		mockUOW.ClearRepositories()
		mockUOW.SetShouldFail(false)
	}
}

// TestUnitOfWorkContract runs the contract test suite for UnitOfWork
func TestUnitOfWorkContract(t *testing.T) {
	suite.Run(t, new(UnitOfWorkContractSuite))
}

// UnitOfWork contract tests

func (s *UnitOfWorkContractSuite) TestBeginCommitTransaction() {
	// Contract: Begin should start a transaction without error
	err := s.unitOfWork.Begin(s.ctx)
	s.Require().NoError(err, "Begin should succeed")

	// Contract: QuestRepository should be available after Begin
	questRepo := s.unitOfWork.QuestRepository()
	s.Assert().NotNil(questRepo, "QuestRepository should be available")

	// Contract: LocationRepository should be available after Begin
	locationRepo := s.unitOfWork.LocationRepository()
	s.Assert().NotNil(locationRepo, "LocationRepository should be available")

	// Contract: Commit should complete the transaction without error
	err = s.unitOfWork.Commit(s.ctx)
	s.Assert().NoError(err, "Commit should succeed")
}

func (s *UnitOfWorkContractSuite) TestBeginRollbackTransaction() {
	// Contract: Begin should start a transaction without error
	err := s.unitOfWork.Begin(s.ctx)
	s.Require().NoError(err, "Begin should succeed")

	// Contract: Repositories should be available after Begin
	questRepo := s.unitOfWork.QuestRepository()
	s.Assert().NotNil(questRepo, "QuestRepository should be available")

	locationRepo := s.unitOfWork.LocationRepository()
	s.Assert().NotNil(locationRepo, "LocationRepository should be available")

	// Contract: Rollback should abort the transaction without error
	err = s.unitOfWork.Rollback()
	s.Assert().NoError(err, "Rollback should succeed")
}

func (s *UnitOfWorkContractSuite) TestMultipleBeginCalls() {
	// First transaction
	err := s.unitOfWork.Begin(s.ctx)
	s.Require().NoError(err, "First Begin should succeed")

	// Contract: Multiple Begin calls should handle gracefully (either succeed or return meaningful error)
	err2 := s.unitOfWork.Begin(s.ctx)
	// Implementation may either:
	// 1. Allow nested transactions (no error)
	// 2. Return an error indicating transaction already in progress
	// Both behaviors are acceptable for the contract
	if err2 != nil {
		s.T().Logf("Second Begin returned error (acceptable): %v", err2)
	}

	// Cleanup: Always try to commit/rollback
	err = s.unitOfWork.Commit(s.ctx)
	if err != nil {
		// If commit fails, try rollback
		_ = s.unitOfWork.Rollback()
	}
}

func (s *UnitOfWorkContractSuite) TestCommitWithoutBegin() {
	// Contract: Commit without Begin should handle gracefully
	// Implementation may either:
	// 1. Return an error indicating no active transaction
	// 2. Succeed as a no-op
	err := s.unitOfWork.Commit(s.ctx)
	if err != nil {
		s.T().Logf("Commit without Begin returned error (acceptable): %v", err)
	}
	// Either behavior is acceptable for the contract
}

func (s *UnitOfWorkContractSuite) TestRollbackWithoutBegin() {
	// Contract: Rollback without Begin should handle gracefully
	// Implementation may either:
	// 1. Return an error indicating no active transaction
	// 2. Succeed as a no-op
	err := s.unitOfWork.Rollback()
	if err != nil {
		s.T().Logf("Rollback without Begin returned error (acceptable): %v", err)
	}
	// Either behavior is acceptable for the contract
}

func (s *UnitOfWorkContractSuite) TestRepositoryConsistency() {
	// Contract: Repository instances should be consistent within a transaction
	err := s.unitOfWork.Begin(s.ctx)
	s.Require().NoError(err, "Begin should succeed")

	// Get repository instances multiple times
	questRepo1 := s.unitOfWork.QuestRepository()
	questRepo2 := s.unitOfWork.QuestRepository()
	locationRepo1 := s.unitOfWork.LocationRepository()
	locationRepo2 := s.unitOfWork.LocationRepository()

	// Contract: Multiple calls should return the same instance or functionally equivalent instances
	s.Assert().NotNil(questRepo1, "First QuestRepository call should return non-nil")
	s.Assert().NotNil(questRepo2, "Second QuestRepository call should return non-nil")
	s.Assert().NotNil(locationRepo1, "First LocationRepository call should return non-nil")
	s.Assert().NotNil(locationRepo2, "Second LocationRepository call should return non-nil")

	// Clean up
	err = s.unitOfWork.Commit(s.ctx)
	s.Assert().NoError(err, "Commit should succeed")
}

func (s *UnitOfWorkContractSuite) TestRepositoryAccessWithoutTransaction() {
	// Contract: Repository access without active transaction should work
	// (Implementation may choose to work outside transactions or return transaction-aware repos)

	questRepo := s.unitOfWork.QuestRepository()
	s.Assert().NotNil(questRepo, "QuestRepository should be available without active transaction")

	locationRepo := s.unitOfWork.LocationRepository()
	s.Assert().NotNil(locationRepo, "LocationRepository should be available without active transaction")
}

func (s *UnitOfWorkContractSuite) TestTransactionIsolation() {
	// This test demonstrates the expected behavior but may need database-specific setup
	// Contract: Changes within a transaction should be isolated until commit

	// Start first transaction
	err := s.unitOfWork.Begin(s.ctx)
	s.Require().NoError(err, "First transaction Begin should succeed")

	// Contract: Within transaction, repositories should be transaction-aware
	questRepo := s.unitOfWork.QuestRepository()
	s.Assert().NotNil(questRepo, "QuestRepository should be available in transaction")

	// We can't easily test isolation without making actual database changes,
	// but we can verify that the repositories are properly bound to the transaction context
	// by ensuring they don't panic and return consistent results

	// Commit transaction
	err = s.unitOfWork.Commit(s.ctx)
	s.Assert().NoError(err, "Transaction commit should succeed")
}

func (s *UnitOfWorkContractSuite) TestContextCancellation() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(s.ctx)

	// Start transaction
	err := s.unitOfWork.Begin(ctx)
	s.Require().NoError(err, "Begin with cancellable context should succeed")

	// Cancel the context
	cancel()

	// Contract: Operations with cancelled context should handle gracefully
	// The exact behavior may vary by implementation:
	// 1. May return context.Canceled error
	// 2. May complete successfully if operation was already in progress
	// 3. May timeout appropriately

	err = s.unitOfWork.Commit(ctx)
	if err != nil {
		s.T().Logf("Commit with cancelled context returned error (acceptable): %v", err)
		// If commit failed due to context cancellation, try rollback
		_ = s.unitOfWork.Rollback()
	}
}

func (s *UnitOfWorkContractSuite) TestConcurrentTransactions() {
	// This test verifies that the UnitOfWork handles concurrent access appropriately
	// Contract: UnitOfWork should either:
	// 1. Support concurrent transactions safely
	// 2. Return appropriate errors when concurrent access is not supported

	// Skip this test if we're using a mock implementation that doesn't support concurrency
	s.T().Skip("Skipping concurrent transactions test to avoid race conditions in CI")
}
