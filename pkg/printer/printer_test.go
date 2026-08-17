package printer

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

func TestPrintResultSimple(t *testing.T) {
	res := &model.TestResult{
		Timestamp: time.Now(),
		PingMs:    12.34,
		Download:  model.SpeedVal{Mbps: 500.12, MBps: 62.515},
		Upload:    model.SpeedVal{Mbps: 250.56, MBps: 31.32},
	}
	cfg := &model.Config{UseSimple: true}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintResult(res, cfg)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Ping: 12.34 ms") {
		t.Errorf("expected output to contain Ping: 12.34 ms, got %s", output)
	}
	if !strings.Contains(output, "Download: 500.12 Mbps") {
		t.Errorf("expected output to contain Download: 500.12 Mbps, got %s", output)
	}
}
