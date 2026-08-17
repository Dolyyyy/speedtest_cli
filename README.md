# ⚡ speedtest_cli - High Performance Speedtest CLI

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-8A2BE2?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)
[![Architecture](https://img.shields.io/badge/Architecture-x86__64%20%7C%20ARM64-orange?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)
[![CI](https://img.shields.io/badge/CI-Passing-brightgreen?style=for-the-badge&logo=githubactions&logoColor=white)](https://github.com/Dolyyyy/speedtest_cli/actions)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![Open Source](https://img.shields.io/badge/Open_Source-%E2%9D%A4-red?style=for-the-badge)](https://github.com/Dolyyyy/speedtest_cli)

A lightweight, multi-platform, high-performance network speed test CLI tool written in **Go**. Features a modern color terminal interface, unthrottled 100G Datacenter Backbone mode, automatic server discovery with 3x retry resilience, 10G/50G/100G multi-stream TCP saturation, zero-dependency host diagnostics, local benchmark history, and structured JSON output.

🔗 **Repository:** [https://github.com/Dolyyyy/speedtest_cli](https://github.com/Dolyyyy/speedtest_cli)

---

## 🚀 Quick Start (1-Line Commands)

### 📦 Install on System (Linux / macOS / WSL)
```bash
curl -fsSL https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.sh | bash
```

### ⚡ Run Instantly without Installing (Zero-Trace)
```bash
curl -fsSL https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/run.sh | bash
```

### 🪟 Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.ps1 | iex
```

> **Tip (Custom Name):** To install without overwriting an existing Ookla binary:  
> `curl -fsSL https://raw.githubusercontent.com/Dolyyyy/speedtest_cli/main/install.sh | BINARY_NAME=speedtest_cli bash`

---

## 🖥️ Terminal Preview

```text
✔ Connected : Fast Fiber ISP (FR) - IP: 198.51.100.42                    
✔ Target server : Scaleway (Paris, France) [ID: 61933 | Host/IP: st1.scaleway.com]

⚡ SPEEDTEST CLI v1.1.0
────────────────────────────────────────────────
✔ Ping: 1.80 ms | Jitter: 0.12 ms   
✔ Download : 1927.03 Mbps (16 streams)                 
✔ Upload   : 1980.22 Mbps (16 streams)               

╭────────────────────────────────────────────────────────────────────────╮
│ 📊 SPEEDTEST RESULTS                                                   │
├────────────────────────────────────────────────────────────────────────┤
│ 💻 Host:    server-01                                            │
│    System:  linux/amd64 (16 CPU Cores)                                 │
│    Memory:  20.5 GB free / 23.5 GB total                               │
│ 🌐 Client:  Fast Fiber ISP (198.51.100.42)                                    │
│ 📡 Server:  Scaleway (Ivry-sur-Seine, France - 1.0 km)   │
├────────────────────────────────────────────────────────────────────────┤
│ ⚡ Latency: 1.80 ms  (Jitter: 0.12 ms)                                │
│ 📥 Download: 1927.03 Mbps [████████████░░░░░░]                          │
│ 📤 Upload:   1980.22 Mbps [████████████░░░░░░]                          │
├────────────────────────────────────────────────────────────────────────┤
│ 🔗 https://github.com/Dolyyyy/speedtest_cli                            │
╰────────────────────────────────────────────────────────────────────────╯
```

---

## 🛠️ Usage & Examples

```bash
# Default: Fast auto-benchmarking (16 parallel TCP streams)
speedtest

# 100G Backbone Mode: Benchmark unthrottled 100G Datacenter Backbones
speedtest --100g

# High-Bandwidth: Maximize multi-stream saturation for 10G / 50G / 100G links
speedtest --threads 64

# View local benchmark history
speedtest --history

# Target a specific Speedtest server ID
speedtest -s 61933

# List available nearby servers with Host/IP and distance
speedtest --list

# Display bandwidth in Megabytes per second (MB/s)
speedtest --bytes

# Benchmark private custom infrastructure server
speedtest --custom 10.0.0.5:8080

# Clean machine-readable JSON output
speedtest --json

# Minimal 3-line format (Ping, Download, Upload)
speedtest --simple
```

---

## 📋 Command Line Options

| Flag | Short | Description |
| :--- | :--- | :--- |
| `--help` | `-h` | Display interactive CLI help guide |
| `--100g` | `-B` | Enable unthrottled 100G Datacenter Backbone benchmark mode |
| `--history` | `-H` | Display local benchmark history table |
| `--clear-history` | | Clear all saved local benchmark history records |
| `--threads <N>` | `-t <N>` | Number of parallel TCP streams for 10G/50G/100G links (1-128, default: 16) |
| `--server <ID>` | `-s <ID>` | Run speedtest against a specific server ID |
| `--custom <host>` | | Benchmark custom speedtest server host (e.g. `host:port`) |
| `--list` | `-l` | List available speedtest servers near your location |
| `--bytes` | | Display download/upload rates in MB/s (Mo/s) instead of Mbps |
| `--json` | | Output results in raw JSON format for scripting & pipelines |
| `--simple` | | Output concise 3-line format (Ping, Download, Upload) |
| `--version` | `-v` | Display `speedtest_cli` version |

---

## 📄 JSON Output Format (`--json`)

```json
{
  "timestamp": "2026-08-17T22:04:52+02:00",
  "host": {
    "hostname": "server-01",
    "os": "linux",
    "kernel": "6.6.137-microsoft-standard-WSL2",
    "arch": "amd64",
    "cpu_cores": 12,
    "total_ram": "15.6 GB",
    "avail_ram": "14.8 GB",
    "local_ip": "192.0.2.1",
    "interface": "eth0",
    "load_avg": "0.15, 0.10, 0.05",
    "go_version": "go1.23.0"
  },
  "client": {
    "ip": "198.51.100.42",
    "isp": "Fast Fiber ISP",
    "country": "FR"
  },
  "server": {
    "id": "69726",
    "name": "Paris",
    "sponsor": "AdKyNet SAS",
    "country": "France",
    "host": "speedtest.par1.adky.net",
    "distance_km": 1.0,
    "latency_ms": 1.13
  },
  "ping_ms": 1.13,
  "jitter_ms": 0.12,
  "download": {
    "mbps": 7528.89,
    "mb_s": 941.11
  },
  "upload": {
    "mbps": 4819.37,
    "mb_s": 602.42
  }
}
```

---

## ⚡ Core Architecture & Engineering

- **Zero External Dependencies**: Pure Go with stdlib metrics for CPU, RAM, OS, Kernel, and network interfaces.
- **100G Backbone Engine**: Multi-stream zero-copy workers targeting unthrottled datacenter backbones.
- **10G/50G/100G Socket Optimization**: 4MB TCP socket read/write buffers, `TCP_NODELAY`, and configurable multi-stream concurrency up to 128 threads.
- **Strict Domain Layer Separation**: 100% of data structures and domain models reside in `pkg/model/`.
- **Strategy Pattern Formatter**: Clean polymorphic rendering engine in `pkg/printer/` supporting TUI, JSON, and Simple formats.
- **Resilient Discovery & Latency Fallback**: 3x retry loop with exponential backoff on server discovery and latency measurement.
- **Ultra Lightweight**: Single static CGO-free binary (~6.0 MB) with near-instant execution.

---

## 📁 Repository Structure

```
speedtest_cli/
├── cmd/
│   └── speedtest/         # Binary CLI entrypoint
├── pkg/
│   ├── config/            # CLI configuration & flag argument parsing
│   ├── engine/            # Speedtest execution engine, 100G runner & Tier-1 discovery
│   ├── history/           # Local benchmark history storage (~/.speedtest_history.json)
│   ├── model/             # Domain models (Host, Client, Server, Speed, Result, Config, Spinner, Version)
│   ├── printer/           # Strategy pattern output formatters (TUI, JSON, Simple)
│   └── ui/                # Modern TUI design (Colors, Spinner, Dashboard, History, Help)
├── binaries/              # Pre-compiled static release binaries (Linux, macOS, Windows)
├── .github/
│   └── workflows/ci.yml   # Multi-platform GitHub Actions CI matrix
├── build.sh               # Cross-platform compilation script
├── install.sh             # Linux/WSL/macOS 1-liner installer
├── run.sh                 # Zero-trace 1-liner runner
├── install.ps1            # Windows PowerShell 1-liner installer
├── CONTRIBUTING.md        # Contribution guidelines
├── LICENSE                # MIT License
└── README.md
```

---

## 🛠️ Building from Source

Requirements: **Go 1.23+**

```bash
# Clone the repository
git clone https://github.com/Dolyyyy/speedtest_cli.git
cd speedtest_cli

# Run all unit tests
go test -v ./...

# Build binary for current platform
go build -o speedtest ./cmd/speedtest

# Build all cross-platform static release binaries into binaries/
./build.sh
```

---

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for pull request guidelines, code conventions, and testing instructions.

---

## 📄 License

Distributed under the [MIT License](LICENSE).
