.PHONY: build run test clean migrate-up migrate-down docker-up docker-down

# Build the application
build:
	go build -o bin/example-service cmd/example-service/main.go

# Run the application
run:
	go run cmd/example-service/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5002/community_builds_example_service_db?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost:5002/community_builds_example_service_db?sslmode=disable" down

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development with hot reload
dev:
	air 