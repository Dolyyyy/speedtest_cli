package ui

import (
	"testing"
)

func TestStripANSI(t *testing.T) {
	colored := "\x1b[36mHello World\x1b[0m"
	stripped := StripANSI(colored)
	if stripped != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", stripped)
	}

	titleColored := ColorTitle("📊 SPEEDTEST RESULTS")
	titleStripped := StripANSI(titleColored)
	if titleStripped != "📊 SPEEDTEST RESULTS" {
		t.Errorf("expected '📊 SPEEDTEST RESULTS', got '%s'", titleStripped)
	}
}
