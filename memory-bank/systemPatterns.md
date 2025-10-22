# Quest Manager - System Patterns

## 🏗️ Architecture Overview

Quest Manager implements **Clean Architecture** with **Domain-Driven Design** principles, creating a maintainable, testable, and scalable system.

## 🎯 Core Architectural Patterns

### 1. Clean Architecture (Uncle Bob)

**Structure**: 4 concentric layers with dependency inversion

```
┌─────────────────────────────────────────────────────────┐
│                   Presentation Layer                     │
│                  (HTTP, Middleware)                      │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                  Application Layer                       │
│              (Use Cases, Commands, Queries)              │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                    Domain Layer                          │
│           (Business Logic, Aggregates, Events)           │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│                Infrastructure Layer                      │
│          (PostgreSQL, gRPC Clients, Repositories)        │
└─────────────────────────────────────────────────────────┘
```

**Key Principles**:
- Dependencies point **inward** (outer layers depend on inner layers)
- Domain layer has **no external dependencies**
- Business logic isolated from infrastructure concerns
- Testable without external systems

### 2. Domain-Driven Design (DDD)

#### Aggregates
- **Quest**: Central business entity with lifecycle management
- **Location**: Geospatial entity with coordinate operations
- **GeoCoordinate**: Value object for location calculations

#### Domain Events
```go
// Quest Events
QuestCreated{ID, Title, CreatedAt, ...}
QuestAssigned{QuestID, UserID, AssignedAt, ...}  
QuestStatusChanged{QuestID, OldStatus, NewStatus, ...}

// Location Events
LocationCreated{ID, Coordinate, CreatedAt, ...}
LocationUpdated{ID, Coordinate, UpdatedAt, ...}
```

#### Value Objects
- **GeoCoordinate**: Immutable location with validation
- **BoundingBox**: Spatial query optimization
- **QuestStatus**: Enum with business rules

### 3. Hexagonal Architecture (Ports & Adapters)

**Ports (Interfaces)**:
```go
type QuestRepository interface {
    Save(ctx context.Context, quest *quest.Quest) error
    GetByID(ctx context.Context, id uuid.UUID) (*quest.Quest, error)
    SearchByRadius(ctx context.Context, center GeoCoordinate, radius float64) ([]*quest.Quest, error)
}

type UnitOfWork interface {
    Begin(ctx context.Context) error
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
    QuestRepository() ports.QuestRepository
    LocationRepository() ports.LocationRepository
}

type EventPublisher interface {
    Publish(ctx context.Context, events []domain.Event) error
}
```

**Adapters (Implementations)**:
- PostgreSQL repositories
- gRPC auth client
- HTTP handlers
- Event storage

### 4. CQRS (Command Query Responsibility Segregation)

#### Commands (Write Operations)
```go
// Command Handlers
CreateQuestCommand → CreateQuestHandler
AssignQuestCommand → AssignQuestHandler
ChangeQuestStatusCommand → ChangeQuestStatusHandler

// Characteristics:
- Modify state
- Use transactions
- Publish domain events
- Return void or minimal data
```

#### Queries (Read Operations)
```go
// Query Handlers
ListQuestsQuery → ListQuestsHandler
GetQuestByIDQuery → GetQuestByIDHandler
SearchQuestsByRadiusQuery → SearchQuestsByRadiusHandler

// Characteristics:
- Read-only operations
- No transactions needed
- Optimized for performance
- Return data models
```

### 5. Event-Driven Architecture

#### Event Flow
```
Domain Operation
    ↓
Add Domain Event
    ↓
Publish Event (in transaction)
    ↓
Event Persisted
    ↓
Transaction Commits
    ↓
Events Cleared
```

#### Event Storage
- Events stored in PostgreSQL `events` table
- Transactional consistency with aggregates
- Future: Message queue integration (RabbitMQ/Kafka)

## 🔧 Implementation Patterns

### 1. Container Pattern (Dependency Injection)

**Lazy Initialization**:
```go
func (c *Container) GetAuthClient(ctx context.Context) ports.AuthClient {
    if c.authClient == nil {
        conn, err := grpc.NewClient(c.configs.AuthGRPC, ...)
        if err != nil {
            panic(fmt.Errorf("failed to create auth gRPC client: %w", err))
        }
        c.RegisterCloser(connCloser{conn})
        c.authClient = authclient.NewUserAuthClient(grpcClient)
    }
    return c.authClient
}
```

