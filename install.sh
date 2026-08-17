#!/usr/bin/env bash
set -e

REPO="Dolyyyy/speedtest_cli"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="${1:-${BINARY_NAME:-speedtest}}"

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    linux)
        BINARY_FILE="speedtest-linux-${ARCH}"
        ;;
    darwin)
        BINARY_FILE="speedtest-darwin-${ARCH}"
        ;;
    *)
        echo "❌ Unsupported operating system: $OS"
        exit 1
        ;;
esac

# Fallback to local user bin directory if /usr/local/bin is not writable
if [ "$EUID" -ne 0 ] && [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

TARGET_PATH="${INSTALL_DIR}/${BINARY_NAME}"

echo "🚀 Installing speedtest_cli..."
echo "🔗 Repo: https://github.com/${REPO}"

# Check for local binary first (during development / testing)
if [ -f "binaries/${BINARY_FILE}" ]; then
    echo "📦 Installing from local binary binaries/${BINARY_FILE}..."
    cp "binaries/${BINARY_FILE}" "$TARGET_PATH"
else
    DOWNLOAD_URL="https://github.com/${REPO}/raw/main/binaries/${BINARY_FILE}?t=$(date +%s)"
    echo "📥 Downloading ${BINARY_FILE} from GitHub..."
    curl -fsSL "$DOWNLOAD_URL" -o "$TARGET_PATH"
fi

chmod +x "$TARGET_PATH"

echo "✅ speedtest_cli installed successfully to ${TARGET_PATH}!"
echo ""
echo "Quick start:"
echo "  ${BINARY_NAME}          # Run full speedtest"
echo "  ${BINARY_NAME} --bytes  # Speed in MB/s"
echo "  ${BINARY_NAME} --help   # View all options"
echo ""
echo "🔗 Star us on GitHub: https://github.com/${REPO}"
