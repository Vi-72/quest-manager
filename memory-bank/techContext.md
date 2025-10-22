# Quest Manager - Technical Context

## 🛠️ Technology Stack

### Core Technologies
- **Language**: Go 1.23+
- **Database**: PostgreSQL 14+ with spatial extensions
- **Authentication**: External gRPC Auth Service
- **API**: REST with OpenAPI 3.0 specification
- **ORM**: GORM v2 for database operations

### Key Dependencies
```go
// Core Framework
github.com/go-chi/chi/v5          // HTTP router
github.com/gofrs/uuid/v5          // UUID generation
github.com/oapi-codegen/oapi-codegen // OpenAPI code generation

// Database & ORM
gorm.io/gorm                      // ORM
gorm.io/driver/postgres           // PostgreSQL driver

// gRPC & Authentication
google.golang.org/grpc            // gRPC client
github.com/Vi-72/quest-auth       // Auth service SDK

// Validation & Error Handling
github.com/getkin/kin-openapi/openapi3 // OpenAPI validation
github.com/pkg/errors             // Error wrapping

// Testing
github.com/stretchr/testify       // Test assertions
github.com/DATA-DOG/go-sqlmock    // SQL mocking
```

## 🏗️ Development Setup

### Prerequisites
```bash
# Required Software
Go 1.23+
PostgreSQL 14+
Docker & Docker Compose
Make

# Optional (for development)
golangci-lint
oapi-codegen
```

### Environment Configuration
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=quest_manager
DB_SSL_MODE=disable

# HTTP Server
HTTP_PORT=8080

# Authentication
AUTH_GRPC=localhost:50051
ENABLE_AUTH_MIDDLEWARE=true

# Development Mode
DEV_AUTH_HEADER_NAME=X-Dev-User-ID
DEV_AUTH_STATIC_USER_ID=00000000-0000-0000-0000-000000000001

# Event Processing
EVENT_GOROUTINE_LIMIT=10
```

### Quick Start
```bash
# 1. Clone and setup
git clone <repository>
cd quest-manager
cp config.example .env

# 2. Start dependencies
docker compose up -d postgres

# 3. Run application
go run ./cmd/app

# 4. Verify
curl http://localhost:8080/health
```

## 🔧 Build & Development Tools

### Make Commands
```bash
# Development
make run              # Run application
make build            # Build binary
make clean            # Clean build artifacts

# Testing
make test             # Run all tests
make test-unit        # Unit tests only
make test-integration # Integration tests
make test-coverage    # Coverage report
make test-stats       # Test statistics

# Code Quality
make lint             # Run golangci-lint
make fmt              # Format code
make generate         # Generate OpenAPI code

# Database
make db-migrate       # Run migrations
make db-seed          # Seed test data
```

### Code Generation
```bash
# OpenAPI Code Generation
make generate
# or
oapi-codegen -config configs/server.cfg.yaml api/http/quests/v1/openapi.yaml

# Generated Files
api/http/quests/v1/servers.gen.go    # HTTP server
internal/generated/                  # Generated models
```

## 🗄️ Database Schema

### Core Tables
```sql
-- Quests table with denormalized coordinates
CREATE TABLE quests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    status VARCHAR(20) NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'posted', 'assigned', 'in_progress', 'declined', 'completed')),
    
    -- Denormalized coordinates for fast access
    target_latitude DECIMAL(10,8) NOT NULL CHECK (target_latitude BETWEEN -90 AND 90),
    target_longitude DECIMAL(11,8) NOT NULL CHECK (target_longitude BETWEEN -180 AND 180),
    execution_latitude DECIMAL(10,8) NOT NULL CHECK (execution_latitude BETWEEN -90 AND 90),
    execution_longitude DECIMAL(11,8) NOT NULL CHECK (execution_longitude BETWEEN -180 AND 180),
    
    -- Optional named location references
    target_location_id UUID REFERENCES locations(id),
    execution_location_id UUID REFERENCES locations(id),
    
    -- User references
    creator UUID NOT NULL,
    assignee UUID,
    
    -- Equipment and skills
    equipment TEXT[],
    skills TEXT[],
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    assigned_at TIMESTAMP,
    completed_at TIMESTAMP
);

-- Named locations for reuse
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    latitude DECIMAL(10,8) NOT NULL CHECK (latitude BETWEEN -90 AND 90),
    longitude DECIMAL(11,8) NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    address TEXT,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Domain events for audit and integration
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_name VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);
```

### Performance Indexes
```sql
-- Quest indexes
CREATE INDEX idx_quests_status ON quests(status);
CREATE INDEX idx_quests_creator ON quests(creator);
CREATE INDEX idx_quests_assignee ON quests(assignee);
CREATE INDEX idx_quests_target_location ON quests(target_latitude, target_longitude);
CREATE INDEX idx_quests_execution_location ON quests(execution_latitude, execution_longitude);

-- Location indexes
CREATE INDEX idx_locations_coords ON locations(latitude, longitude);
CREATE INDEX idx_locations_name ON locations(name);

