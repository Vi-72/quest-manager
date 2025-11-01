# System Components - Quest Manager

## 🏗️ Architecture Overview

Quest Manager follows **Clean Architecture** and **Domain-Driven Design (DDD)** principles with clear separation of concerns.

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer (Adapters)                │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐  │
│  │  Handlers   │  │  Middleware  │  │  Error Handle  │  │
│  └─────────────┘  └──────────────┘  └────────────────┘  │
└───────────────────────────┬─────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────┐
│              Application Layer (Use Cases)              │
│  ┌─────────────────┐         ┌─────────────────┐        │
│  │    Commands     │         │     Queries     │        │
│  │  (Write ops)    │         │   (Read ops)    │        │
│  └─────────────────┘         └─────────────────┘        │
└───────────────────────────┬─────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────┐
│                 Domain Layer (Business Logic)           │
│  ┌──────────┐  ┌──────────┐  ┌─────────────────────┐    │
│  │  Quest   │  │ Location │  │  Domain Events      │    │
│  │Aggregate │  │  Entity  │  │  (quest.created,    │    │
│  └──────────┘  └──────────┘  │   quest.assigned)   │    │
└───────────────────────────────┴─────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────┐
│           Infrastructure Layer (Adapters Out)           │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐  │
│  │ PostgreSQL  │  │  Auth gRPC   │  │  Event Store   │  │
│  │Repositories │  │    Client    │  │                │  │
│  └─────────────┘  └──────────────┘  └────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## 📦 Core Components

### 1. Domain Layer (`internal/core/domain/`)

#### Quest Aggregate (`model/quest/`)
**Purpose:** Quest business logic and invariants

**Key Files:**
- `quest.go` - Quest aggregate root
- `status.go` - Status enum and transitions

**Responsibilities:**
- Validate quest creation
- Enforce status transition rules
- Handle quest assignment logic
- Generate domain events
- Maintain business invariants

**Example:**
```go
type Quest struct {
    id          uuid.UUID
    title       string
    description string
    status      Status
    assignee    *uuid.UUID
    // ... more fields
}

func (q *Quest) AssignTo(userID uuid.UUID) error {
    // Business rule: can only assign in 'created' or 'posted' status
    if q.status != StatusCreated && q.status != StatusPosted {
        return errors.New("quest cannot be assigned in current status")
    }
    
    q.assignee = &userID
    q.status = StatusAssigned
    q.AddDomainEvent(NewQuestAssignedEvent(...))
    return nil
}
```

---

#### Location Entity (`model/location/`)
**Purpose:** Geographic location with coordinates

**Key Files:**
- `location.go` - Location entity
- `events.go` - Location domain events

**Responsibilities:**
- Store coordinates (latitude, longitude)
- Optional address information
- Generate location events
- Coordinate validation

---

#### Kernel (`model/kernel/`)
**Purpose:** Shared value objects

**Key Files:**
- `geo_coordinate.go` - Geographic coordinates

**Responsibilities:**
- Validate coordinate ranges
- Calculate distances (Haversine)
- Calculate bounding boxes for radius searches

**Example:**
```go
func (c GeoCoordinate) DistanceTo(other GeoCoordinate) float64 {
    // Haversine formula implementation
    return distance
}

func (c GeoCoordinate) BoundingBoxForRadius(radiusKm float64) BoundingBox {
    // Calculate bounding box for efficient DB queries
    return box
}
```

---

### 2. Application Layer (`internal/core/application/`)

#### Commands (`usecases/commands/`)
**Purpose:** Write operations that modify state

**Key Handlers:**
- `CreateQuestCommandHandler` - Create new quest
- `AssignQuestCommandHandler` - Assign quest to user
- `ChangeQuestStatusCommandHandler` - Change quest status

