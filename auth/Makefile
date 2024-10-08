# Simple Makefile for a Go project

# Build the application
# Example: make build
all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
# Example: make run
run:
	@go run cmd/api/main.go

# Create DB container
# Example: make docker-run
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
# Example: make docker-down
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
# Example: make test
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
# Example: make itest
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
# Example: make clean
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
# Example: make watch
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Create a new migration
# Example: make make-migration
make-migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate:
	@migrate -path ./migrations -database postgres://melkey:password1234@localhost:5432/blueprint?sslmode=disable up

.PHONY: all build run test clean watch make-migration migrate
