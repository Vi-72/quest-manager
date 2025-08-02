package repository

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestAllRepositories runs all repository test suites
func TestAllRepositories(t *testing.T) {
	// Run Quest Repository tests
	suite.Run(t, &QuestRepositoryTestSuite{})

	// Run Location Repository tests
	suite.Run(t, &LocationRepositoryTestSuite{})

	// Run Event Repository tests
	suite.Run(t, &EventRepositoryTestSuite{})
}
