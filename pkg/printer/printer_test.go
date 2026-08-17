package printer

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

func TestJSONPrinter(t *testing.T) {
	res := &model.TestResult{
		Timestamp: time.Now(),
		PingMs:    10.5,
		Download:  model.NewSpeedVal(100.0),
		Upload:    model.NewSpeedVal(50.0),
	}

	var buf bytes.Buffer
	p := &JSONPrinter{Out: &buf}
	err := p.Print(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"ping_ms": 10.5`) {
		t.Errorf("expected JSON to contain ping_ms 10.5, got %s", output)
	}
}

func TestSimplePrinterMbps(t *testing.T) {
	res := &model.TestResult{
		Timestamp: time.Now(),
		PingMs:    12.34,
		Download:  model.NewSpeedVal(500.12),
		Upload:    model.NewSpeedVal(250.56),
	}

	var buf bytes.Buffer
	p := &SimplePrinter{Out: &buf, UseBytes: false}
	err := p.Print(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Ping: 12.34 ms") {
		t.Errorf("expected output to contain Ping: 12.34 ms, got %s", output)
	}
	if !strings.Contains(output, "Download: 500.12 Mbps") {
		t.Errorf("expected output to contain Download: 500.12 Mbps, got %s", output)
	}
}

func TestSimplePrinterMBps(t *testing.T) {
	res := &model.TestResult{
		Timestamp: time.Now(),
		PingMs:    12.34,
		Download:  model.NewSpeedVal(800.0), // 100 MB/s
		Upload:    model.NewSpeedVal(400.0), // 50 MB/s
	}

	var buf bytes.Buffer
	p := &SimplePrinter{Out: &buf, UseBytes: true}
	err := p.Print(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Download: 100.00 MB/s") {
		t.Errorf("expected output to contain Download: 100.00 MB/s, got %s", output)
	}
}
