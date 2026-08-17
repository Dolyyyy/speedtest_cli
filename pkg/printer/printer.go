package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/Dolyyyy/speedtest_cli/pkg/ui"
)

// ResultPrinter defines interface for output formatting strategies
type ResultPrinter interface {
	Print(res *model.TestResult) error
}

// JSONPrinter outputs raw JSON
type JSONPrinter struct {
	Out io.Writer
}

func (p *JSONPrinter) Print(res *model.TestResult) error {
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(p.Out, string(data))
	return err
}

// SimplePrinter outputs minimalist text lines
type SimplePrinter struct {
	Out      io.Writer
	UseBytes bool
}

func (p *SimplePrinter) Print(res *model.TestResult) error {
	if p.UseBytes {
		_, err := fmt.Fprintf(p.Out, "Ping: %.2f ms\nDownload: %.2f MB/s\nUpload: %.2f MB/s\n", res.PingMs, res.Download.MBps, res.Upload.MBps)
		return err
	}
	_, err := fmt.Fprintf(p.Out, "Ping: %.2f ms\nDownload: %.2f Mbps\nUpload: %.2f Mbps\n", res.PingMs, res.Download.Mbps, res.Upload.Mbps)
	return err
}

// TUIPrinter outputs futuristic color dashboard
type TUIPrinter struct {
	UseBytes bool
}

func (p *TUIPrinter) Print(res *model.TestResult) error {
	ui.PrintDashboard(res, p.UseBytes)
	return nil
}

// GetPrinter returns appropriate ResultPrinter based on configuration
func GetPrinter(cfg *model.Config) ResultPrinter {
	if cfg.UseJSON {
		return &JSONPrinter{Out: os.Stdout}
	}
	if cfg.UseSimple {
		return &SimplePrinter{Out: os.Stdout, UseBytes: cfg.UseBytes}
	}
	return &TUIPrinter{UseBytes: cfg.UseBytes}
}

// PrintResult delegates execution to configured printer strategy
func PrintResult(res *model.TestResult, cfg *model.Config) error {
	p := GetPrinter(cfg)
	return p.Print(res)
}
