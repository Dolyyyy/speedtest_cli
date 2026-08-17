# Contributing to speedtest_cli

Thank you for your interest in contributing to **speedtest_cli**! We welcome contributions from developers of all skill levels.

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/dl/) installed on your machine.
- Git.

### Fork & Clone

1. Fork the repository on GitHub: `https://github.com/Dolyyyy/speedtest_cli`
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/speedtest_cli.git
   cd speedtest_cli
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/Dolyyyy/speedtest_cli.git
   ```

---

## 🛠️ Development & Building

### Project Layout

```
speedtest_cli/
├── cmd/
│   └── speedtest/         # Application entrypoint
├── pkg/
│   ├── config/            # Flag parsing & configuration
│   ├── engine/            # Speedtest runner & server discovery
│   ├── model/             # 100% of domain structs & data models
│   ├── printer/           # Output formatters (TUI, JSON, Simple)
│   └── ui/                # TUI elements, colors, animated spinner
├── binaries/              # Multi-platform compiled static binaries
├── build.sh               # Multi-platform cross-compilation script
└── .github/workflows/     # GitHub Actions CI/CD workflow
```

### Build Locally

```bash
go build -o speedtest ./cmd/speedtest
./speedtest --help
```

---

## 🧪 Testing Requirements

Before submitting any Pull Request, you **must** ensure that all unit tests pass and code compiles cleanly without race conditions.

### Run Unit Tests

```bash
go test -v -race ./...
```

### Verify Multi-Platform Build

Ensure cross-compilation for all architectures succeeds:

```bash
chmod +x build.sh
./build.sh
```

---

## 🔄 Pull Request Guidelines

1. **Create a Feature Branch**:
   ```bash
   git checkout -b feature/my-cool-feature
   ```

2. **Follow Go Design Practices**:
   - Keep domain types exclusively in `pkg/model/`.
   - Maintain clear separation between UI, configuration, engine, and printer packages.
   - Write unit tests for new logic in `*_test.go` files.

3. **Pass GitHub Actions CI/CD**:
   - All PRs are automatically checked by GitHub Actions across **Linux**, **macOS**, and **Windows**.
   - Your PR will only be merged once all CI unit tests and multi-platform compilation checks pass cleanly.

4. **Commit Conventions**:
   Use clear, conventional commit messages:
   - `feat: ...` for new features
   - `fix: ...` for bug fixes
   - `test: ...` for adding tests
   - `docs: ...` for documentation updates

---

## 📄 License

By contributing to `speedtest_cli`, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
