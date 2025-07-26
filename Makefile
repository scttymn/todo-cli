BINARY_NAME=todo
VERSION=v0.1.0
BUILD_DIR=build

# Default target
all: clean build

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build for all platforms
build: $(BUILD_DIR)
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	@echo "Build complete!"
	@ls -la $(BUILD_DIR)/

# Build for current platform only
local:
	go build -o $(BINARY_NAME)

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)

# Install locally (current platform)
install: local
	sudo mv $(BINARY_NAME) /usr/local/bin/

.PHONY: all build local clean install