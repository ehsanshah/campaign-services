# ==============================================================================
# 🛠️ Variables
# ==============================================================================

APP_NAME = campaign-service
BINARY_NAME = campaign-server
GO_MAIN_FILE = src/cmd/server/main.go
DOCKER_FILE_PATH = deployments/docker/Dockerfile
DOCKER_IMAGE_TAG = latest

# مسیرهای پروتو
PROTO_SRC_DIR = src/api/proto
PROTO_OUT_DIR = src/pkg/pb
# پیدا کردن تمام فایل‌های پروتو
PROTO_FILES := $(shell find $(PROTO_SRC_DIR) -name "*.proto")

# تنظیمات داکر کامپوز (زیرساخت)
DOCKER_COMPOSE_FILE = deployments/docker/docker-compose.yml
DB_CONTAINER_NAME = campaign_postgres
DB_USER = user
DB_NAME = campaign_db
MIGRATION_FILE = migrations/01_init.sql

# ==============================================================================
# 📋 Commands
# ==============================================================================

.PHONY: help all proto update-submodule run build clean run-infra stop-infra db-init docker-build

help:
	@echo "🛠️  Available Commands:"
	@echo "  make update-submodule  -> 🔄 Update git submodules (protos)"
	@echo "  make proto             -> 📄 Generate gRPC code from .proto files"
	@echo "  make run-infra         -> 🐘 Start Postgres & Adminer"
	@echo "  make db-init           -> 💽 Apply SQL Migrations to Postgres"
	@echo "  make run               -> 🚀 Run the Go Application locally"
	@echo "  make build             -> 🔨 Compile Go binary"
	@echo "  make docker-build      -> 🐳 Build Docker Image for this Service"
	@echo "  make stop-infra        -> 🛑 Stop all containers"
	@echo "  make clean             -> 🧹 Remove binaries and generated files"

all: proto build

# ==============================================================================
# 🔗 Git & Proto
# ==============================================================================

update-submodule:
	@echo "🔄 Updating git submodules..."
	git submodule update --init --recursive --remote

proto: update-submodule
	@echo "🗑️  Cleaning old generated files..."
	rm -rf $(PROTO_OUT_DIR)
	mkdir -p $(PROTO_OUT_DIR)
	@echo "🚀 Generating gRPC code..."
	protoc \
		--proto_path=$(PROTO_SRC_DIR) \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "✅ Proto generation complete!"

# ==============================================================================
# 🐘 Infrastructure (Docker & DB)
# ==============================================================================

run-infra:
	@echo "🐳 Starting Database & Tools..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "⏳ Waiting for Database to be ready..."
	@sleep 5

stop-infra:
	@echo "🛑 Stopping Infrastructure..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

db-init:
	@echo "💽 Applying Migrations to $(DB_NAME)..."
	cat $(MIGRATION_FILE) | docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)
	@echo "✅ Database initialized successfully!"

# ==============================================================================
# 🚀 Application Development
# ==============================================================================

tidy:
	@echo "📦 Tidy up go modules..."
	go mod tidy

run: tidy
	@echo "🚀 Starting $(APP_NAME)..."
	go run $(GO_MAIN_FILE)

build: tidy
	@echo "🔨 Building binary..."
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) $(GO_MAIN_FILE)
	@echo "✅ Build complete: bin/$(BINARY_NAME)"

# ==============================================================================
# 🐳 Docker Build (درخواست جدید شما)
# ==============================================================================

docker-build:
	@echo "🐳 Building Docker Image: $(APP_NAME):$(DOCKER_IMAGE_TAG)..."
	docker build -f $(DOCKER_FILE_PATH) -t $(APP_NAME):$(DOCKER_IMAGE_TAG) .
	@echo "✅ Docker Image built successfully!"

clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin
	rm -rf $(PROTO_OUT_DIR)