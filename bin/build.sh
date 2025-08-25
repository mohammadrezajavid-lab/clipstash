#!/bin/bash

VERSION="v1.0.0"

# This script builds the application for all major platforms,
# ensuring CGO is disabled for true cross-compilation.

echo "Building version $VERSION..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="${SCRIPT_DIR}/.."
OUT_DIR="${ROOT_DIR}/bin/bin"
mkdir -p "${OUT_DIR}"

echo "Building for Linux ..."
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/clipstash-linux-amd64-${VERSION}"  "${ROOT_DIR}/main.go"
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o "${OUT_DIR}/clipstash-linux-arm64-${VERSION}"  "${ROOT_DIR}/main.go"

echo "Building for Windows ..."
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "${OUT_DIR}/clipstash-windows-amd64-${VERSION}.exe" "${ROOT_DIR}/main.go"
env CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o "${OUT_DIR}/clipstash-windows-arm64-${VERSION}.exe" "${ROOT_DIR}/main.go"

echo "Building for macOS ..."
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o "${OUT_DIR}/clipstash-darwin-amd64-${VERSION}"  "${ROOT_DIR}/main.go"
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o "${OUT_DIR}/clipstash-darwin-arm64-${VERSION}"  "${ROOT_DIR}/main.go"

echo "Build complete! Your binaries for version $VERSION are in ${OUT_DIR}."