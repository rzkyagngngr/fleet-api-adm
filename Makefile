.PHONY: run build tidy dev

# Run the application
run:
	go run cmd/main.go

# Build binary
build:
	go build -o bin/app cmd/main.go

# Tidy modules
tidy:
	go mod tidy

# Run with hot reload (requires: go install github.com/air-verse/air@latest)
dev:
	air

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
