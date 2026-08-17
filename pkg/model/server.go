package model

// ServerInfo contains details about the target speedtest server
type ServerInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Sponsor  string  `json:"sponsor"`
	Country  string  `json:"country"`
	Host     string  `json:"host,omitempty"`
	Distance float64 `json:"distance_km"`
	Latency  float64 `json:"latency_ms"`
}

// ServerItem represents a server summary for the --list command
type ServerItem struct {
	ID       string  `json:"id"`
	Sponsor  string  `json:"sponsor"`
	Name     string  `json:"name"`
	Country  string  `json:"country"`
	Host     string  `json:"host,omitempty"`
	Distance float64 `json:"distance_km"`
}
