package engine

import (
	"testing"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

func TestNewRunner(t *testing.T) {
	cfg := &model.Config{
		Threads: 16,
	}

	runner := NewRunner(cfg)
	if runner == nil {
		t.Fatal("expected runner instance to be created, got nil")
	}

	if runner.Cfg.Threads != 16 {
		t.Errorf("expected 16 threads in runner config, got %d", runner.Cfg.Threads)
	}
}

func TestFormatDurationMs(t *testing.T) {
	d := 15500 * time.Microsecond // 15.5 ms
	ms := formatDurationMs(d)
	if ms != 15.5 {
		t.Errorf("expected 15.5 ms, got %f", ms)
	}
}
