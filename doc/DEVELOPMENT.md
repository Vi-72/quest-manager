# Development Guide - Quest Manager

## 🚀 Getting Started

### Prerequisites

- **Go:** 1.21 or higher
- **PostgreSQL:** 14 or higher
- **Docker:** (optional) for containerized development
- **Make:** (optional) for build automation

### Initial Setup

1. **Clone repository**
```bash
git clone <repository-url>
cd quest-manager
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup database**
```bash
# Using Docker
docker-compose up -d postgres

# Or manual PostgreSQL
createdb quest_manager
```

4. **Configure environment**
```bash
cp config.example .env
# Edit .env with your settings
```

5. **Run application**
```bash
go run cmd/app/main.go
```

---

## 🛠️ Development Workflow

### 1. Code → Test → Build Cycle

```bash
# Make changes to code

# Run tests
go test ./...

# Check compilation
go build ./...

# Run linter
golangci-lint run

# Format code
goimports -w .
```

### 2. Feature Development

**Standard workflow:**
1. Create feature branch: `git checkout -b feature/my-feature`
2. Write domain tests first (TDD)
3. Implement domain logic
4. Write contract tests
5. Implement use cases
6. Write HTTP tests
7. Implement HTTP handlers
8. Update OpenAPI spec
9. Run all tests
10. Create PR

---

## 🧪 Testing During Development

### Quick Test (domain only)
```bash
go test ./tests/domain/... -v
```

### Full Test Suite
```bash
go test ./... -v -p 1
```

### Integration Tests (requires PostgreSQL)
```bash
go test -tags=integration ./tests/integration/... -v -p 1
```

### Specific Test
```bash
go test ./tests/domain -run TestQuest_AssignTo -v
```

### With Coverage
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📝 Code Style Guide

### Go Conventions

**File naming:**
- `snake_case.go` for file names
- One primary type per file
- Group related functions together

**Naming:**
- `PascalCase` for exported types and functions
- `camelCase` for unexported types and functions
- Interface names: `Repository`, `Handler`, `Client` (no `-er` suffix unless natural)

**Comments:**
- English only
- Document all exported functions
- Explain "why", not "what"

### Project-Specific Rules

**Error handling:**
```go
// Good: wrap errors with context
return errs.WrapInfrastructureError("failed to save quest", err)

// Bad: return raw error
return err
```

**Domain logic:**
```go
// Good: business logic in domain
func (q *Quest) AssignTo(userID uuid.UUID) error {
    if q.status != StatusCreated && q.status != StatusPosted {
        return errors.New("invalid status for assignment")
    }
    // ...
}

// Bad: business logic in handler
if quest.Status != "created" && quest.Status != "posted" {
    return errors.New("invalid status")
}
```

**Dependency injection:**
```go
// Good: inject dependencies
func NewHandler(repo Repository, publisher EventPublisher) Handler {
    return &handler{repo: repo, publisher: publisher}
}

// Bad: create dependencies inside
func NewHandler() Handler {
    repo := NewRepository() // ❌ tight coupling
    return &handler{repo: repo}
}
```

---

## 🏗️ Adding New Features

### Adding New Endpoint

**1. Update OpenAPI spec** (`api/http/quests/v1/openapi.yaml`)
```yaml
paths:
  /quests/my-endpoint:
    get:
      summary: My new endpoint
      operationId: myEndpoint
      responses:
        '200':
          description: Success
```

**2. Regenerate OpenAPI code**
```bash
go generate ./api/http/quests/v1/
```

**3. Create domain logic** (`internal/core/domain/model/quest/`)
```go
func (q *Quest) DoSomething() error {
    // Business rules here
    return nil
}
```

**4. Create use case** (`internal/core/application/usecases/`)
```go
type MyCommandHandler interface {
    Handle(ctx context.Context, cmd MyCommand) (MyResult, error)
}
```

**5. Implement handler** (`internal/adapters/in/http/`)
```go
func (a *ApiHandler) MyEndpoint(ctx context.Context, req MyRequest) (MyResponse, error) {
    // Implementation
}
```

**6. Wire in composition root** (`cmd/container.go`)
```go
myHandler := usecases.NewMyCommandHandler(unitOfWork, eventPublisher)
```

**7. Write tests**
- Domain tests
- Contract tests  
- Handler tests
- HTTP tests

---

## 🔍 Debugging

### Enable Debug Logging
```bash
export LOG_LEVEL=debug
go run cmd/app/main.go
```

### Debug with Delve
```bash
dlv debug cmd/app/main.go
```

### SQL Query Logging
```bash
export DB_LOG_LEVEL=debug
# Shows all SQL queries with parameters
```

### Inspect HTTP Requests
```bash
# Run with verbose middleware logging
go run cmd/app/main.go
# Check logs for [requestID] entries
```

---

## 📊 Code Quality Tools

### Linting
```bash
# Run golangci-lint
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### Formatting
```bash
# Format imports
goimports -w .

# Format code
gofmt -w .

# Or use both
go fmt ./... && goimports -w .
```

