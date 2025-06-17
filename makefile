

# Makefile
.PHONY: build install clean test

BINARY_NAME=devopsctl
BUILD_DIR=build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

dev-install: build
	@echo "Installing $(BINARY_NAME) for development..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) ~/bin/ 2>/dev/null || cp $(BUILD_DIR)/$(BINARY_NAME) ~/.local/bin/

help:
	@echo "Available commands:"
	@echo "  build       - Build the binary"
	@echo "  install     - Build and install system-wide"
	@echo "  dev-install - Build and install to user bin directory"
	@echo "  clean       - Remove build artifacts"
	@echo "  test        - Run tests"
	@echo "  help        - Show this help message"