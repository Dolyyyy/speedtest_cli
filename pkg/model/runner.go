package model

import (
	"github.com/showwin/speedtest-go/speedtest"
)

// Runner orchestrates speedtest discovery and high-performance benchmarking
type Runner struct {
	Cfg    *Config
	Client *speedtest.Speedtest
}
