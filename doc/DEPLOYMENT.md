# Deployment Guide - Quest Manager

## 🚀 Deployment Options

Quest Manager supports multiple deployment strategies:
1. Docker Compose (recommended for development/staging)
2. Kubernetes (recommended for production)
3. Standalone binary (simple deployments)

---

## 🐳 Docker Deployment

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+

### Using Docker Compose

**1. Configure environment:**
```bash
# Create .env file
cat > .env << EOF
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=quest_manager
DB_SSLMODE=disable

AUTH_GRPC=quest-auth:50051
MIDDLEWARE_ENABLE_AUTH=true

SERVER_PORT=8080
EVENT_GOROUTINE_LIMIT=10
EOF
```

**2. Start services:**
```bash
docker-compose up -d
```

**3. Check logs:**
```bash
docker-compose logs -f quest-manager
```

**4. Stop services:**
```bash
docker-compose down
```

### Docker Compose Configuration

See `docker-compose.yml` for full configuration:
- PostgreSQL database with persistent volume
- Quest Manager API
- Network configuration
- Health checks

---

## 🏗️ Building Docker Image

### Multi-Stage Build

**Build image:**
```bash
docker build -t quest-manager:1.4.0 .
docker build -t quest-manager:latest .
```

**Dockerfile stages:**
1. **Builder:** Compile Go binary
2. **Runtime:** Minimal Alpine image with binary

**Optimizations:**
- Layer caching for dependencies
- Multi-stage build for small image size
- Non-root user for security

### Push to Registry

```bash
# Tag for registry
docker tag quest-manager:1.4.0 registry.example.com/quest-manager:1.4.0

# Push
docker push registry.example.com/quest-manager:1.4.0
```

---

## ☸️ Kubernetes Deployment

### Basic Deployment

**1. Create namespace:**
```bash
kubectl create namespace quest-manager
```

**2. Create secrets:**
```bash
kubectl create secret generic quest-db-credentials \
  --from-literal=username=postgres \
  --from-literal=password=your-secure-password \
  -n quest-manager
```

**3. Deploy PostgreSQL:**
```yaml
# postgres-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: quest-manager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:14-alpine
        env:
        - name: POSTGRES_DB
          value: quest_manager
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: quest-db-credentials
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: quest-db-credentials
              key: password
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: quest-manager
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
```

**4. Deploy Quest Manager:**
```yaml
# quest-manager-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: quest-manager
  namespace: quest-manager
spec:
  replicas: 3
  selector:
    matchLabels:
      app: quest-manager
  template:
    metadata:
      labels:
        app: quest-manager
    spec:
      containers:
      - name: quest-manager
        image: registry.example.com/quest-manager:1.4.0
        env:
        - name: DB_HOST
          value: postgres
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: quest-db-credentials
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: quest-db-credentials
              key: password
        - name: DB_NAME
          value: quest_manager
        - name: DB_SSLMODE
          value: disable
        - name: AUTH_GRPC
          value: quest-auth:50051
        - name: MIDDLEWARE_ENABLE_AUTH
          value: "true"
        - name: SERVER_PORT
          value: "8080"
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: quest-manager
  namespace: quest-manager
spec:
  selector:
    app: quest-manager
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

**5. Apply configurations:**
```bash
kubectl apply -f postgres-deployment.yaml
kubectl apply -f quest-manager-deployment.yaml
```

---

## 🔒 Production Deployment

### Security Checklist

- [ ] Enable SSL for database (`DB_SSLMODE=require`)
- [ ] Use strong database password
- [ ] Enable authentication (`MIDDLEWARE_ENABLE_AUTH=true`)
- [ ] Configure Auth service (`AUTH_GRPC`)
- [ ] Use secrets management (not env files)
- [ ] Enable HTTPS for API
- [ ] Set up firewall rules
- [ ] Regular security updates
- [ ] Monitor for vulnerabilities

### Database Security

**1. Create dedicated user:**
```sql
CREATE USER quest_app WITH PASSWORD 'strong-password';
GRANT CONNECT ON DATABASE quest_manager TO quest_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO quest_app;
```

**2. Enable SSL:**
```bash
DB_SSLMODE=require
# Or even stricter
DB_SSLMODE=verify-full
```

**3. Connection pooling:**
```go
// Configure in code
db.DB().SetMaxOpenConns(25)
db.DB().SetMaxIdleConns(5)
db.DB().SetConnMaxLifetime(5 * time.Minute)
```

---

## 📊 Monitoring & Observability

### Health Checks

**Liveness probe:**
```bash
curl http://localhost:8080/health
```

**Readiness probe:**
Check database connectivity and auth service availability.

### Metrics (Future)

Recommended metrics to track:
- Request rate (requests/sec)
- Response time (p50, p95, p99)
- Error rate (4xx, 5xx)
- Database connection pool usage
- Quest creation/assignment rates

### Logging

**Structured logging:**
```go
log.Printf("INFO publishing domain event event_name=%s event_id=%s", 
    event.Name(), event.ID())
