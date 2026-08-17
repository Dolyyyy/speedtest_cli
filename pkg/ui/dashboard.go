package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

// stripANSI removes ANSI color code sequences from a string
func stripANSI(s string) string {
	var b strings.Builder
	inSeq := false
	for _, r := range s {
		if r == '\x1b' {
			inSeq = true
			continue
		}
		if inSeq {
			if (r >= 'a' && r <= 'Z') || (r >= 'A' && r <= 'Z') || r == '~' {
				inSeq = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// truncateVis truncates a string to a maximum visual cell width if needed
func truncateVis(s string, maxVisWidth int) string {
	if runewidth.StringWidth(s) <= maxVisWidth {
		return s
	}
	var b strings.Builder
	currW := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if currW+rw > maxVisWidth-3 {
			break
		}
		b.WriteRune(r)
		currW += rw
	}
	b.WriteString("...")
	return b.String()
}

// PrintDashboard displays futuristic, pixel-perfect dashboard card
func PrintDashboard(res *model.TestResult, useBytes bool) {
	const boxWidth = 72

	borderColor := color.New(color.FgHiCyan, color.Bold).SprintFunc()

	topBorder := borderColor("╭" + strings.Repeat("─", boxWidth+2) + "╮")
	sepBorder := borderColor("├" + strings.Repeat("─", boxWidth+2) + "┤")
	botBorder := borderColor("╰" + strings.Repeat("─", boxWidth+2) + "╯")

	printRow := func(content string) {
		plainText := stripANSI(content)
		visWidth := runewidth.StringWidth(plainText)
		pad := boxWidth - visWidth
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("%s %s%s %s\n", borderColor("│"), content, strings.Repeat(" ", pad), borderColor("│"))
	}

	fmt.Println()
	fmt.Println(topBorder)
	printRow(ColorTitle("📊 SPEEDTEST RESULTS"))
	fmt.Println(sepBorder)

	// Host & System Information (Multi-row for full clarity)
	if res.Host.Hostname != "" {
		printRow(fmt.Sprintf("%s %s", ColorMuted("💻 Host:    "), ColorVal(truncateVis(res.Host.Hostname, 58))))
		hostMeta := fmt.Sprintf("%s/%s (%d CPU Cores)", res.Host.OS, res.Host.Arch, res.Host.CPUCores)
		printRow(fmt.Sprintf("   %s %s", ColorMuted("OS/CPU:  "), ColorMuted(truncateVis(hostMeta, 58))))
		if res.Host.TotalRAM != "" {
			var ramMeta string
			if res.Host.AvailRAM != "" {
				ramMeta = fmt.Sprintf("%s free / %s total", res.Host.AvailRAM, res.Host.TotalRAM)
			} else {
				ramMeta = res.Host.TotalRAM
			}
			printRow(fmt.Sprintf("   %s %s", ColorMuted("Memory:  "), ColorMuted(truncateVis(ramMeta, 58))))
		}
	}

	// Client Information
	clientStr := fmt.Sprintf("%s (%s)", res.Client.ISP, res.Client.IP)
	printRow(fmt.Sprintf("%s %s", ColorMuted("🌐 Client:  "), ColorVal(truncateVis(clientStr, 58))))

	// Server Information
	serverLoc := res.Server.Name
	if res.Server.Country != "" {
		serverLoc += ", " + res.Server.Country
	}
	serverRaw := fmt.Sprintf("%s (%s - %.1f km)", res.Server.Sponsor, serverLoc, res.Server.Distance)
	printRow(fmt.Sprintf("%s %s", ColorMuted("📡 Server:  "), ColorVal(truncateVis(serverRaw, 58))))

	fmt.Println(sepBorder)

	// Latency
	pingStr := fmt.Sprintf("%.2f ms  (Jitter: %.2f ms)", res.PingMs, res.JitterMs)
	printRow(fmt.Sprintf("⚡ %s %s", ColorMuted("Latency: "), ColorVal(pingStr)))

	// Download & Upload Speeds with Visual Gauge Bar
	var dlSpeedText, ulSpeedText string
	if useBytes {
		dlSpeedText = fmt.Sprintf("%.2f MB/s", res.Download.MBps)
		ulSpeedText = fmt.Sprintf("%.2f MB/s", res.Upload.MBps)
	} else {
		dlSpeedText = fmt.Sprintf("%.2f Mbps", res.Download.Mbps)
		ulSpeedText = fmt.Sprintf("%.2f Mbps", res.Upload.Mbps)
	}

	dlBar := renderSpeedBar(res.Download.Mbps, color.FgHiGreen)
	ulBar := renderSpeedBar(res.Upload.Mbps, color.FgHiYellow)

	printRow(fmt.Sprintf("📥 %s %-14s %s", ColorMuted("Download:"), ColorSuccess(dlSpeedText), dlBar))
	printRow(fmt.Sprintf("📤 %s %-14s %s", ColorMuted("Upload:  "), ColorWarning(ulSpeedText), ulBar))

	fmt.Println(sepBorder)
	printRow(fmt.Sprintf("🔗 %s", ColorMuted("https://github.com/Dolyyyy/speedtest_cli")))
	fmt.Println(botBorder)
	fmt.Println()
}

// renderSpeedBar renders a visual progress gauge bar
func renderSpeedBar(mbps float64, barColor color.Attribute) string {
	maxScale := 1000.0
	if mbps > 1000.0 {
		maxScale = 10000.0
	}
	if mbps > 10000.0 {
		maxScale = 100000.0
	}

	ratio := mbps / maxScale
	if ratio > 1.0 {
		ratio = 1.0
	}

	barWidth := 18
	filled := int(math.Round(ratio * float64(barWidth)))
	if filled == 0 && mbps > 0 {
		filled = 1
	}

	cFunc := color.New(barColor).SprintFunc()
	cMuted := color.New(color.FgHiBlack).SprintFunc()

	filledBar := cFunc(strings.Repeat("█", filled))
	emptyBar := cMuted(strings.Repeat("░", barWidth-filled))

	return "[" + filledBar + emptyBar + "]"
}

// PrintServerList prints table of available servers in English
func PrintServerList(servers []model.ServerItem) {
	fmt.Println()
	fmt.Printf("%s\n\n", ColorTitle("📍 AVAILABLE SERVERS:"))
	fmt.Printf("%-8s %-25s %-20s %-10s\n", ColorHeader("ID"), ColorHeader("SPONSOR"), ColorHeader("LOCATION"), ColorHeader("DISTANCE"))
	fmt.Println(strings.Repeat("─", 68))

	maxShow := 15
	if len(servers) < maxShow {
		maxShow = len(servers)
	}

	for i := 0; i < maxShow; i++ {
		s := servers[i]
		fmt.Printf("%-8s %-25.25s %-20.20s %.1f km\n", s.ID, s.Sponsor, s.Name+" ("+s.Country+")", s.Distance)
	}
	fmt.Println()
	fmt.Printf("%s\n\n", ColorMuted("🔗 https://github.com/Dolyyyy/speedtest_cli"))
}

// PrintHistoryTable prints formatted table of local speedtest history
func PrintHistoryTable(items []model.HistoryItem, useBytes bool) {
	fmt.Println()
	fmt.Printf("%s\n\n", ColorTitle("📜 BENCHMARK HISTORY:"))
	if len(items) == 0 {
		fmt.Println(ColorMuted("No saved speedtest history records found."))
		fmt.Println()
		return
	}

	if useBytes {
		fmt.Printf("%-20s %-25s %-10s %-12s %-12s\n", ColorHeader("DATE/TIME"), ColorHeader("SERVER"), ColorHeader("PING"), ColorHeader("DOWNLOAD"), ColorHeader("UPLOAD"))
	} else {
		fmt.Printf("%-20s %-25s %-10s %-12s %-12s\n", ColorHeader("DATE/TIME"), ColorHeader("SERVER"), ColorHeader("PING"), ColorHeader("DOWNLOAD"), ColorHeader("UPLOAD"))
	}
	fmt.Println(strings.Repeat("─", 82))

	for _, item := range items {
		var dlStr, ulStr string
		if useBytes {
			dlStr = fmt.Sprintf("%.2f MB/s", item.Download.MBps)
			ulStr = fmt.Sprintf("%.2f MB/s", item.Upload.MBps)
		} else {
			dlStr = fmt.Sprintf("%.2f Mbps", item.Download.Mbps)
			ulStr = fmt.Sprintf("%.2f Mbps", item.Upload.Mbps)
		}

		serverStr := truncateVis(item.Server, 24)
		pingStr := fmt.Sprintf("%.2f ms", item.PingMs)

		fmt.Printf("%-20s %-25s %-10s %-12s %-12s\n", item.Timestamp, serverStr, pingStr, ColorSuccess(dlStr), ColorWarning(ulStr))
	}
	fmt.Println()
	fmt.Printf("%s\n\n", ColorMuted("🔗 https://github.com/Dolyyyy/speedtest_cli"))
}
