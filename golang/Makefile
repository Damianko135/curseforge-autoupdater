.PHONY: build clean test install deps run

# Build the application
build:
	go build -o curseforge-updater.exe ./cmd

# Build for multiple platforms
build-all:
	GOOS=windows GOARCH=amd64 go build -o dist/curseforge-updater-windows-amd64.exe ./cmd
	GOOS=linux GOARCH=amd64 go build -o dist/curseforge-updater-linux-amd64 ./cmd
	GOOS=darwin GOARCH=amd64 go build -o dist/curseforge-updater-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build -o dist/curseforge-updater-darwin-arm64 ./cmd

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Run the application
run:
	go run ./cmd

# Clean build artifacts
clean:
	rm -f curseforge-updater.exe
	rm -rf dist/

# Install the application to GOPATH/bin
install:
	go install ./cmd

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Security check (requires gosec)
security:
	gosec ./...