### Static Analysis
```bash
# Find suspicious code
go vet ./...

# Check for common mistakes
staticcheck ./...
```

### Test Coverage
```bash
# Generate coverage report
./scripts/coverage-check.sh

# View in browser
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🐛 Common Issues & Solutions

### Issue 1: "connection refused" in tests
**Cause:** PostgreSQL not running

**Solution:**
```bash
docker-compose up -d postgres
# Or start PostgreSQL manually
```

---

### Issue 2: "oapi-codegen: command not found"
**Cause:** OpenAPI code generator not installed

**Solution:**
```bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
```

---

### Issue 3: Tests failing with "transaction already in progress"
**Cause:** UnitOfWork reuse between tests

**Solution:**
- Ensure each test creates new UnitOfWork
- Use `-p 1` for integration tests
- Clean up in TearDownTest

---

### Issue 4: "user ID not found in context"
**Cause:** Missing authentication middleware

**Solution:**
- Ensure auth middleware is registered
- For tests, use mock auth client
- Check middleware order

---

## 🔄 Development Modes

### Mode 1: Full Stack (with Auth)
```bash
# Start Auth service
cd ../quest-auth && go run cmd/main.go

# Start Quest Manager
export AUTH_GRPC=localhost:50051
export MIDDLEWARE_ENABLE_AUTH=true
go run cmd/app/main.go
```

### Mode 2: No Auth (Development)
```bash
export MIDDLEWARE_ENABLE_AUTH=false
go run cmd/app/main.go
# All requests use mock authentication
```

### Mode 3: Docker Compose
```bash
docker-compose up
# Full stack with database
```

---

## 📦 Project Structure

```
quest-manager/
├── api/                      # API definitions
│   └── http/quests/v1/      # OpenAPI spec & generated code
├── cmd/                      # Application entry points
│   ├── app/main.go          # Main application
│   └── *.go                 # Composition root, config, router
├── internal/                # Private application code
│   ├── adapters/            # Adapters (in/out)
│   │   ├── in/http/        # HTTP handlers, middleware
│   │   └── out/            # PostgreSQL, gRPC clients
│   ├── core/               # Core business logic
│   │   ├── application/    # Use cases
│   │   ├── domain/         # Domain models
│   │   └── ports/          # Interfaces
│   └── pkg/                # Shared packages
│       ├── ddd/            # DDD building blocks
│       └── errs/           # Error utilities
├── tests/                   # All tests
│   ├── domain/             # Domain tests
│   ├── contracts/          # Contract tests
│   └── integration/        # Integration tests
├── doc/                     # Documentation
└── scripts/                 # Utility scripts
```

---

## 🎯 Development Best Practices

### 1. Test-Driven Development (TDD)
- Write test first (red)
- Implement minimal code (green)
- Refactor (refactor)

### 2. Domain-First Design
- Start with domain models
- Define business rules in domain
- Keep domain independent

### 3. Small Commits
- One feature per commit
- Descriptive commit messages
- Keep changes focused

### 4. Code Review
- Self-review before PR
- Run all tests locally
- Check linter output

### 5. Documentation
- Update README for major changes
- Create CHANGELOG for breaking changes
- Comment complex logic

---

## 🔧 Useful Make Commands

```bash
# Build
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Generate OpenAPI code
make generate

# Run application
make run

# Clean artifacts
make clean
```

See `Makefile` for all available commands.

---

## 📚 Learning Resources

### Go Best Practices
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Architecture Patterns
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

### Testing
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testify Suite](https://pkg.go.dev/github.com/stretchr/testify/suite)

---

## 🔗 Related

- [Architecture](ARCHITECTURE.md) - System design
- [API Documentation](API.md) - API endpoints
- [Testing](TESTING.md) - Testing strategies
- [Configuration](CONFIGURATION.md) - Environment setup

---

**Happy Coding!** 🎉

