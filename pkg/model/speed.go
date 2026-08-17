package model

// SpeedVal represents a speed value in Mbps and MB/s
type SpeedVal struct {
	Mbps float64 `json:"mbps"`
	MBps float64 `json:"mb_s"`
}
