package history

import (
	"testing"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

func TestSaveAndLoadHistory(t *testing.T) {
	_ = ClearHistory()

	res := &model.TestResult{
		Timestamp: time.Now(),
		Server: model.ServerInfo{
			Sponsor: "TestSponsor",
			Name:    "TestCity",
		},
		PingMs:   12.34,
		JitterMs: 1.2,
		Download: model.NewSpeedVal(500.0),
		Upload:   model.NewSpeedVal(250.0),
	}

	err := SaveResult(res)
	if err != nil {
		t.Fatalf("failed to save result: %v", err)
	}

	items, err := LoadHistory()
	if err != nil {
		t.Fatalf("failed to load history: %v", err)
	}

	if len(items) == 0 {
		t.Fatal("expected at least 1 history item, got 0")
	}

	if items[0].Download.Mbps != 500.0 {
		t.Errorf("expected download 500.0 Mbps, got %f", items[0].Download.Mbps)
	}

	_ = ClearHistory()
}
