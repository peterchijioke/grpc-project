.PHONY: help build-server build-client run-server run-client generate test clean docker-build docker-run docker-compose-up docker-compose-down

help:
	@echo "Available targets:"
	@echo "  generate         - Regenerate gRPC code from proto files"
	@echo "  build-server     - Build the server binary"
	@echo "  build-client     - Build the client binary"
	@echo "  run-server       - Run the server"
	@echo "  run-client       - Run the client (use CMD=signup|login|greet|stream)"
	@echo "  test             - Run tests"
	@echo "  clean            - Remove build artifacts"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-run       - Run Docker container"
	@echo "  docker-compose-up - Start server + PostgreSQL with docker-compose"
	@echo "  docker-compose-down - Stop all services"

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto

build-server:
	go build -o bin/server ./cmd/server

build-client:
	go build -o bin/client ./cmd/client

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client -cmd=$(CMD) -email=$(EMAIL) -password=$(PASSWORD) -name=$(NAME)

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf *.test *.out coverage.*

docker-build:
	docker build --no-cache -t grpc-server .

docker-run:
	docker run -p 50051:50051 grpc-server

docker-compose-up:
	docker-compose build --no-cache && docker-compose up -d

docker-compose-down:
	docker-compose down