**Context-Aware Dependencies**:
```go
// All getter methods accept context.Context
func (c *Container) GetAuthConn(ctx context.Context) *grpc.ClientConn
func (c *Container) GetQuestRepository(ctx context.Context) ports.QuestRepository
func (c *Container) GetUnitOfWork(ctx context.Context) ports.UnitOfWork
```

### 2. Unit of Work Pattern

**Transactional Consistency**:
```go
func (h *CreateQuestHandler) Handle(ctx context.Context, cmd CreateQuestCommand) error {
    uow := h.container.GetUnitOfWork(ctx)
    
    if err := uow.Begin(ctx); err != nil {
        return err
    }
    defer uow.Rollback(ctx)
    
    // Business logic
    quest := quest.NewQuest(cmd.Title, cmd.Difficulty, ...)
    
    // Persist changes
    if err := uow.QuestRepository().Save(ctx, quest); err != nil {
        return err
    }
    
    // Publish events
    if err := h.eventPublisher.Publish(ctx, quest.Events()); err != nil {
        return err
    }
    
    return uow.Commit(ctx)
}
```

### 3. Repository Pattern

**Abstraction Layer**:
```go
type QuestRepository interface {
    Save(ctx context.Context, quest *quest.Quest) error
    GetByID(ctx context.Context, id uuid.UUID) (*quest.Quest, error)
    SearchByRadius(ctx context.Context, center GeoCoordinate, radius float64) ([]*quest.Quest, error)
    ListByStatus(ctx context.Context, status quest.Status) ([]*quest.Quest, error)
}

// PostgreSQL Implementation
type postgresQuestRepository struct {
    db *gorm.DB
}

func (r *postgresQuestRepository) Save(ctx context.Context, quest *quest.Quest) error {
    return r.db.WithContext(ctx).Save(quest).Error
}
```

### 4. Middleware Pattern

**HTTP Middleware Chain**:
```go
func (c *Container) Middlewares(swagger *openapi3.T) []func(http.Handler) http.Handler {
    middlewares := []func(http.Handler) http.Handler{
        middleware.Recovery,
        middleware.Logger,
    }
    
    if c.configs.Middleware.EnableAuth {
        authMW := httpmiddleware.NewAuthMiddleware(c.GetAuthClient(ctx))
        middlewares = append(middlewares, authMW.Auth)
    }
    
    if c.configs.Middleware.EnableValidation {
        validationMW := httpmiddleware.NewValidationMiddleware(swagger)
        middlewares = append(middlewares, validationMW.Validate)
    }
    
    return middlewares
}
```

## 🗄️ Data Patterns

### 1. Hybrid Location Storage

**Denormalized Coordinates**:
```sql
-- Fast access for display and search
CREATE TABLE quests (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    target_latitude DECIMAL(10,8) NOT NULL,
    target_longitude DECIMAL(11,8) NOT NULL,
    execution_latitude DECIMAL(10,8) NOT NULL,
    execution_longitude DECIMAL(11,8) NOT NULL,
    -- Optional FK to named locations
    target_location_id UUID REFERENCES locations(id),
    execution_location_id UUID REFERENCES locations(id)
);
```

**Named Locations**:
```sql
-- Reusable locations with metadata
CREATE TABLE locations (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    address TEXT,
    description TEXT
);
```

### 2. Spatial Indexing

**Optimized Queries**:
```sql
-- Spatial indexes for geospatial operations
CREATE INDEX idx_target_location ON quests(target_latitude, target_longitude);
CREATE INDEX idx_execution_location ON quests(execution_latitude, execution_longitude);
CREATE INDEX idx_locations_coords ON locations(latitude, longitude);
```