**Pattern:**
```go
func (h *handler) Handle(ctx context.Context, cmd Command) (Result, error) {
    var result Result
    
    err := h.txManager.RunInTransaction(ctx, func(ctx context.Context, repos ports.Repositories) error {
        // 1. Load aggregate
        aggregate, err := repos.Quest.GetByID(ctx, cmd.ID)
        if err != nil {
            return err
        }
        
        // 2. Execute domain logic
        if err := aggregate.DoSomething(cmd.Params); err != nil {
            return err
        }
        
        // 3. Save aggregate
        if err := repos.Quest.Save(ctx, aggregate); err != nil {
            return err
        }
        
        // 4. Publish events synchronously in same transaction
        if err := repos.Event.Publish(ctx, aggregate.GetDomainEvents()...); err != nil {
            return err
        }
        
        result = aggregate
        return nil
    })
    
    return result, err
}
```

---

#### Queries (`usecases/queries/`)
**Purpose:** Read operations that don't modify state

**Key Handlers:**
- `ListQuestsQueryHandler` - List all quests
- `GetQuestByIDQueryHandler` - Get single quest
- `SearchQuestsByRadiusQueryHandler` - Geographic search
- `ListAssignedQuestsQueryHandler` - User's assigned quests

**Pattern:**
```go
func (h *handler) Handle(ctx context.Context, query Query) ([]Result, error) {
    // Simple read from repository
    return h.repository.FindByCondition(ctx, query.Filter)
}
```

**CQRS:** Commands and Queries are separated for clarity and scalability.

---

### 3. Ports Layer (`internal/core/ports/`)

**Purpose:** Interfaces for inbound and outbound adapters

**Key Interfaces:**
- `QuestRepository` - Quest persistence
- `LocationRepository` - Location persistence
- `TransactionManager` - Transaction management (closure-based)
- `EventPublisher` - Event publishing (synchronous only)
- `AuthClient` - Authentication service

**Hexagonal Architecture:** Domain depends on ports, not implementations.

---

### 4. HTTP Adapters (`internal/adapters/in/http/`)

#### API Handler (`api_handler.go`)
**Purpose:** Coordinate HTTP requests to use case handlers

**Responsibilities:**
- Implement OpenAPI server interface
- Convert HTTP requests to commands/queries
- Convert domain results to HTTP responses
- Handle errors and return Problem Details

#### Handlers (per endpoint)
- `create_quest_handler.go` - POST /quests
- `assign_quest_handler.go` - POST /quests/{id}/assign
- `list_assigned_quests_handler.go` - GET /quests/assigned
- `get_quest_by_id_handler.go` - GET /quests/{id}
- `list_quests_handler.go` - GET /quests
- `change_quest_status_handler.go` - PATCH /quests/{id}/status
- `search_quests_by_radius_handler.go` - GET /quests/search-radius

**Pattern:** One handler per endpoint for maintainability.

---

#### Middleware (`middleware/`)

**Authentication Middleware** (`authentication.go`)
- Extract Bearer token from header
- Validate token via Auth service
- Inject user ID into context
- Return 401 on auth failure

**OpenAPI Validation Middleware** (`openapi_validation.go`)
- Validate request against OpenAPI schema
- Check required fields
- Validate field formats and ranges
- Return 400 on validation failure

**Order:** Authentication → Validation → Handler

---

#### Error Handling (`errors/`)

**Problem Details (RFC 7807):**
- `problem_details.go` - Problem Details structure
- `bad_request.go` - 400 errors
- `not_found.go` - 404 errors
- `conflict.go` - 409 errors
- `converters.go` - Convert domain errors to HTTP

**Pattern:**
```go
problem := NewProblem(
    http.StatusBadRequest,
    "Validation Error",
    "Invalid difficulty: must be easy, medium, or hard",
)
problem.WriteResponse(w)
```

---

### 5. Infrastructure Adapters (`internal/adapters/out/`)

#### PostgreSQL Repositories (`postgres/`)

**Quest Repository** (`questrepo/`)
- CRUD operations for quests
- Complex queries (by status, assignee, bounding box)
- Join with locations for addresses
- Transaction support

