.PHONY: run test build docker-up docker-down

run:
	go run ./cmd/api

worker:
	go run ./cmd/worker

test:
	go test -v ./...

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

docker-up:
	docker compose -f ../deploy/docker-compose.yml up --build -d

docker-down:
	docker compose -f ../deploy/docker-compose.yml down
