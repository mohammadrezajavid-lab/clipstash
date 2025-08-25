#!/bin/bash

VERSION="v1.0.0"

# This script builds the application for all major platforms,
# ensuring CGO is disabled for true cross-compilation.

echo "Building version $VERSION..."

echo "Building for Linux..."
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/clipstash-linux-amd64-$VERSION ../clipstash.go
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/clipstash-linux-arm64-$VERSION ../clipstash.go

echo "Building for Windows..."
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/clipstash-windows-amd64-$VERSION.exe ../clipstash.go
env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./bin/clipstash-windows-arm64-$VERSION.exe ../clipstash.go

echo "Building for macOS..."
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/clipstash-mac-amd64-$VERSION ../clipstash.go
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/clipstash-mac-arm64-$VERSION ../clipstash.go

echo "Build complete! Your binaries for version $VERSION are in the ./bin directory."