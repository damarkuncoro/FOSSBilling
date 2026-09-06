.PHONY: all build test dev docker-up docker-down clean

all: test build

# Run all 25 Go unit & integration test suites
test:
	@echo "🧪 Running Go test suites in tests-backend-go/..."
	cd tests-backend-go && go test -v ./...

# Build all Go binaries and frontend bundles
build:
	@echo "🔨 Building Go backend binaries..."
	cd backend-go && CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/api ./cmd/api
	cd backend-go && CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/worker ./cmd/worker
	cd backend-go && CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/cli ./cmd/cli
	@echo "🔨 Building Administrator Portal..."
	cd frontend-administrator && npm run build
	@echo "🔨 Building Customer Portal..."
	cd frontend-client && npm run build
	@echo "✅ All artifacts built successfully!"

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
	@echo "🧹 Clean complete!"
