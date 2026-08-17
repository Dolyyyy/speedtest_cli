#!/usr/bin/env bash
set -e

REPO="Dolyyyy/speedtest_cli"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "❌ Unsupported arch: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux) BIN="speedtest-linux-${ARCH}" ;;
    darwin) BIN="speedtest-darwin-${ARCH}" ;;
    *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
esac

TMP_BIN="$(mktemp /tmp/speedtest.XXXXXX)"
trap 'rm -f "$TMP_BIN"' EXIT

curl -fsSL "https://github.com/${REPO}/raw/main/binaries/${BIN}?t=$(date +%s)" -o "$TMP_BIN"
chmod +x "$TMP_BIN"

"$TMP_BIN" "$@"
