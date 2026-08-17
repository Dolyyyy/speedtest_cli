package ui

import (
	"fmt"
	"strings"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/fatih/color"
)

// PrintDashboard displays futuristic dashboard card in English with Host Info and GitHub URL
func PrintDashboard(res *model.TestResult, useBytes bool) {
	border := color.CyanString("┌────────────────────────────────────────────────────────┐")
	sep := color.CyanString("├────────────────────────────────────────────────────────┤")
	bottom := color.CyanString("└────────────────────────────────────────────────────────┘")

	fmt.Println()
	fmt.Println(border)
	fmt.Printf("│  %s                   │\n", ColorTitle("📊 SPEEDTEST RESULTS"))
	fmt.Println(sep)
	fmt.Printf("│  %s %-45s │\n", ColorMuted("Host:    "), ColorVal(res.Host.String()))
	fmt.Printf("│  %s %-45s │\n", ColorMuted("Client:  "), ColorVal(fmt.Sprintf("%s (%s)", res.Client.ISP, res.Client.IP)))
	fmt.Printf("│  %s %-45s │\n", ColorMuted("Server:  "), ColorVal(fmt.Sprintf("%s - %s (%s)", res.Server.Sponsor, res.Server.Name, res.Server.Country)))
	fmt.Println(sep)
	fmt.Printf("│  ⚡ %s  %-39s │\n", ColorMuted("Ping:    "), ColorVal(fmt.Sprintf("%.2f ms  (Jitter: %.2f ms)", res.PingMs, res.JitterMs)))

	if useBytes {
		fmt.Printf("│  📥 %s  %-39s │\n", ColorMuted("Download:"), ColorSuccess(fmt.Sprintf("%.2f MB/s", res.Download.MBps)))
		fmt.Printf("│  📤 %s  %-39s │\n", ColorMuted("Upload:  "), ColorWarning(fmt.Sprintf("%.2f MB/s", res.Upload.MBps)))
	} else {
		fmt.Printf("│  📥 %s  %-39s │\n", ColorMuted("Download:"), ColorSuccess(fmt.Sprintf("%.2f Mbps", res.Download.Mbps)))
		fmt.Printf("│  📤 %s  %-39s │\n", ColorMuted("Upload:  "), ColorWarning(fmt.Sprintf("%.2f Mbps", res.Upload.Mbps)))
	}
	fmt.Println(sep)
	fmt.Printf("│  🔗 %-50s │\n", ColorMuted("https://github.com/Dolyyyy/speedtest_cli"))
	fmt.Println(bottom)
	fmt.Println()
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
