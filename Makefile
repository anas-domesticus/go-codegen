BINARY_NAME := go-codegen
SRC := ./...

.PHONY: all test build clean

all: test build

# Run unit tests
test:
	@echo "Running unit tests..."
	go test -v $(SRC)

# Build the CLI
build:
	@echo "Building the CLI..."
	mkdir -p .out
	go build -o .out/$(BINARY_NAME) main.go

# Installs the tool
install: build
	cp .out/$(BINARY_NAME) $(HOME)/.local/bin

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f .out/$(BINARY_NAME)