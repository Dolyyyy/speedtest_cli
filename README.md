# speedtest_cli - Next-Gen High Performance Speedtest CLI

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-8A2BE2?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)
[![Architecture](https://img.shields.io/badge/Architecture-x86__64%20%7C%20ARM64-orange?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![Open Source](https://img.shields.io/badge/Open_Source-%E2%9D%A4-red?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)

A lightweight, multi-platform, high-performance network speed test CLI tool written in **Go**. Features a modern color terminal interface, automatic server selection, multi-stream benchmarking, latency/jitter diagnostics, and structured JSON output.

---

## 🚀 Instant 1-Command Installation

### Linux / WSL / macOS
```bash
curl -fsSL "https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.sh?t=$(date +%s)" | sudo bash
```

### Windows (PowerShell)
```powershell
iwr -useb "https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.ps1?t=$(get-date -UFormat %s)" | iex
```

Once installed, simply type `speedtest` in your terminal:
```bash
speedtest
```

---

## 🛠️ Usage & Examples

```bash
# Default: Fast automatic speedtest with modern TUI card
speedtest

# Display speed in Megabytes per second (MB/s)
speedtest --bytes

# List nearest available speedtest servers
speedtest --list

# Target a specific server by ID
speedtest -s 12345

# Export benchmark result in raw JSON for automated scripts
speedtest --json

# Minimal 3-line format (Ping, Download, Upload)
speedtest --simple
```

---

## 📋 Command Line Options

| Flag | Short | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display interactive CLI help guide |
| `--list` | `-l` | List available speedtest servers near your location |
| `--server <ID>` | `-s <ID>` | Run speedtest against a specific server ID |
| `--bytes` | | Display download/upload rates in MB/s (Mo/s) instead of Mbps |
| `--json` | | Output result in raw JSON format for scripting & logging |
| `--simple` | | Output concise 3-line format (Ping, Download, Upload) |
| `--version` | `-v` | Display `speedtest_cli` version |

---

## ⚡ Features

- **Blazing Fast**: Written in pure Go with zero runtime dependencies.
- **Single Command Setup**: Instant installation script downloads pre-compiled static binaries directly.
- **Multi-Architecture Support**: Built for Linux (amd64, arm64), macOS (Intel, Apple Silicon), and Windows (amd64).
- **Futuristic TUI**: Modern color palette, animated progress indicators, and clean summary dashboards.
- **Accurate Network Metrics**: Multi-threaded TCP download & upload testing with low overhead latency & jitter measurements.
- **Scripting Ready**: Supports `--json` and `--simple` outputs for seamless CI/CD integration and monitoring scripts.

---

## 📁 Repository Structure

```
speedtest_cli/
├── cmd/
│   └── speedtest/         # Binary entrypoint
├── pkg/
│   ├── config/            # CLI configuration & argument parsing
│   ├── engine/            # Speedtest engine & runner workflow
│   ├── model/             # Domain types (Client, Server, Speed, Result, Version)
│   ├── printer/           # Formatters (TUI, JSON, Simple)
│   └── ui/                # Modern TUI design (Colors, Spinner, Dashboard, Help)
├── binaries/              # Pre-compiled multi-platform static binaries
│   ├── speedtest-linux-amd64
│   ├── speedtest-linux-arm64
│   ├── speedtest-darwin-amd64
│   ├── speedtest-darwin-arm64
│   └── speedtest-windows-amd64.exe
├── build.sh               # Multi-platform compilation script
├── install.sh             # Linux/WSL/macOS 1-liner installer
├── install.ps1            # Windows PowerShell 1-liner installer
└── README.md
```

---

## 🛠️ Building from Source

Requirements: **Go 1.23+**

```bash
git clone https://github.com/Dolyyyy/speedtest_cli.git
cd speedtest_cli

# Build for current system
go build -o speedtest ./cmd/speedtest

# Build all cross-platform binaries into binaries/
./build.sh
```

---

## 📄 License

Distributed under the [MIT License](LICENSE).
