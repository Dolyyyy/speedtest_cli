package ui

import (
	"fmt"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

// NewSpinner creates a spinner instance
func NewSpinner(msg string, disabled bool) *model.Spinner {
	return &model.Spinner{
		Frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		Message:  msg,
		Disabled: disabled,
		StopChan: make(chan struct{}),
	}
}

// StartSpinner begins spinner animation with real-time speed ticker
func StartSpinner(s *model.Spinner) {
	if s.Disabled {
		return
	}
	s.Wg.Add(1)
	go func() {
		defer s.Wg.Done()
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-s.StopChan:
				return
			case <-ticker.C:
				s.Mu.Lock()
				msg := s.Message
				if s.RateGetter != nil {
					rateBps := s.RateGetter()
					if rateBps > 0 {
						if s.UseBytes {
							rateMBps := rateBps / 1000000.0
							msg = fmt.Sprintf("Testing %s speed (%d streams) - %s %.2f MB/s", s.TestType, s.Threads, ColorTitle("⚡"), rateMBps)
						} else {
							rateMbps := (rateBps * 8.0) / 1000000.0
							msg = fmt.Sprintf("Testing %s speed (%d streams) - %s %.2f Mbps", s.TestType, s.Threads, ColorTitle("⚡"), rateMbps)
						}
					}
				}
				fmt.Printf("\r\033[K%s  %s", ColorTitle(s.Frames[i%len(s.Frames)]), msg)
				s.Mu.Unlock()
				i++
			}
		}
	}()
}

// UpdateSpinnerMessage updates active spinner text
func UpdateSpinnerMessage(s *model.Spinner, msg string) {
	if s.Disabled {
		return
	}
	s.Mu.Lock()
	s.Message = msg
	s.Mu.Unlock()
}

// StopSpinner completes spinner execution with success icon
func StopSpinner(s *model.Spinner, successMsg string) {
	if s.Disabled {
		return
	}
	close(s.StopChan)
	s.Wg.Wait()
	fmt.Printf("\r\033[K%s %s\n", ColorSuccess("✔"), successMsg)
}

// FailSpinner completes spinner execution with error icon
func FailSpinner(s *model.Spinner, errMsg string) {
	if s.Disabled {
		return
	}
	close(s.StopChan)
	s.Wg.Wait()
	fmt.Printf("\r\033[K%s %s\n", ColorError("✖"), errMsg)
}
