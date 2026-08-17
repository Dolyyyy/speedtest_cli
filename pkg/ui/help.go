package ui

import (
	"fmt"
	"strings"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/fatih/color"
)

// PrintHeader displays modern banner
func PrintHeader(quiet bool) {
	if quiet {
		return
	}
	fmt.Println()
	fmt.Printf("%s %s %s\n", ColorTitle("⚡"), color.New(color.Bold).Sprint("SPEEDTEST CLI"), ColorAccent("v"+model.Version))
	fmt.Println(color.CyanString(strings.Repeat("─", 48)))
}

// PrintHelp displays styled CLI help guide in English with GitHub URL
func PrintHelp() {
	fmt.Println()
	fmt.Printf("%s - High-Performance 10G/100G Capable Speedtest CLI in Go (v%s)\n", ColorTitle("⚡ speedtest_cli"), model.Version)
	fmt.Printf("%s\n\n", ColorMuted("🔗 https://github.com/Dolyyyy/speedtest_cli"))
	fmt.Printf("%s\n", ColorSuccess("USAGE:"))
	fmt.Println("  speedtest [options]")
	fmt.Println()
	fmt.Printf("%s\n", ColorSuccess("EXAMPLES:"))
	fmt.Println("  speedtest                    # Instant auto-benchmarking (16 parallel streams)")
	fmt.Println("  speedtest --threads 32       # Maximize multi-stream for 10G/25G/100G link saturation")
	fmt.Println("  speedtest --history          # View local benchmark history table")
	fmt.Println("  speedtest --bytes            # Display speed in MB/s instead of Mbps")
	fmt.Println("  speedtest --list             # List nearby servers")
	fmt.Println("  speedtest -s 12345           # Target server #12345")
	fmt.Println("  speedtest --custom 10.0.0.5  # Benchmark private infrastructure server")
	fmt.Println("  speedtest --json             # Export result in raw JSON format")
	fmt.Println()
	fmt.Printf("%s\n", ColorSuccess("OPTIONS:"))
	fmt.Printf("  %-25s %s\n", ColorWarning("-h, --help"), "Show this help guide")
	fmt.Printf("  %-25s %s\n", ColorWarning("-H, --history"), "Display local speedtest benchmark history")
	fmt.Printf("  %-25s %s\n", ColorWarning("--clear-history"), "Clear all saved local speedtest history records")
	fmt.Printf("  %-25s %s\n", ColorWarning("-t, --threads <N>"), "Number of parallel TCP connections for 10G/100G links (default: 16)")
	fmt.Printf("  %-25s %s\n", ColorWarning("-s, --server <ID>"), "Specify a target server ID")
	fmt.Printf("  %-25s %s\n", ColorWarning("--custom <host>"), "Benchmark custom speedtest server host (e.g., host:port)")
	fmt.Printf("  %-25s %s\n", ColorWarning("-l, --list"), "List nearest available speedtest servers")
	fmt.Printf("  %-25s %s\n", ColorWarning("--bytes"), "Display speed in MB/s instead of Mbps")
	fmt.Printf("  %-25s %s\n", ColorWarning("--json"), "Output results in JSON format")
	fmt.Printf("  %-25s %s\n", ColorWarning("--simple"), "Minimalist 3-line output format")
	fmt.Printf("  %-25s %s\n", ColorWarning("-v, --version"), "Display CLI version")
	fmt.Println()
}
