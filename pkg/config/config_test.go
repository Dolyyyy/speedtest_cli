package config

import (
	"testing"
)

func TestParseFlagsDefaults(t *testing.T) {
	cfg := ParseFlagsArgs([]string{})

	if cfg.Threads != 16 {
		t.Errorf("expected default threads to be 16, got %d", cfg.Threads)
	}

	if cfg.UseBytes {
		t.Errorf("expected default UseBytes to be false, got true")
	}

	if cfg.UseJSON {
		t.Errorf("expected default UseJSON to be false, got true")
	}
}

func TestParseFlagsThreadsBoundary(t *testing.T) {
	cfg := ParseFlagsArgs([]string{"--threads", "100"})

	if cfg.Threads != 64 {
		t.Errorf("expected max threads cap to be 64, got %d", cfg.Threads)
	}
}

func TestParseFlagsCustomServer(t *testing.T) {
	cfg := ParseFlagsArgs([]string{"--custom", "10.0.0.5:8080", "--bytes"})

	if cfg.CustomHost != "10.0.0.5:8080" {
		t.Errorf("expected custom host 10.0.0.5:8080, got %s", cfg.CustomHost)
	}

	if !cfg.UseBytes {
		t.Errorf("expected UseBytes to be true, got false")
	}
}
