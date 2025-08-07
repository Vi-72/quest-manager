# Путь до OpenAPI схемы
OPENAPI_FILE=api/openapi/openapi.yml
OPENAPI_CONFIG=configs/server.cfg.yaml

# Целевые папки
GEN_DIR=internal/generated/servers

# Бинарное имя (если хочешь собирать проект)
BINARY_NAME=task-server

# ========================
# API
# ========================

.PHONY: gen-api
gen-api:
	oapi-codegen -config $(OPENAPI_CONFIG) $(OPENAPI_FILE)
# ========================
# BUILD
# ========================

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd/app

.PHONY: run
run:
	go run ./cmd/app/main.go

# ========================
# CLEAN
# ========================

.PHONY: clean
clean:
	rm -rf $(BINARY_NAME) $(GEN_DIR)

# ========================
# TESTS
# ========================

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	@echo "🧪 Running unit tests..."
	go test ./tests/domain -v



.PHONY: test-repository
test-repository:
	@echo "🗄️ Running repository integration tests (PostgreSQL)..."
	go test -tags=integration ./tests/integration/cases/repository -v



.PHONY: test-integration
test-integration:
	@echo "🔗 Running integration tests..."
	go test -tags=integration ./tests/integration/... -v



.PHONY: test-fast
test-fast:
	@echo "⚡ Running fast tests only..."
	go test -short ./tests/domain -v

.PHONY: test-coverage
test-coverage:
	@echo "📊 Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-coverage-integration
test-coverage-integration:
	@echo "📊 Generating integration test coverage..."
	go test -tags=integration -coverprofile=coverage-integration.out ./tests/integration/...
	go tool cover -html=coverage-integration.out -o coverage-integration.html

.PHONY: test-bench
test-bench:
	@echo "🚀 Running benchmark tests..."
	go test -bench=. -benchmem ./...

.PHONY: test-all
test-all: test-unit test-repository test-integration
	@echo "✅ All tests completed!"

.PHONY: test-watch
test-watch:
	@echo "👀 Watching for changes and running tests..."
	# Requires 'entr' tool: brew install entr
	find . -name "*.go" | entr -c make test-fast

# ========================
# SCRIPTS
# ========================

.PHONY: test-stats
test-stats:
	@echo "📈 Running test statistics script..."
	@chmod +x scripts/test-stats.sh
	@./scripts/test-stats.sh

.PHONY: test-stats-new
test-stats-new:
	@echo "📊 Running new test statistics script..."
	@chmod +x scripts/test-stats-new.sh
	@./scripts/test-stats-new.sh

.PHONY: coverage-check
coverage-check:
	@echo "🔍 Running coverage check script..."
	@chmod +x scripts/coverage-check.sh
	@./scripts/coverage-check.sh

.PHONY: coverage-report
coverage-report:
	@echo "📋 Running coverage report script..."
	@chmod +x scripts/coverage-report.sh
	@./scripts/coverage-report.sh

# ========================
# DEV SHORTCUT
# ========================

.PHONY: dev
dev: gen-api run
