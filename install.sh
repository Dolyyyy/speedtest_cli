#!/bin/bash
# speedtest_cli - Multi-architecture installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.sh | sudo bash

set -e

echo "🚀 Installing speedtest_cli..."
echo "🔗 Repo: https://github.com/Dolyyyy/speedtest_cli"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux)
        case "$ARCH" in
            x86_64)
                BINARY="speedtest-linux-amd64"
                ;;
            aarch64|arm64)
                BINARY="speedtest-linux-arm64"
                ;;
            *)
                echo "❌ Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    darwin)
        case "$ARCH" in
            x86_64)
                BINARY="speedtest-darwin-amd64"
                ;;
            arm64)
                BINARY="speedtest-darwin-arm64"
                ;;
            *)
                echo "❌ Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac
        ;;
    *)
        echo "❌ Unsupported operating system: $OS"
        exit 1
        ;;
esac

TARGET_DIR="/usr/local/bin"
if [ ! -w "$TARGET_DIR" ] && [ "$EUID" -ne 0 ]; then
    TARGET_DIR="$HOME/.local/bin"
    mkdir -p "$TARGET_DIR"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" 2>/dev/null && pwd || true)"

rm -f "$TARGET_DIR/speedtest" 2>/dev/null || true

if [ -n "$SCRIPT_DIR" ] && [ -f "$SCRIPT_DIR/binaries/$BINARY" ]; then
    echo "📦 Installing from local binary binaries/$BINARY..."
    cp "$SCRIPT_DIR/binaries/$BINARY" "$TARGET_DIR/speedtest"
else
    echo "📥 Downloading $BINARY from GitHub..."
    curl -fsSL "https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/binaries/$BINARY?t=$(date +%s)" -o "$TARGET_DIR/speedtest"
fi

chmod +x "$TARGET_DIR/speedtest"

echo "✅ speedtest_cli installed successfully to $TARGET_DIR/speedtest!"
echo ""
echo "Quick start:"
echo "  speedtest          # Run full speedtest"
echo "  speedtest --bytes  # Speed in MB/s"
echo "  speedtest --help   # View all options"
echo ""
echo "🔗 Star us on GitHub: https://github.com/Dolyyyy/speedtest_cli"
