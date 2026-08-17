package model

import (
	"sync"
)

// Spinner represents an animated CLI loading indicator
type Spinner struct {
	Frames     []string
	Message    string
	Disabled   bool
	StopChan   chan struct{}
	Wg         sync.WaitGroup
	Mu         sync.Mutex
	RateGetter func() float64
	UseBytes   bool
	Threads    int
	TestType   string
}
