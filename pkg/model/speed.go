package model

import "fmt"

// SpeedVal represents a speed value in Mbps and MB/s
type SpeedVal struct {
	Mbps float64 `json:"mbps"`
	MBps float64 `json:"mb_s"`
}

// NewSpeedVal constructs a SpeedVal from Mbps and computes MB/s
func NewSpeedVal(mbps float64) SpeedVal {
	return SpeedVal{
		Mbps: Round(mbps, 2),
		MBps: Round(mbps/8.0, 2),
	}
}

// String formats SpeedVal according to byte flag preference
func (s SpeedVal) String(useBytes bool) string {
	if useBytes {
		return fmt.Sprintf("%.2f MB/s", s.MBps)
	}
	return fmt.Sprintf("%.2f Mbps", s.Mbps)
}
