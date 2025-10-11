# Configuration Guide - Quest Manager

## 🔧 Environment Variables

### Database Configuration

| Variable      | Description       | Default         | Required |
|---------------|-------------------|-----------------|----------|
| `DB_HOST`     | PostgreSQL host   | `localhost`     | ✅        |
| `DB_PORT`     | PostgreSQL port   | `5432`          | ✅        |
| `DB_USER`     | Database user     | `postgres`      | ✅        |
| `DB_PASSWORD` | Database password | -               | ✅        |
| `DB_NAME`     | Database name     | `quest_manager` | ✅        |
| `DB_SSLMODE`  | SSL mode          | `disable`       | ❌        |

### Server Configuration

| Variable      | Description      | Default   | Required |
|---------------|------------------|-----------|----------|
| `SERVER_PORT` | HTTP server port | `8080`    | ❌        |
| `SERVER_HOST` | HTTP server host | `0.0.0.0` | ❌        |

### Authentication Configuration

| Variable    | Description               | Default | Required |
|-------------|---------------------------|---------|----------|
| `AUTH_GRPC` | Auth service gRPC address | -       | ⚠️       |

**⚠️ Note:** If `AUTH_GRPC` is not set, authentication is disabled (dev mode only).

### Middleware Configuration

| Variable                  | Description          | Default        | Required |
|---------------------------|----------------------|----------------|----------|
| `MIDDLEWARE_ENABLE_AUTH`  | Enable JWT auth      | `true`         | ❌        |
| `DEV_AUTH_HEADER_NAME`    | Dev mode auth header | `X-User-ID`    | ❌        |
| `DEV_AUTH_STATIC_USER_ID` | Dev mode static user | `00000000-...` | ❌        |

### Event Publishing

| Variable                | Description          | Default | Required |
|-------------------------|----------------------|---------|----------|
| `EVENT_GOROUTINE_LIMIT` | Max event goroutines | `10`    | ❌        |

---

## 📁 Configuration Files

### `config.example`
Example configuration file (copy to `.env`).

### `configs/server.cfg.yaml`
Server-specific configuration (experimental).

---

## 🚀 Configuration Loading

Configuration is loaded from environment variables using `cmd/config.go`:

```go
type Config struct {
    // Database
    DbHost     string
    DbPort     string
    DbUser     string
    DbPassword string
    DbName     string
    DbSslMode  string
    
    // Server
    ServerPort string
    ServerHost string
    
    // Auth
    AuthGRPC string
    
    // Middleware
    Middleware MiddlewareConfig
    
    // Events
    EventGoroutineLimit int
}
```

**Loading Priority:**
1. Environment variables
2. Default values in code
3. Fallback to safe defaults

---

## 🐳 Docker Configuration

### `docker-compose.yml`
Orchestrates services:
- PostgreSQL database
- Quest Manager API
- (Future: Quest Auth service)

**Usage:**
```bash
docker-compose up -d
```

### `Dockerfile`
Multi-stage build for Quest Manager:
1. Build stage: compile Go binary
2. Runtime stage: minimal Alpine image

**Build:**
```bash
docker build -t quest-manager:latest .
```

---

## 🔒 Security Configuration

### Production Settings
```bash
# Database with SSL
DB_SSLMODE=require

# Enable authentication
MIDDLEWARE_ENABLE_AUTH=true
AUTH_GRPC=auth-service:50051
```

### Development Settings
```bash
# Database without SSL
DB_SSLMODE=disable

# Optional: disable auth for local development
MIDDLEWARE_ENABLE_AUTH=false
# Or use dev mode with mock auth
DEV_AUTH_STATIC_USER_ID=00000000-0000-0000-0000-000000000001
```

### Testing Settings
```bash
# Separate test database
DB_NAME=quest_test
DB_SSLMODE=disable

# Auth disabled for most tests (using mocks)
MIDDLEWARE_ENABLE_AUTH=false
```

---

## 🗄️ Database Setup

### Create Database
```sql
CREATE DATABASE quest_manager;
```

### Migrations
Migrations are handled automatically on startup via GORM AutoMigrate:
- `quests` table
- `locations` table
- `events` table