```

**Log levels:**
- `ERROR` - Errors requiring attention
- `WARN` - Warnings (deprecated, potential issues)
- `INFO` - Important events (quest created, assigned)
- `DEBUG` - Detailed debugging info

### Tracing (Future)

Consider adding:
- OpenTelemetry integration
- Request ID tracking
- Distributed tracing

---

## 🔄 Database Migrations

### Current: Auto-Migration
```go
// cmd/database.go
db.AutoMigrate(
    &questrepo.QuestModel{},
    &locationrepo.LocationModel{},
    &eventrepo.EventModel{},
)
```

### Future: Versioned Migrations

**Recommended tools:**
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)

**Migration workflow:**
```bash
# Create migration
migrate create -ext sql -dir migrations -seq add_quest_priority

# Apply migrations
migrate -path migrations -database "postgres://..." up

# Rollback
migrate -path migrations -database "postgres://..." down 1
```

---

## 🌐 Reverse Proxy Setup

### Nginx Configuration

```nginx
upstream quest_manager {
    server quest-manager-1:8080;
    server quest-manager-2:8080;
    server quest-manager-3:8080;
}

server {
    listen 80;
    server_name api.example.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    location /api/v1/ {
        proxy_pass http://quest_manager;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 10s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
    }
    
    location /health {
        proxy_pass http://quest_manager;
        access_log off;
    }
}
```

---

## 🔐 Secrets Management

### Development
```bash
# .env file (gitignored)
DB_PASSWORD=dev-password
```

### Kubernetes
```bash
# Use Kubernetes secrets
kubectl create secret generic quest-secrets \
  --from-literal=db-password=prod-password \
  --from-literal=auth-token=secret-token
```

### Cloud Providers

**AWS Secrets Manager:**
```go
// Load secrets on startup
secret := getSecretFromAWS("quest-manager/db-password")
cfg.DbPassword = secret
```

**Google Secret Manager:**
```go
secret := getSecretFromGCP("projects/*/secrets/db-password")
cfg.DbPassword = secret
```

---

## 📈 Scaling

### Horizontal Scaling

**Quest Manager is stateless** and can be scaled horizontally:

```bash
# Docker Compose
docker-compose up -d --scale quest-manager=3

# Kubernetes
kubectl scale deployment quest-manager --replicas=5
```

**Considerations:**
- No shared state between instances
- Each request gets new UnitOfWork
- Database connection pool per instance
- Load balancer distributes traffic

### Database Scaling

**Read Replicas:**
- Use read replicas for queries
- Master for commands (writes)
- Configure in repository layer

**Connection Pooling:**
```go
db.DB().SetMaxOpenConns(25)  // Max connections per instance
db.DB().SetMaxIdleConns(5)   // Idle connections
```

---

## 🧪 Deployment Testing

### Pre-Deployment Tests
```bash
# Run full test suite
go test ./... -v -p 1

# Integration tests
go test -tags=integration ./tests/integration/... -v -p 1

# Build verification
go build ./...

# Linter
golangci-lint run
```

### Smoke Tests (Post-Deployment)
```bash
# Health check
curl https://api.example.com/health

# Create quest
curl -X POST https://api.example.com/api/v1/quests \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", ...}'

# Get quests
curl https://api.example.com/api/v1/quests \
  -H "Authorization: Bearer ${TOKEN}"
```

---

## 🔄 Deployment Checklist

### Before Deployment
- [ ] All tests passing
- [ ] Code reviewed
- [ ] CHANGELOG updated
- [ ] Version bumped
- [ ] Database migrations prepared
- [ ] Configuration verified
- [ ] Secrets configured
- [ ] Backup database

### During Deployment
- [ ] Build Docker image
- [ ] Tag with version
- [ ] Push to registry
- [ ] Update K8s deployment
- [ ] Monitor rollout
- [ ] Verify health checks

### After Deployment
- [ ] Run smoke tests
- [ ] Check logs for errors
- [ ] Monitor metrics
- [ ] Verify database connectivity
- [ ] Test critical paths
- [ ] Document deployment

---

## 🆘 Rollback Procedure

### Docker Compose
```bash
# Stop current version
docker-compose down

# Checkout previous version
git checkout v1.3.0

# Rebuild and start
docker-compose up -d --build
```

### Kubernetes
```bash
# Rollback to previous revision
kubectl rollout undo deployment/quest-manager -n quest-manager

# Or specific revision
kubectl rollout undo deployment/quest-manager --to-revision=2 -n quest-manager

# Check status
kubectl rollout status deployment/quest-manager -n quest-manager
```

---

## 📊 Performance Tuning

### Database Optimization

**Indexes:**
```sql
CREATE INDEX idx_quests_status ON quests(status);
CREATE INDEX idx_quests_assignee ON quests(assignee) WHERE assignee IS NOT NULL;
CREATE INDEX idx_locations_coords ON locations(latitude, longitude);
```

**Connection Pool:**
```bash
# Increase for high traffic
export DB_MAX_OPEN_CONNS=50
export DB_MAX_IDLE_CONNS=10
```

### Application Tuning

**Event Publishing:**
```bash
# Increase goroutine limit for high event volume
export EVENT_GOROUTINE_LIMIT=50
```

**Go Runtime:**
```bash
# Increase max processors
export GOMAXPROCS=4
```

---

## 🔗 Related

- [Configuration](CONFIGURATION.md) - Environment configuration
- [Development](DEVELOPMENT.md) - Local development setup
- [Architecture](ARCHITECTURE.md) - System design

---

**Deployment Version:** 1.4.0  
**Last Updated:** October 9, 2025

