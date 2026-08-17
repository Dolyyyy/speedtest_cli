package model

import "fmt"

// SpeedVal represents network throughput in both Mbps and MB/s
type SpeedVal struct {
	Mbps float64 `json:"mbps"`
	MBps float64 `json:"mb_s"`
}

// NewSpeedVal constructs SpeedVal with automatic MB/s calculation
func NewSpeedVal(mbps float64) SpeedVal {
	if mbps < 0 {
		mbps = 0
	}
	return SpeedVal{
		Mbps: Round(mbps, 2),
		MBps: Round(mbps/8.0, 2),
	}
}

// String formats speed metric based on unit preference
func (s SpeedVal) String(useBytes bool) string {
	if useBytes {
		return fmt.Sprintf("%.2f MB/s", s.MBps)
	}
	return fmt.Sprintf("%.2f Mbps", s.Mbps)
}
