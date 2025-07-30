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
# DEV SHORTCUT
# ========================

.PHONY: dev
dev: gen-api run
