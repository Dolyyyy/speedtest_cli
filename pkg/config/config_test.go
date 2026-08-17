package config

import (
	"testing"
)

func TestParseFlagsDefaults(t *testing.T) {
	cfg := ParseFlagsArgs([]string{})
	if cfg.Threads != 16 {
		t.Errorf("expected default threads 16, got %d", cfg.Threads)
	}
	if cfg.UseBytes != false {
		t.Errorf("expected default useBytes false, got %v", cfg.UseBytes)
	}
}

func TestParseFlagsThreadsBoundary(t *testing.T) {
	cfgLow := ParseFlagsArgs([]string{"--threads", "0"})
	if cfgLow.Threads != 1 {
		t.Errorf("expected min threads cap to be 1, got %d", cfgLow.Threads)
	}

	cfgHigh := ParseFlagsArgs([]string{"--threads", "200"})
	if cfgHigh.Threads != 128 {
		t.Errorf("expected max threads cap to be 128, got %d", cfgHigh.Threads)
	}
}

func TestParseFlagsCustomServer(t *testing.T) {
	cfg := ParseFlagsArgs([]string{"--custom", "10.0.0.1:8080"})
	if cfg.CustomHost != "10.0.0.1:8080" {
		t.Errorf("expected custom host 10.0.0.1:8080, got %s", cfg.CustomHost)
	}
}
