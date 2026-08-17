package model

// HistoryItem represents a single saved speedtest record
type HistoryItem struct {
	Timestamp string   `json:"timestamp"`
	Server    string   `json:"server"`
	PingMs    float64  `json:"ping_ms"`
	JitterMs  float64  `json:"jitter_ms"`
	Download  SpeedVal `json:"download"`
	Upload    SpeedVal `json:"upload"`
}
