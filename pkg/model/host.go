package model

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// HostInfo represents detailed metadata about the local machine executing speedtest
type HostInfo struct {
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	Kernel    string `json:"kernel,omitempty"`
	Arch      string `json:"arch"`
	CPUCores  int    `json:"cpu_cores"`
	TotalRAM  string `json:"total_ram,omitempty"`
	AvailRAM  string `json:"avail_ram,omitempty"`
	LocalIP   string `json:"local_ip,omitempty"`
	Interface string `json:"interface,omitempty"`
	LoadAvg   string `json:"load_avg,omitempty"`
	GoVersion string `json:"go_version"`
}

// FetchHostInfo gathers local system metrics without external dependencies
func FetchHostInfo() HostInfo {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}

	info := HostInfo{
		Hostname:  hostname,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		CPUCores:  runtime.NumCPU(),
		GoVersion: runtime.Version(),
	}

	ifaceName, localIP := getLocalNetworkInfo()
	info.Interface = ifaceName
	info.LocalIP = localIP

	if runtime.GOOS == "linux" {
		if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
			info.Kernel = strings.TrimSpace(string(data))
		}

		info.LoadAvg = getLinuxLoadAvg()

		totalRAM, availRAM := getLinuxMemStats()
		info.TotalRAM = totalRAM
		info.AvailRAM = availRAM
	}

	return info
}

func getLocalNetworkInfo() (ifaceName string, localIP string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", ""
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return iface.Name, ip.String()
			}
		}
	}
	return "", ""
}

func getLinuxLoadAvg() string {
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			return fmt.Sprintf("%s, %s, %s", fields[0], fields[1], fields[2])
		}
	}
	return ""
}

func getLinuxMemStats() (total string, avail string) {
	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		var totalKb, availKb uint64
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					totalKb, _ = strconv.ParseUint(fields[1], 10, 64)
				}
			} else if strings.HasPrefix(line, "MemAvailable:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					availKb, _ = strconv.ParseUint(fields[1], 10, 64)
				}
			}
		}
		if totalKb > 0 {
			total = fmt.Sprintf("%.1f GB", float64(totalKb)/(1024*1024))
		}
		if availKb > 0 {
			avail = fmt.Sprintf("%.1f GB", float64(availKb)/(1024*1024))
		}
	}
	return total, avail
}

// String formats HostInfo as a multi-line or clean inline description
func (h HostInfo) String() string {
	var hostDesc string
	if h.LocalIP != "" {
		hostDesc = fmt.Sprintf("%s (%s | %s/%s, %d CPU)", h.Hostname, h.LocalIP, h.OS, h.Arch, h.CPUCores)
	} else {
		hostDesc = fmt.Sprintf("%s (%s/%s, %d CPU)", h.Hostname, h.OS, h.Arch, h.CPUCores)
	}

	if h.TotalRAM != "" {
		if h.AvailRAM != "" {
			hostDesc += fmt.Sprintf(" | RAM: %s free / %s", h.AvailRAM, h.TotalRAM)
		} else {
			hostDesc += fmt.Sprintf(" | RAM: %s", h.TotalRAM)
		}
	}
	return hostDesc
}
