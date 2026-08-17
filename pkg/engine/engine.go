package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/Dolyyyy/speedtest_cli/pkg/ui"
	"github.com/showwin/speedtest-go/speedtest"
)

type ipApiMeta struct {
	Status      string `json:"status"`
	CountryCode string `json:"countryCode"`
	ISP         string `json:"isp"`
	Org         string `json:"org"`
	Query       string `json:"query"`
}

// NewRunner initializes a Runner instance with tuned high-throughput network parameters
func NewRunner(cfg *model.Config) *model.Runner {
	// Force IPv4 tcp4 DialContext with fallback to dual-stack if tcp4 is unavailable
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			// Attempt IPv4 connection first
			conn, err := dialer.DialContext(ctx, "tcp4", addr)
			if err == nil {
				return conn, nil
			}
			// Fallback to standard network if tcp4 unavailable (IPv6-only host)
			return dialer.DialContext(ctx, network, addr)
		},
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

// fetchIPv4ClientInfo fetches accurate IPv4 address and ISP details
func fetchIPv4ClientInfo() (*model.ClientInfo, error) {
	client := &http.Client{
		Timeout: 4 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return (&net.Dialer{Timeout: 3 * time.Second}).DialContext(ctx, "tcp4", addr)
			},
		},
	}

	resp, err := client.Get("http://ip-api.com/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var meta ipApiMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil || meta.Status != "success" {
		return nil, fmt.Errorf("invalid response")
	}

	isp := meta.ISP
	if meta.Org != "" && !strings.Contains(meta.ISP, meta.Org) {
		isp = meta.ISP + " (" + meta.Org + ")"
	}
	
	// Normalize common French ISP names for clean display
	lowerIsp := strings.ToLower(isp)
	if strings.Contains(lowerIsp, "orange") {
		isp = "Orange"
	} else if strings.Contains(lowerIsp, "sfr") {
		isp = "SFR"
	} else if strings.Contains(lowerIsp, "free") {
		isp = "Free"
	} else if strings.Contains(lowerIsp, "bouygues") {
		isp = "Bouygues Telecom"
	}

	return &model.ClientInfo{
		IP:      meta.Query,
		ISP:     isp,
		Country: meta.CountryCode,
	}, nil
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
			Distance: model.Round(s.Distance, 1),
		}
	}
	return items, nil
}

// Run executes complete benchmark workflow
func Run(r *model.Runner) (*model.TestResult, error) {
	quiet := r.Cfg.IsQuiet()

	// 1. User Info & Server Discovery (IPv4 Priority)
	spUser := ui.NewSpinner("Connecting & detecting connection metadata (IPv4 preference)...", quiet)
	ui.StartSpinner(spUser)

	var clientInfo model.ClientInfo

	// Attempt IPv4 metadata lookup first
	v4Info, errV4 := fetchIPv4ClientInfo()
	if errV4 == nil && v4Info != nil && v4Info.IP != "" {
		clientInfo = *v4Info
	} else {
		// Fallback to speedtest.FetchUserInfo()
		user, err := r.Client.FetchUserInfo()
		if err != nil {
			ui.FailSpinner(spUser, "Failed to retrieve connection information")
			return nil, err
		}
		clientInfo = model.ClientInfo{
			IP:      user.IP,
			ISP:     user.Isp,
			Country: user.Country,
			Lat:     user.Lat,
			Lon:     user.Lon,
		}
	}

	serverList, err := r.Client.FetchServers()
	if err != nil || len(serverList) == 0 {
		ui.FailSpinner(spUser, "Failed to retrieve test servers")
		return nil, fmt.Errorf("no test servers available")
	}

	ui.StopSpinner(spUser, fmt.Sprintf("Connected : %s (%s) - IP: %s", clientInfo.ISP, clientInfo.Country, clientInfo.IP))

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

	pingMs := formatDurationMs(server.Latency)
	jitterMs := formatDurationMs(server.Jitter)

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

	dlSpeed := model.NewSpeedVal(server.DLSpeed.Mbps())
	ui.StopSpinner(spDL, fmt.Sprintf("Download : %s (%d streams)", dlSpeed.String(r.Cfg.UseBytes), r.Cfg.Threads))

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

	ulSpeed := model.NewSpeedVal(server.ULSpeed.Mbps())
	ui.StopSpinner(spUL, fmt.Sprintf("Upload   : %s (%d streams)", ulSpeed.String(r.Cfg.UseBytes), r.Cfg.Threads))

	result := &model.TestResult{
		Timestamp: time.Now(),
		Client:    clientInfo,
		Server: model.ServerInfo{
			ID:       server.ID,
			Name:     server.Name,
			Sponsor:  server.Sponsor,
			Country:  server.Country,
			Distance: model.Round(server.Distance, 2),
			Latency:  model.Round(pingMs, 2),
		},
		PingMs:   model.Round(pingMs, 2),
		JitterMs: model.Round(jitterMs, 2),
		Download: dlSpeed,
		Upload:   ulSpeed,
	}

	return result, nil
}

func formatDurationMs(d time.Duration) float64 {
	if d <= 0 {
		return 0
	}
	return float64(d.Nanoseconds()) / 1e6
}
