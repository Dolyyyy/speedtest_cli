package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestJSONResultSerialization(t *testing.T) {
	res := TestResult{
		Timestamp: time.Now(),
		Client: ClientInfo{
			IP:      "1.2.3.4",
			ISP:     "TestISP",
			Country: "FR",
		},
		Server: ServerInfo{
			ID:       "100",
			Name:     "Paris",
			Sponsor:  "TestSponsor",
			Country:  "France",
			Distance: 10.5,
			Latency:  15.2,
		},
		PingMs:   15.2,
		JitterMs: 1.5,
		Download: SpeedVal{
			Mbps: 950.5,
			MBps: 118.81,
		},
		Upload: SpeedVal{
			Mbps: 450.2,
			MBps: 56.27,
		},
	}

	data, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("failed to marshal TestResult: %v", err)
	}

	var unmarshaled TestResult
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal TestResult: %v", err)
	}

	if unmarshaled.Client.IP != "1.2.3.4" {
		t.Errorf("expected client IP 1.2.3.4, got %s", unmarshaled.Client.IP)
	}

	if unmarshaled.Download.Mbps != 950.5 {
		t.Errorf("expected download Mbps 950.5, got %f", unmarshaled.Download.Mbps)
	}
}