**Location Repository** (`locationrepo/`)
- CRUD operations for locations
- Geographic queries
- Coordinate precision handling

**Event Repository** (`eventrepo/`)
- Persist domain events synchronously
- Events stored within same transaction as domain changes
- Outbox pattern for event reliability

**Transaction Manager** (`transaction_manager.go`)
- Closure-based transaction management
- GORM transaction lifecycle
- Repository instances scoped to transaction

---

#### Auth Client (`client/auth/`)

**Purpose:** Integration with Quest Auth service

**Key Files:**
- `client.go` - gRPC client wrapper
- `factory.go` - Client factory

**Responsibilities:**
- Validate JWT tokens
- Extract user ID from tokens
- Handle auth errors (expired, invalid)
- gRPC communication

---

## 🎯 Component Interactions

### Create Quest Flow
```
HTTP Request
    ↓
HTTP Handler (api_handler.go)
    ↓
CreateQuestCommandHandler
    ↓
Quest Aggregate (domain logic)
    ↓
Quest Repository (persistence)
    ↓
Event Publisher (events)
    ↓
HTTP Response
```

### Assign Quest Flow (with Auth)
```
HTTP Request + JWT Token
    ↓
Auth Middleware (validate token, extract user ID)
    ↓
HTTP Handler (get user ID from context)
    ↓
AssignQuestCommandHandler
    ↓
Quest Aggregate (business rules)
    ↓
Quest Repository (save)
    ↓
Event Publisher (quest.assigned event)
    ↓
HTTP Response
```

### Query Flow
```
HTTP Request
    ↓
Auth Middleware
    ↓
HTTP Handler
    ↓
QueryHandler (no business logic)
    ↓
Repository (direct read)
    ↓
HTTP Response
```

---

## 🧩 Design Patterns

### 1. Repository Pattern
Abstraction over data persistence.

### 2. Unit of Work Pattern
Manage transactions across multiple repositories.

### 3. CQRS (Command Query Responsibility Segregation)
Separate read and write operations.

### 4. Domain Events Pattern
Communicate changes within the system.

### 5. Dependency Injection
All dependencies injected via constructor.

### 6. Factory Pattern
Create complex objects (UnitOfWork, Handlers).

### 7. Middleware Chain
Process requests through multiple layers.

---

## 📚 Component Dependencies

### Domain Layer
**Depends on:** Nothing (pure business logic)  
**Depended by:** Application layer

### Application Layer
**Depends on:** Domain, Ports  
**Depended by:** HTTP Adapters

### Ports Layer
**Depends on:** Domain types  
**Depended by:** Application, Infrastructure

### HTTP Adapters
**Depends on:** Application, Ports  
**Depended by:** HTTP Router

### Infrastructure Adapters
**Depends on:** Ports  
**Depended by:** Composition Root

**Rule:** Dependencies point inward (Dependency Inversion Principle)

---

## 🔄 Lifecycle Management

### Application Startup
```go
1. Load configuration
2. Connect to PostgreSQL
3. Run migrations
4. Connect to Auth gRPC service
5. Create repositories
6. Create handlers
7. Create HTTP router with middlewares
8. Start HTTP server
```

### Request Lifecycle
```go
1. HTTP request arrives
2. Authentication middleware (validate token)
3. OpenAPI validation middleware
4. Route to handler
5. Handler creates command/query
6. Use case handler executes logic
7. Domain logic executed
8. Repository persists changes
9. Events published
10. Response returned
```

### Shutdown
```go
1. Stop accepting new requests
2. Wait for in-flight requests
3. Close database connections
4. Close gRPC connections
5. Graceful exit
```

---

## 🔗 Related

- [Architecture](ARCHITECTURE.md) - Overall system design
- [API Documentation](API.md) - API endpoints
- [Testing](TESTING.md) - Component testing strategies

---

**Key Principle:** Each component has a single responsibility and clear boundaries.

