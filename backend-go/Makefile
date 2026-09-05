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
