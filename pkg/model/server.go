package model

// ServerInfo represents target speedtest server details
type ServerInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Sponsor  string  `json:"sponsor"`
	Country  string  `json:"country"`
	Distance float64 `json:"distance_km"`
	Latency  float64 `json:"latency_ms"`
}

// ServerItem represents a listed server in --list output
type ServerItem struct {
	ID       string  `json:"id"`
	Sponsor  string  `json:"sponsor"`
	Name     string  `json:"name"`
	Country  string  `json:"country"`
	Distance float64 `json:"distance_km"`
}
