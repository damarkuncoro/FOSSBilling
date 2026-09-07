.PHONY: all build test dev demo docker-up docker-down clean help

all: test build

# Help command to list all targets
help:
	@echo "Available targets:"
	@echo "  make build       - Build all binaries (api, worker, cli) and frontend bundles"
	@echo "  make test        - Run all Go backend test suites (44 packages)"
	@echo "  make test-front  - Run all React frontend test suites"
	@echo "  make test-all    - Run complete backend and frontend test suites"
	@echo "  make dev         - Run backend-go in development mode (requires air)"
	@echo "  make demo        - Execute live business simulation"
	@echo "  make cli-status  - Run CLI system and provider driver audit"
	@echo "  make docs-dev    - Run VitePress documentation portal in dev mode"
	@echo "  make docs-build  - Build production VitePress documentation site"
	@echo "  make docker-up   - Start full stack with Docker Compose"
	@echo "  make docker-down - Stop Docker Compose stack"
	@echo "  make clean       - Remove build artifacts"

# Run all Go unit & integration test suites
test:
	@echo "🧪 Running Go test suites in tests-backend-go/..."
	cd tests-backend-go && go test -v ./...

# Run frontend test suites
test-front:
	@echo "🧪 Running Frontend Administrator tests..."
	cd frontend-administrator && npm test -- --run
	@echo "🧪 Running Frontend Client tests..."
	cd frontend-client && npm test -- --run

# Run full ecosystem tests
test-all: test test-front
	@echo "✅ All backend and frontend test suites passed!"

# Build all Go binaries and frontend bundles
build:
	@echo "🔨 Building Go backend binaries..."
	mkdir -p backend-go/bin
	cd backend-go && CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/api ./cmd/api
	cd backend-go && CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/worker ./cmd/worker
	cd backend-go && CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/cli ./cmd/cli
	@echo "🔨 Building Administrator Portal..."
	cd frontend-administrator && npm install && npm run build
	@echo "🔨 Building Customer Portal..."
	cd frontend-client && npm install && npm run build
	@echo "✅ All artifacts built successfully!"

# Run backend in development mode with live reload (requires 'air')
dev:
	@echo "🚀 Starting backend-go in development mode..."
	cd backend-go && air -c .air.toml

# Run end-to-end live business simulation in Go
demo:
	@echo "🚀 Running E2E live simulation..."
	cd backend-go && go run ./cmd/demo

# Start full multi-container stack with Docker Compose
docker-up:
	@echo "🐳 Starting full FOSSBilling stack with Docker Compose..."
	docker compose -f deploy/docker-compose.yml up -d --build

# Stop Docker Compose stack
docker-down:
	@echo "🛑 Stopping Docker Compose stack..."
	docker compose -f deploy/docker-compose.yml down

# Clean temporary build artifacts
clean:
	rm -rf backend-go/bin/
	rm -rf frontend-administrator/dist/
	rm -rf frontend-client/dist/
	rm -rf docs/.vitepress/dist/
	@echo "🧹 Clean complete!"

# Documentation portal dev and build commands
docs-dev:
	@echo "📖 Starting VitePress documentation portal..."
	cd docs && npm run docs:dev

docs-build:
	@echo "🔨 Building VitePress documentation site..."
	cd docs && npm run docs:build