**Bounding Box + Haversine**:
```go
// 1. Bounding box pre-filter (fast)
// 2. Haversine formula for precise distance (accurate)
func (r *postgresQuestRepository) SearchByRadius(ctx context.Context, center GeoCoordinate, radius float64) ([]*quest.Quest, error) {
    // Calculate bounding box
    bbox := center.BoundingBox(radius)
    
    // Query with spatial index
    var quests []*quest.Quest
    err := r.db.WithContext(ctx).
        Where("target_latitude BETWEEN ? AND ? AND target_longitude BETWEEN ? AND ?",
            bbox.MinLat, bbox.MaxLat, bbox.MinLon, bbox.MaxLon).
        Find(&quests).Error
    
    // Filter by precise distance
    return filterByDistance(quests, center, radius), err
}
```

### 3. Event Sourcing Ready

**Event Storage**:
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    event_name VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    version INTEGER NOT NULL
);
```

## 🔐 Security Patterns

### 1. JWT Authentication

**Token Validation Flow**:
```
Client Request + JWT
    ↓
Auth Middleware
    ↓
gRPC Call → Quest Auth Service
    ↓
Validate Token
    ↓
Extract User ID
    ↓
Inject into Context
    ↓
Handler Uses User ID
```

### 2. User Context Pattern

**Security Principle**: User ID extracted from JWT, not request parameters

```go
// ❌ Insecure - user can impersonate others
func AssignQuest(questID, userID string) error

// ✅ Secure - user ID from authenticated context
func AssignQuest(ctx context.Context, questID string) error {
    userID := auth.GetUserIDFromContext(ctx)
    // Business logic with authenticated user
}
```

### 3. Input Validation

**Multi-Layer Validation**:
```go
// 1. OpenAPI Schema Validation (HTTP Layer)
validationmiddleware.Validate(r) // Format, ranges, required fields

// 2. Domain Validation (Business Layer)
quest, err := quest.NewQuest(dto.Title, dto.Difficulty, ...)
// Business rules, enum validation

// 3. Resource Validation (Application Layer)
quest, err := repository.GetByID(questID)
// Existence checks, authorization
```

## 🧪 Testing Patterns

### 1. Test Pyramid

**Unit Tests** (Domain Layer):
```go
func TestQuest_Assign(t *testing.T) {
    quest := quest.NewQuest("Test Quest", quest.DifficultyEasy, ...)
    userID := uuid.New()
    
    err := quest.Assign(userID)
    
    assert.NoError(t, err)
    assert.Equal(t, quest.Status(), quest.StatusAssigned)
    assert.Equal(t, quest.Assignee(), userID)
    assert.Len(t, quest.Events(), 1)
}
```

**Contract Tests** (Interface Layer):
```go
func TestQuestRepositoryContract(t *testing.T) {
    // Test interface compliance with mocks
    mockRepo := mocks.NewQuestRepository(t)
    // Verify interface methods work correctly
}
```

**Integration Tests** (Full Stack):
```go
func TestCreateQuestE2E(t *testing.T) {
    // Full HTTP request → database → response
    resp := httptest.NewRequest("POST", "/api/v1/quests", questData)
    assert.Equal(t, 201, resp.StatusCode)
}
```

### 2. Test Data Builders

**Domain Test Data**:
```go
func NewTestQuest() *quest.Quest {
    return quest.NewQuest(
        "Test Quest",
        quest.DifficultyEasy,
        NewTestCoordinate(),
        NewTestCoordinate(),
        []string{"equipment1"},
        []string{"skill1"},
        uuid.New(),
    )
}
```

## 🚀 Performance Patterns

### 1. Connection Pooling

**Database Connections**:
```go
func NewDatabase(config DatabaseConfig) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(dsn), &gorm.Config{
        ConnPool: &sql.DB{
            MaxOpenConns: 25,
            MaxIdleConns: 5,
            ConnMaxLifetime: 5 * time.Minute,
        },
    })
}
```

### 2. Async Event Processing

**Goroutine Pool**:
```go
type EventPublisher struct {
    workerPool chan struct{}
    eventRepo  ports.EventRepository
}

func (p *EventPublisher) Publish(ctx context.Context, events []domain.Event) error {
    select {
    case p.workerPool <- struct{}{}:
        go func() {
            defer func() { <-p.workerPool }()
            p.publishEvents(ctx, events)
        }()
        return nil
    default:
        return errors.New("event publishing queue full")
    }
}
```

---

**These patterns form the foundation of Quest Manager's architecture, ensuring maintainability, testability, and scalability while following industry best practices.**

