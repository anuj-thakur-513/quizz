# Define variables
APP_NAME = quizz

# Default target
all: build

# Build the application
build: lint, install
	go build -o $(APP_NAME)

# Run the application
dev:
	air

# Start the server (example)
start: lint
	go run cmd/server/main.go

# Install dependencies (if any)
install:
	go mod tidy

# Lint (using a tool like golangci-lint)
lint:
	golangci-lint run
