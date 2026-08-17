package printer

import (
	"encoding/json"
	"fmt"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/Dolyyyy/speedtest_cli/pkg/ui"
)

// PrintResult outputs test result according to configuration flags
func PrintResult(res *model.TestResult, cfg *model.Config) {
	if cfg.UseJSON {
		data, err := json.MarshalIndent(res, "", "  ")
		if err == nil {
			fmt.Println(string(data))
		}
		return
	}

	if cfg.UseSimple {
		if cfg.UseBytes {
			fmt.Printf("Ping: %.2f ms\nDownload: %.2f MB/s\nUpload: %.2f MB/s\n", res.PingMs, res.Download.MBps, res.Upload.MBps)
		} else {
			fmt.Printf("Ping: %.2f ms\nDownload: %.2f Mbps\nUpload: %.2f Mbps\n", res.PingMs, res.Download.Mbps, res.Upload.Mbps)
		}
		return
	}

	ui.PrintDashboard(res, cfg.UseBytes)
}
