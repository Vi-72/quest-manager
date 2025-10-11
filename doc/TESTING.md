# Testing Strategy - Quest Manager

## 🧪 Overview

Quest Manager uses a comprehensive testing strategy with multiple test layers to ensure code quality, reliability, and maintainability.

---

## 📊 Test Pyramid

```
        ╱╲
       ╱E2E╲         3 tests    - Full system integration
      ╱──────╲
     ╱  HTTP  ╲      21 tests   - API endpoint validation
    ╱──────────╲
   ╱  Handler   ╲    18 tests   - Use case orchestration
  ╱──────────────╲
 ╱   Contracts    ╲  24 tests   - Interface compliance
╱──────────────────╲
│      Domain       │ 45+ tests  - Business logic
└────────────────────┘
```

**Total:** 110+ tests covering all layers

---

## 🎯 Test Layers

### 1. Domain Tests (`tests/domain/`)
**Purpose:** Test business logic in isolation  
**Dependencies:** None (pure domain logic)  
**Speed:** Very Fast (< 1ms per test)

**What we test:**
- Quest creation and validation
- Status transitions and business rules
- Assignment logic
- Domain events generation
- Coordinate calculations
- Value object validation

**Example:**
```go
func TestQuest_AssignTo_ValidStatusBoundaries(t *testing.T) {
    // Create quest in 'created' status
    quest := createTestQuest(quest.StatusCreated)
    
    // Assign to user
    err := quest.AssignTo(userID)
    
    // Assert: should succeed and change status to 'assigned'
    assert.NoError(t, err)
    assert.Equal(t, quest.StatusAssigned, quest.Status)
}
```

**Coverage:** Business rules, value objects, entities, aggregates

---

### 2. Contract Tests (`tests/contracts/`)
**Purpose:** Verify interface compliance between layers  
**Dependencies:** Mocks and contracts  
**Speed:** Fast (< 10ms per test)

**What we test:**
- Command handlers follow expected behavior
- Query handlers return correct data structures
- Repositories implement port interfaces
- Event publisher contracts
- UnitOfWork transactions

**Example:**
```go
func TestAssignQuestCommandHandlerContract(t *testing.T) {
    // Setup mock repositories
    mockUoW := mocks.NewMockUnitOfWork(...)
    handler := commands.NewAssignQuestCommandHandler(mockUoW, nil)
    
    // Execute command
    result, err := handler.Handle(ctx, cmd)
    
    // Assert: handler follows contract
    assert.NoError(t, err)
    assert.Equal(t, questID, result.ID)
}
```

**Coverage:** Use cases, handlers, port interfaces

---

### 3. Handler Tests (`tests/integration/tests/quest_handler_tests/`)
**Purpose:** Test use case handlers with real database  
**Dependencies:** PostgreSQL, repositories  
**Speed:** Medium (~50ms per test)

**What we test:**
- Command handlers orchestration
- Query handlers data retrieval
- Transaction management
- Error handling (not found, validation errors)
- Repository integration

**Example:**
```go
func (s *Suite) TestAssignQuest() {
    // Pre-condition: create quest
    createdQuest := createRandomQuestStep(ctx, s.createHandler)
    
    // Act: assign quest via handler
    result := assignQuestStep(ctx, s.assignHandler, questID, userID)
    
    // Assert: verify in database
    quest := getQuestByIDStep(ctx, s.getHandler, questID)
    assert.Equal(t, userID, *quest.Assignee)
}
```

**Coverage:** Application layer, use cases, repositories

---

### 4. HTTP Tests (`tests/integration/tests/quest_http_tests/`)
**Purpose:** Test HTTP API layer and OpenAPI validation  
**Dependencies:** Full HTTP stack, database  
**Speed:** Medium (~100ms per test)

**What we test:**
- HTTP request/response handling
- OpenAPI schema validation
- Content-Type handling
- Error response formatting (Problem Details)
- Parameter validation