-- Event indexes
CREATE INDEX idx_events_aggregate ON events(aggregate_id, aggregate_type);
CREATE INDEX idx_events_name ON events(event_name);
CREATE INDEX idx_events_occurred_at ON events(occurred_at);
```

## 🔐 Security Configuration

### JWT Authentication
```go
// Production Configuration
ENABLE_AUTH_MIDDLEWARE=true
AUTH_GRPC=localhost:50051

// Development Configuration
ENABLE_AUTH_MIDDLEWARE=false
DEV_AUTH_HEADER_NAME=X-Dev-User-ID
DEV_AUTH_STATIC_USER_ID=00000000-0000-0000-0000-000000000001
```

### gRPC Client Configuration
```go
type AuthClientConfig struct {
    GRPCAddress string
    Timeout     time.Duration
    Retries     int
}

func NewAuthClient(config AuthClientConfig) ports.AuthClient {
    conn, err := grpc.Dial(config.GRPCAddress, 
        grpc.WithInsecure(), // Use TLS in production
        grpc.WithTimeout(config.Timeout),
    )
    if err != nil {
        panic(fmt.Errorf("failed to connect to auth service: %w", err))
    }
    
    return authclient.NewUserAuthClient(conn)
}
```

## 🧪 Testing Infrastructure

### Test Database Setup
```bash
# Docker Compose for testing
version: '3.8'
services:
  postgres-test:
    image: postgres:14
    environment:
      POSTGRES_DB: quest_manager_test
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    ports:
      - "5433:5432"
    volumes:
      - postgres_test_data:/var/lib/postgresql/data

volumes:
  postgres_test_data:
```

### Test Configuration
```go
// Test environment variables
func setupTestConfig() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5433,
            User:     "postgres",
            Password: "secret",
            Name:     "quest_manager_test",
            SSLMode:  "disable",
        },
        Auth: AuthConfig{
            GRPC: "localhost:50051", // Mock auth service
        },
        Middleware: MiddlewareConfig{
            EnableAuth:       false, // Use mock auth
            EnableValidation: true,
            EnableLogging:    false,
            EnableRecovery:   true,
        },
    }
}
```

### Test Categories
```bash
# Unit Tests (Domain Layer)
go test ./tests/domain -v

# Contract Tests (Interface Layer)
go test ./tests/contracts -v

# Integration Tests (Full Stack)
go test -tags=integration ./tests/integration/... -v -p 1

# HTTP Tests (API Layer)
go test -tags=integration ./tests/integration/tests/quest_http_tests -v

# E2E Tests (Complete Flow)
go test -tags=integration ./tests/integration/tests/quest_e2e_tests -v
```

## 📊 Monitoring & Observability

### Health Checks
```go
// Health endpoint
func (h *APIHandler) Health(w http.ResponseWriter, r *http.Request) {
    health := map[string]string{
        "status":    "healthy",
        "timestamp": time.Now().UTC().Format(time.RFC3339),
        "version":   "1.5.0",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}
```

### Logging Configuration
```go
// Structured logging
type Logger struct {
    *logrus.Logger
}

func NewLogger() *Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })
    logger.SetLevel(logrus.InfoLevel)
    return &Logger{logger}
}
```

### Error Handling
```go
// Problem Details (RFC 7807)
type Problem struct {
    Type     string `json:"type"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail"`
    Instance string `json:"instance,omitempty"`
}

func (p *Problem) WriteResponse(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/problem+json")
    w.WriteHeader(p.Status)
    json.NewEncoder(w).Encode(p)
}
```

## 🚀 Deployment Configuration

### Docker Configuration
```dockerfile
# Multi-stage build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o quest-manager ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/quest-manager .
COPY --from=builder /app/config.example .
EXPOSE 8080
CMD ["./quest-manager"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  quest-manager:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=secret
      - DB_NAME=quest_manager
      - AUTH_GRPC=auth-service:50051
    depends_on:
      - postgres
      - auth-service

  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: quest_manager
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  auth-service:
    image: quest-auth:latest
    ports:
      - "50051:50051"

volumes:
  postgres_data:
```

## 🔧 Development Constraints

### Code Style
- **Go Version**: 1.23+ required
- **Linting**: golangci-lint with strict rules
- **Formatting**: goimports for import organization
- **Comments**: English language, godoc format

### Performance Requirements
- **Response Time**: <100ms for typical operations
- **Memory Usage**: <100MB per instance
- **Database Connections**: Max 25 concurrent connections
- **Event Processing**: Async with goroutine pool (max 10 workers)

### Security Requirements
- **Authentication**: JWT tokens required for all endpoints
- **Input Validation**: Multi-layer validation (OpenAPI + Domain + Resource)
- **SQL Injection**: Prevented via GORM parameterized queries
- **Error Handling**: No sensitive data in error messages

---

**This technical context provides the foundation for development, testing, and deployment of Quest Manager, ensuring consistent practices across the development team.**

