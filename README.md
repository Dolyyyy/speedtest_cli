# speedtest_cli - Next-Gen High Performance Speedtest CLI

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-8A2BE2?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)
[![Architecture](https://img.shields.io/badge/Architecture-x86__64%20%7C%20ARM64-orange?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![Open Source](https://img.shields.io/badge/Open_Source-%E2%9D%A4-red?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)

A lightweight, multi-platform, high-performance network speed test CLI tool written in **Go**. Features a modern color terminal interface, automatic server selection, multi-stream 10G/100G benchmarking, latency/jitter diagnostics, local benchmark history, and structured JSON output.

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
# Default: Fast automatic speedtest with modern TUI card (16 parallel streams)
speedtest

# Maximize multi-stream TCP saturation for 10G / 25G / 100G links
speedtest --threads 32

# View local speedtest benchmark history
speedtest --history

# Display speed in Megabytes per second (MB/s)
speedtest --bytes

# List nearest available speedtest servers
speedtest --list

# Target a specific server by ID
speedtest -s 12345

# Benchmark custom infrastructure server
speedtest --custom 10.0.0.5:8080

# Export benchmark result in raw JSON format
speedtest --json

# Minimal 3-line format (Ping, Download, Upload)
speedtest --simple
```

---

## 📋 Command Line Options

| Flag | Short | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display interactive CLI help guide |
| `--history` | `-H` | Display local benchmark history table |
| `--clear-history` | | Clear all saved local benchmark history records |
| `--threads <N>` | `-t <N>` | Number of parallel TCP streams for 10G/100G links (default: 16) |
| `--server <ID>` | `-s <ID>` | Run speedtest against a specific server ID |
| `--custom <host>` | | Benchmark custom speedtest server host (e.g. host:port) |
| `--list` | `-l` | List available speedtest servers near your location |
| `--bytes` | | Display download/upload rates in MB/s (Mo/s) instead of Mbps |
| `--json` | | Output result in raw JSON format for scripting & logging |
| `--simple` | | Output concise 3-line format (Ping, Download, Upload) |
| `--version` | `-v` | Display `speedtest_cli` version |

---

## ⚡ Features

- **Blazing Fast**: Written in pure Go with zero runtime dependencies.
- **Single Command Setup**: Instant installation script downloads pre-compiled static binaries directly.
- **Multi-Architecture Support**: Built for Linux (amd64, arm64), macOS (Intel, Apple Silicon), and Windows (amd64).
- **Futuristic TUI**: Rounded neon borders, visual speed gauge progress bars, and pixel-perfect ANSI layout.
- **Local History**: View past benchmark results with `speedtest --history`.
- **10G/100G Network Saturation**: Tuned TCP socket transport with configurable parallel streams (`--threads`).
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
│   ├── history/           # Local benchmark history storage
│   ├── model/             # 100% domain types (Client, Server, Speed, Result, Host, History, Version)
│   ├── printer/           # Strategy pattern formatters (TUI, JSON, Simple)
│   └── ui/                # Modern TUI design (Colors, Spinner, Dashboard, History, Help)
├── binaries/              # Pre-compiled multi-platform static binaries
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
