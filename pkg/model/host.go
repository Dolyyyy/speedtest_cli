package model

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// HostInfo represents metadata about the local machine executing speedtest
type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Kernel   string `json:"kernel,omitempty"`
	Arch     string `json:"arch"`
	CPUCores int    `json:"cpu_cores"`
	TotalRAM string `json:"total_ram,omitempty"`
}

// FetchHostInfo gathers local system metrics without external dependencies
func FetchHostInfo() HostInfo {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	info := HostInfo{
		Hostname: hostname,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCores: runtime.NumCPU(),
	}

	if runtime.GOOS == "linux" {
		if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
			info.Kernel = strings.TrimSpace(string(data))
		}
		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						if kb, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
							info.TotalRAM = fmt.Sprintf("%.1f GB", float64(kb)/(1024*1024))
						}
					}
					break
				}
			}
		}
	}

	return info
}

// String formats HostInfo as a single line description
func (h HostInfo) String() string {
	parts := []string{fmt.Sprintf("%s (%s/%s, %d CPU cores)", h.Hostname, h.OS, h.Arch, h.CPUCores)}
	if h.TotalRAM != "" {
		parts = append(parts, "RAM: "+h.TotalRAM)
	}
	return strings.Join(parts, " | ")
}
