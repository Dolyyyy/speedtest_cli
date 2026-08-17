package engine

import (
	"testing"

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

func TestRound(t *testing.T) {
	val := round(12.34567, 2)
	if val != 12.35 {
		t.Errorf("expected rounded value 12.35, got %f", val)
	}
}
