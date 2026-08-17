package engine

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/Dolyyyy/speedtest_cli/pkg/ui"
	"github.com/showwin/speedtest-go/speedtest"
)

// NewRunner initializes a Runner instance with tuned high-throughput network parameters
func NewRunner(cfg *model.Config) *model.Runner {
	// High-throughput TCP Transport for 10G / 25G / 100G interfaces
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   500,
		MaxConnsPerHost:       500,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ReadBufferSize:        256 * 1024, // 256 KB TCP buffer
		WriteBufferSize:       256 * 1024, // 256 KB TCP buffer
	}

	userConfig := &speedtest.UserConfig{
		MaxConnections: cfg.Threads,
		T:              transport,
	}

	return &model.Runner{
		Cfg:    cfg,
		Client: speedtest.New(speedtest.WithUserConfig(userConfig)),
	}
}

// FetchServerItems retrieves nearby servers for --list flag
func FetchServerItems(r *model.Runner) ([]model.ServerItem, error) {
	serverList, err := r.Client.FetchServers()
	if err != nil {
		return nil, err
	}

	sort.Slice(serverList, func(i, j int) bool {
		return serverList[i].Distance < serverList[j].Distance
	})

	items := make([]model.ServerItem, len(serverList))
	for i, s := range serverList {
		items[i] = model.ServerItem{
			ID:       s.ID,
			Sponsor:  s.Sponsor,
			Name:     s.Name,
			Country:  s.Country,
			Distance: round(s.Distance, 1),
		}
	}
	return items, nil
}

