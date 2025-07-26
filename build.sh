#!/bin/bash

# Build script for cross-platform compilation
VERSION="v0.1.0"
BINARY_NAME="todo"

# Create build directory
mkdir -p build

# Build for different platforms
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o build/${BINARY_NAME}-windows-amd64.exe

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o build/${BINARY_NAME}-linux-amd64

echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o build/${BINARY_NAME}-darwin-amd64

echo "Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -o build/${BINARY_NAME}-darwin-arm64

echo "Build complete! Binaries are in the 'build' directory:"
ls -la build/