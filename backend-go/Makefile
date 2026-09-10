.PHONY: run worker test build docker-up docker-down

run:
	go run ./cmd/api

worker:
	go run ./cmd/worker

test:
	cd ../tests-backend-go && go test -v ./...

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

smoke-test:
	@echo "🔥 Running Smoke Tests against live Docker ecosystem..."
	@curl -s -f http://localhost:8080/api/v1/guest/system/health || (echo "❌ Health check failed" && exit 1)
	@echo "✅ API is Healthy"
	@curl -s -f http://localhost:3000 || (echo "❌ Admin Portal failed" && exit 1)
	@echo "✅ Admin Portal is Up"
	@curl -s -f http://localhost:3001 || (echo "❌ Client Portal failed" && exit 1)
	@echo "✅ Client Portal is Up"
	@echo "🎯 Smoke tests passed successfully!"
