package model

import "time"

// TestResult encapsulates complete benchmark execution data
type TestResult struct {
	Timestamp time.Time  `json:"timestamp"`
	Host      HostInfo   `json:"host"`
	Client    ClientInfo `json:"client"`
	Server    ServerInfo `json:"server"`
	PingMs    float64    `json:"ping_ms"`
	JitterMs  float64    `json:"jitter_ms"`
	Download  SpeedVal   `json:"download"`
	Upload    SpeedVal   `json:"upload"`
}