// Run executes complete benchmark workflow
func Run(r *model.Runner) (*model.TestResult, error) {
	quiet := r.Cfg.UseJSON || r.Cfg.UseSimple

	// 1. User Info & Server Discovery
	spUser := ui.NewSpinner("Connecting & detecting connection metadata...", quiet)
	ui.StartSpinner(spUser)

	user, err := r.Client.FetchUserInfo()
	if err != nil {
		ui.FailSpinner(spUser, "Failed to retrieve connection information")
		return nil, err
	}

	serverList, err := r.Client.FetchServers()
	if err != nil || len(serverList) == 0 {
		ui.FailSpinner(spUser, "Failed to retrieve test servers")
		return nil, fmt.Errorf("no test servers available")
	}

	ui.StopSpinner(spUser, fmt.Sprintf("Connected : %s (%s) - IP: %s", user.Isp, user.Country, user.IP))

	// 2. Target Server Selection
	var targetServers speedtest.Servers
	if r.Cfg.CustomHost != "" {
		customServer, err := speedtest.CustomServer(r.Cfg.CustomHost)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to custom server %s: %w", r.Cfg.CustomHost, err)
		}
		targetServers = speedtest.Servers{customServer}
	} else if r.Cfg.ServerID != "" {
		sid, err := strconv.Atoi(r.Cfg.ServerID)
		if err == nil {
			for _, s := range serverList {
				if s.ID == r.Cfg.ServerID || s.ID == strconv.Itoa(sid) {
					targetServers = append(targetServers, s)
					break
				}
			}
		}
		if len(targetServers) == 0 {
			return nil, fmt.Errorf("server ID %s not found", r.Cfg.ServerID)
		}
	} else {
		spSearch := ui.NewSpinner("Selecting optimal test server (lowest latency)...", quiet)
		ui.StartSpinner(spSearch)

		targets, err := serverList.FindServer([]int{})
		if err != nil || len(targets) == 0 {
			ui.FailSpinner(spSearch, "Failed to select valid target server")
			return nil, fmt.Errorf("no valid target server found")
		}
		targetServers = targets
		ui.StopSpinner(spSearch, fmt.Sprintf("Target server : %s (%s, %s)", targets[0].Sponsor, targets[0].Name, targets[0].Country))
	}

	server := targetServers[0]

	ui.PrintHeader(quiet)

	// 3. Ping & Jitter Test
	spPing := ui.NewSpinner("Testing latency (Ping & Jitter)...", quiet)
	ui.StartSpinner(spPing)

	err = server.PingTest(nil)
	if err != nil {
		ui.FailSpinner(spPing, "Latency test failed")
		return nil, err
	}

	pingMs := float64(server.Latency.Milliseconds())
	if pingMs == 0 && server.Latency > 0 {
		pingMs = float64(server.Latency.Microseconds()) / 1000.0
	}
	jitterMs := float64(server.Jitter.Milliseconds())
	if jitterMs == 0 && server.Jitter > 0 {
		jitterMs = float64(server.Jitter.Microseconds()) / 1000.0
	}

	ui.StopSpinner(spPing, fmt.Sprintf("Ping: %.2f ms | Jitter: %.2f ms", pingMs, jitterMs))

	// 4. Download Test with multi-stream parallel threads & real-time live speed meter
	spDL := ui.NewSpinner(fmt.Sprintf("Testing download speed (%d parallel TCP streams)...", r.Cfg.Threads), quiet)
	spDL.RateGetter = func() float64 {
		return server.Context.GetEWMADownloadRate()
	}
	spDL.UseBytes = r.Cfg.UseBytes
	spDL.Threads = r.Cfg.Threads
	spDL.TestType = "download"
	ui.StartSpinner(spDL)

	err = server.DownloadTest()
	if err != nil {
		ui.FailSpinner(spDL, "Download test failed")
		return nil, err
	}

	dlMbps := server.DLSpeed.Mbps()
	dlMBps := dlMbps / 8.0

	var dlStr string
	if r.Cfg.UseBytes {
		dlStr = fmt.Sprintf("%.2f MB/s", dlMBps)
	} else {
		dlStr = fmt.Sprintf("%.2f Mbps", dlMbps)
	}
	ui.StopSpinner(spDL, fmt.Sprintf("Download : %s (%d streams)", dlStr, r.Cfg.Threads))

	// 5. Upload Test with multi-stream parallel threads & real-time live speed meter
	spUL := ui.NewSpinner(fmt.Sprintf("Testing upload speed (%d parallel TCP streams)...", r.Cfg.Threads), quiet)
	spUL.RateGetter = func() float64 {
		return server.Context.GetEWMAUploadRate()
	}
	spUL.UseBytes = r.Cfg.UseBytes
	spUL.Threads = r.Cfg.Threads
	spUL.TestType = "upload"
	ui.StartSpinner(spUL)

	err = server.UploadTest()
	if err != nil {
		ui.FailSpinner(spUL, "Upload test failed")
		return nil, err
	}

	ulMbps := server.ULSpeed.Mbps()
	ulMBps := ulMbps / 8.0

	var ulStr string
	if r.Cfg.UseBytes {
		ulStr = fmt.Sprintf("%.2f MB/s", ulMBps)
	} else {
		ulStr = fmt.Sprintf("%.2f Mbps", ulMbps)
	}
	ui.StopSpinner(spUL, fmt.Sprintf("Upload   : %s (%d streams)", ulStr, r.Cfg.Threads))

	result := &model.TestResult{
		Timestamp: time.Now(),
		Client: model.ClientInfo{
			IP:      user.IP,
			ISP:     user.Isp,
			Country: user.Country,
			Lat:     user.Lat,
			Lon:     user.Lon,
		},
		Server: model.ServerInfo{
			ID:       server.ID,
			Name:     server.Name,
			Sponsor:  server.Sponsor,
			Country:  server.Country,
			Distance: round(server.Distance, 2),
			Latency:  round(pingMs, 2),
		},
		PingMs:   round(pingMs, 2),
		JitterMs: round(jitterMs, 2),
		Download: model.SpeedVal{
			Mbps: round(dlMbps, 2),
			MBps: round(dlMBps, 2),
		},
		Upload: model.SpeedVal{
			Mbps: round(ulMbps, 2),
			MBps: round(ulMBps, 2),
		},
	}

	return result, nil
}

func round(val float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(val*p) / p
}