**Example:**
```go
func (s *Suite) TestCreateQuestHTTP() {
    // Act: send HTTP POST request
    req := CreateQuestHTTPRequest(validQuestData)
    resp := ExecuteHTTPRequest(ctx, s.httpRouter, req)
    
    // Assert: HTTP layer
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    quest := parseQuestResponse(resp.Body)
    assert.NotEmpty(t, quest.Id)
}
```

**Coverage:** HTTP handlers, OpenAPI validation, API contracts

---

### 5. Middleware Tests (`tests/integration/tests/quest_middlewares/`)
**Purpose:** Test middleware behavior in isolation  
**Dependencies:** Mock auth clients, HTTP router  
**Speed:** Fast (~50ms per test)

**What we test:**
- Authentication middleware
- Token validation
- User context injection
- Error responses for auth failures

**Example:**
```go
func (s *Suite) TestAllEndpointsRequireAuthentication() {
    // Test all endpoints without token
    for _, endpoint := range allEndpoints {
        resp := sendRequestWithoutToken(endpoint)
        assert.Equal(t, 401, resp.StatusCode)
    }
}
```

**Coverage:** Middleware, authentication, authorization

---

### 6. E2E Tests (`tests/integration/tests/quest_e2e_tests/`)
**Purpose:** Test complete user workflows  
**Dependencies:** Full system (HTTP, handlers, DB, events)  
**Speed:** Slow (~200ms per test)

**What we test:**
- Complete user workflows
- Cross-layer integration
- Event publishing and persistence
- Data consistency

**Example:**
```go
func (s *Suite) TestCreateAndAssignWorkflow() {
    // 1. Create via handler
    quest := createQuestViaHandler(ctx)
    
    // 2. Assign via HTTP API
    assignViaHTTP(ctx, quest.ID)
    
    // 3. Verify in database
    verifyQuestInDB(ctx, quest.ID)
    
    // 4. Verify events persisted
    verifyEventsInDB(ctx, quest.ID)
}
```

**Coverage:** Full system integration, workflows

---

## 🛠️ Testing Patterns

### Pattern 1: Pre-condition → Act → Assert

```go
func (s *Suite) TestExample() {
    // Pre-condition: setup test data
    quest := createQuest(...)
    
    // Act: perform action
    result := handler.Handle(ctx, quest.ID)
    
    // Assert: verify outcome
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Pattern 2: Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected error
    }{
        {"valid input", "test", nil},
        {"empty input", "", ErrEmpty},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := validate(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### Pattern 3: Helper Assertions

```go
// Reusable assertion helpers eliminate boilerplate
httpAssertions := NewQuestHTTPAssertions(s.Assert())
quest := httpAssertions.QuestHTTPCreatedSuccessfully(resp, err)
```

### Pattern 4: Test Data Generators

```go
// Centralized test data creation
questData := testdatagenerators.RandomCreateQuestRequest()
quest := testdatagenerators.SimpleQuestData("title", "desc", ...)
```

---

## 🔧 Test Utilities

### Assertions (`tests/integration/core/assertions/`)

Reusable assertion helpers:
- `QuestHTTPAssertions` - HTTP response validation
- `QuestFieldAssertions` - Quest field validation
- `QuestListAssertions` - List response validation
- `QuestAssignAssertions` - Assignment verification
- `QuestE2EAssertions` - E2E workflow validation

### Case Steps (`tests/integration/core/case_steps/`)

Reusable test steps:
- `CreateQuestStep()` - Create quest via handler
- `AssignQuestStep()` - Assign quest
- `ChangeQuestStatusStep()` - Change status
- `ExecuteHTTPRequest()` - Execute HTTP request
- `ListAssignedQuestsStep()` - Query assigned quests

### Test Data Generators (`tests/integration/core/test_data_generators/`)

Flexible test data builders:
- `RandomCreateQuestRequest()` - Random valid quest
- `SimpleQuestData()` - Quest with specific fields
- `WithRandom()` - Random values builder
- `WithTitle()`, `WithDifficulty()` - Field builders

### Mock Clients (`tests/integration/mock/`)

Configurable mocks for external dependencies:
- `AlwaysSuccessAuthClient` - Always returns success
- `ExpiredTokenAuthClient` - Returns token expired error
- `InvalidTokenAuthClient` - Returns invalid token error
- `ConfigurableAuthClient` - Custom behavior

---

## 🚀 Running Tests

### All Tests
```bash
go test ./... -v -p 1
```

### By Layer
```bash
# Domain tests only
go test ./tests/domain/... -v

