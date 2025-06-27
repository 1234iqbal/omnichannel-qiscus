.PHONY: build run test clean docker-build docker-up docker-down

# Build the application
build:
	go build -o bin/$(APP_NAME) cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

