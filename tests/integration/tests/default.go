package tests

import (
	"os"
	"quest-manager/cmd"

	"github.com/stretchr/testify/suite"
)

// DefaultSuite basic test suite for integration tests
type DefaultSuite struct {
	SuiteDIContainer
	TestDIContainer *TestContainer
}

// NewDefault creates new DefaultSuite
func NewDefault(s suite.TestingSuite) DefaultSuite {
	return DefaultSuite{
		SuiteDIContainer: NewSuite(s),
	}
}

// getEnv returns environment variable value or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupSuite initializes resources before running all tests in the suite
func (s *DefaultSuite) SetupSuite() {
	// Get DB connection string from environment variables with fallback to defaults
	// This allows CI/CD to override values while keeping local development simple
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "quest_manager_test")
	sslMode := getEnv("DB_SSLMODE", "disable")

	// Build connection string using cmd.MakeConnectionString for consistency
	connectionString, err := cmd.MakeConnectionString(host, port, user, password, dbName, sslMode)
	s.Require().NoError(err, "Failed to build database connection string")

	db, _, err := cmd.MustConnectDB(connectionString)
	s.Require().NoError(err, "Failed to connect to test database")

	s.TestDIContainer = NewTestContainer(db)

	// Run migrations
	cmd.MustAutoMigrate(s.TestDIContainer.DB)
}

// TearDownSuite cleans up resources after completing all tests in the suite
func (s *DefaultSuite) TearDownSuite() {
	s.TestDIContainer.CleanupAll()
}

// SetupTest prepares state before each test
func (s *DefaultSuite) SetupTest() {
	// Clean database before each test
	s.TestDIContainer.CleanupAll()
}

// TearDownTest cleans state after each test
func (s *DefaultSuite) TearDownTest() {
	// Wait for event processing completion
	s.TestDIContainer.WaitForEventProcessing(0)
}