# Contract tests only
go test ./tests/contracts/... -v

# Integration tests (requires PostgreSQL)
go test -tags=integration ./tests/integration/... -v -p 1
```

### Specific Test Suite
```bash
# HTTP tests
go test -tags=integration ./tests/integration/tests/quest_http_tests -v

# Middleware tests
go test -tags=integration ./tests/integration/tests/quest_middlewares/authentication -v

# Handler tests
go test -tags=integration ./tests/integration/tests/quest_handler_tests -v
```

### With Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🎯 Test Isolation

### Database Cleanup
Each test suite cleans database before/after tests:
```go
func (s *Suite) SetupTest() {
    s.TestDIContainer.CleanupDatabase()
}

func (s *Suite) TearDownTest() {
    s.TestDIContainer.CleanupDatabase()
}
```

### Transaction Isolation
- Each handler test runs in its own transaction
- Integration tests use `-p 1` to avoid race conditions
- Event publisher uses goroutine limits (5)

### Independent Tests
- No shared state between tests
- Each test creates its own data
- Tests can run in any order

---

## 📋 Testing Checklist

### For New Features
- [ ] Domain tests for business logic
- [ ] Contract tests for handlers
- [ ] Handler tests for orchestration
- [ ] HTTP tests for API layer
- [ ] E2E test for complete workflow
- [ ] Update existing tests if needed

### For Bug Fixes
- [ ] Reproduce bug with failing test
- [ ] Fix the bug
- [ ] Verify test now passes
- [ ] Add regression test if needed

### For Refactoring
- [ ] All existing tests still pass
- [ ] No behavior changes
- [ ] Coverage maintained or improved

---

## 🐛 Debugging Tests

### View Full Test Output
```bash
go test ./tests/domain/... -v
```

### Run Single Test
```bash
go test ./tests/domain -run TestQuest_AssignTo_EmptyUserID -v
```

### With Race Detector
```bash
go test ./... -race -v
```

### Integration Test Debugging
```bash
# Enable SQL logging
export DB_LOG_LEVEL=debug
go test -tags=integration ./tests/integration/... -v
```

---

## 📈 Test Metrics

### Current Coverage
- **Domain Layer:** ~95%
- **Application Layer:** ~90%
- **HTTP Layer:** ~85%
- **Overall:** ~88%

### Test Distribution
- Domain: 45+ tests
- Contracts: 24 tests
- Handlers: 18 tests
- HTTP: 21 tests
- Middleware: 20 tests
- E2E: 3 tests

---

## 🎓 Testing Best Practices

### 1. Test Naming
```go
// Good: describes what is being tested
func TestQuest_AssignTo_AlreadyAssignedToDifferentUser()

// Bad: generic name
func TestAssignment()
```

### 2. Clear Assertions
```go
// Good: specific message
assert.Equal(t, expected, actual, "Quest should be assigned to user A")

// Bad: no message
assert.Equal(t, expected, actual)
```

### 3. Minimal Setup
```go
// Good: only create what's needed
quest := createMinimalQuest()

// Bad: create everything
quest := createQuestWithAllFieldsAndRelations()
```

### 4. Test One Thing
```go
// Good: focused test
func TestAssignToValidStatus()

// Bad: multiple concerns
func TestAssignAndChangeStatusAndList()
```

---

## 🔗 Related

- [Architecture](ARCHITECTURE.md) - Understanding system design helps write better tests
- [API Documentation](API.md) - API contract that tests verify
- [Configuration](CONFIGURATION.md) - Test environment setup

---

**Remember:** Tests are documentation. Write tests that explain how the system works!

