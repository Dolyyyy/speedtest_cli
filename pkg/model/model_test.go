package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFetchHostInfo(t *testing.T) {
	info := FetchHostInfo()
	if info.Hostname == "" {
		t.Error("expected non-empty hostname")
	}
	if info.OS == "" {
		t.Error("expected non-empty OS")
	}
	if info.CPUCores < 1 {
		t.Errorf("expected CPUCores >= 1, got %d", info.CPUCores)
	}

	str := info.String()
	if str == "" {
		t.Error("expected non-empty formatted HostInfo string")
	}
}

func TestJSONResultSerialization(t *testing.T) {
	res := TestResult{
		Timestamp: time.Now(),
		Host: HostInfo{
			Hostname: "test-server",
			OS:       "linux",
			Arch:     "amd64",
			CPUCores: 8,
			TotalRAM: "16.0 GB",
		},
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

	if unmarshaled.Host.Hostname != "test-server" {
		t.Errorf("expected host hostname test-server, got %s", unmarshaled.Host.Hostname)
	}
	if unmarshaled.Client.IP != "1.2.3.4" {
		t.Errorf("expected client IP 1.2.3.4, got %s", unmarshaled.Client.IP)
	}
	if unmarshaled.Download.Mbps != 950.5 {
		t.Errorf("expected download Mbps 950.5, got %f", unmarshaled.Download.Mbps)
	}
}