**Tables:**

**quests:**
```sql
CREATE TABLE quests (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    reward INTEGER NOT NULL,
    duration_minutes INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,
    creator VARCHAR(255) NOT NULL,
    assignee UUID,
    target_location_id UUID,
    execution_location_id UUID,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (target_location_id) REFERENCES locations(id),
    FOREIGN KEY (execution_location_id) REFERENCES locations(id)
);
```

**locations:**
```sql
CREATE TABLE locations (
    id UUID PRIMARY KEY,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    address VARCHAR(500),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

**events:**
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    event_name VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    event_data JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

---

## 🎛️ Configuration by Environment

### Local Development
```bash
# .env.local
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=quest_manager_dev
DB_SSLMODE=disable

SERVER_PORT=8080

# Auth disabled or using dev mode
MIDDLEWARE_ENABLE_AUTH=false
DEV_AUTH_STATIC_USER_ID=00000000-0000-0000-0000-000000000001
```

### Testing (CI/CD)
```bash
# .env.test
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=quest_test
DB_SSLMODE=disable

MIDDLEWARE_ENABLE_AUTH=false
```

### Staging
```bash
# .env.staging
DB_HOST=staging-db.example.com
DB_PORT=5432
DB_USER=quest_app
DB_PASSWORD=${DB_PASSWORD_SECRET}
DB_NAME=quest_manager_staging
DB_SSLMODE=require

SERVER_PORT=8080

MIDDLEWARE_ENABLE_AUTH=true
AUTH_GRPC=auth-staging.example.com:50051

EVENT_GOROUTINE_LIMIT=10
```

### Production
```bash
# .env.production
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_USER=quest_app
DB_PASSWORD=${DB_PASSWORD_SECRET}
DB_NAME=quest_manager
DB_SSLMODE=require

SERVER_PORT=8080

MIDDLEWARE_ENABLE_AUTH=true
AUTH_GRPC=auth.example.com:50051

EVENT_GOROUTINE_LIMIT=20
```

---

## 🔍 Configuration Validation

### Required Fields Check
Application validates required configuration on startup:
```go
func validateConfig(cfg Config) error {
    if cfg.DbHost == "" {
        return errors.New("DB_HOST is required")
    }
    // ... more validations
    return nil
}
```

### Default Values
```go
const (
    DefaultServerPort = "8080"
    DefaultEventGoroutineLimit = 10
    DefaultDevAuthStaticUserID = "00000000-0000-0000-0000-000000000001"
)
```

---

## 🐞 Debug Configuration

### Enable SQL Logging
```bash
DB_LOG_LEVEL=debug  # Log all SQL queries
```

### Verbose Logging
```bash
LOG_LEVEL=debug     # Detailed application logs
```

---

## 📝 Configuration Best Practices

### 1. Never Commit Secrets
- Use `.env` file (gitignored)
- Use secret management in production
- Rotate credentials regularly

### 2. Environment-Specific Configs
- Separate `.env.local`, `.env.staging`, `.env.production`
- Use different databases per environment
- Enable SSL in production

### 3. Fail Fast on Misconfiguration
- Validate on startup
- Log configuration (without secrets)
- Provide clear error messages

### 4. Use Defaults Wisely
- Sensible defaults for development
- Require explicit config for production
- Document all defaults

---

## 🔗 Related

- [Deployment](DEPLOYMENT.md) - How to deploy with different configs
- [Development](DEVELOPMENT.md) - Local development setup
- [Architecture](ARCHITECTURE.md) - System design

---

## ⚙️ Quick Start Configs

### Minimal Local Setup
```bash
export DB_PASSWORD=postgres
export MIDDLEWARE_ENABLE_AUTH=false
go run cmd/app/main.go
```

### With Docker
```bash
docker-compose up -d
# Uses docker-compose.yml configuration
```

### With Authentication
```bash
export AUTH_GRPC=localhost:50051
export MIDDLEWARE_ENABLE_AUTH=true
go run cmd/app/main.go
```

---

**Remember:** Configuration drives behavior. Always validate and document configuration changes!

