#!/bin/bash
set -e

echo "🚀 Compiling speedtest_cli multi-platform static binaries..."

mkdir -p binaries

export CGO_ENABLED=0

# Linux x86_64
echo "📦 Building Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o binaries/speedtest-linux-amd64 ./cmd/speedtest

# Linux ARM64
echo "📦 Building Linux arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o binaries/speedtest-linux-arm64 ./cmd/speedtest

# macOS Intel
echo "📦 Building Darwin amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o binaries/speedtest-darwin-amd64 ./cmd/speedtest

# macOS Apple Silicon
echo "📦 Building Darwin arm64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o binaries/speedtest-darwin-arm64 ./cmd/speedtest

# Windows x86_64
echo "📦 Building Windows amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o binaries/speedtest-windows-amd64.exe ./cmd/speedtest

echo "✅ All binaries built successfully in binaries/!"
ls -lh binaries/
