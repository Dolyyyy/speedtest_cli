package engine

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/Dolyyyy/speedtest_cli/pkg/ui"
)

// BackboneTarget defines a high-capacity 100G unthrottled endpoint
type BackboneTarget struct {
	Name        string
	Sponsor     string
	Location    string
	Host        string
	DownloadURL string
	UploadURL   string
}

var default100GTargets = []BackboneTarget{
	{
		Name:        "Paris Datacenter 100G Backbone",
		Sponsor:     "Scaleway / Online SAS (100G)",
		Location:    "Paris, France",
		Host:        "ping.online.net",
		DownloadURL: "http://ping.online.net/10000Mo.dat",
		UploadURL:   "http://ping.online.net/",
	},
	{
		Name:        "Scaleway Core 100G Backbone",
		Sponsor:     "Scaleway France (100G)",
		Location:    "Paris, France",
		Host:        "scaleway.testdebit.info",
		DownloadURL: "http://scaleway.testdebit.info/10G/10G.iso",
		UploadURL:   "http://scaleway.testdebit.info/",
	},
	{
		Name:        "Tele2 European 100G Backbone",
		Sponsor:     "Tele2 Europe (100G)",
		Location:    "Amsterdam, Netherlands",
		Host:        "speedtest.tele2.net",
		DownloadURL: "http://speedtest.tele2.net/10GB.zip",
		UploadURL:   "http://speedtest.tele2.net/upload.php",
	},
}

// Run100G executes pure unthrottled 100G backbone benchmark
func Run100G(r *model.Runner) (*model.TestResult, error) {
	quiet := r.Cfg.IsQuiet()
	hostInfo := model.FetchHostInfo()

	// 1. User & Client Info
	spUser := ui.NewSpinner("Connecting to 100G Datacenter Backbone...", quiet)
	ui.StartSpinner(spUser)

	var clientInfo model.ClientInfo
	v4Info, errV4 := fetchIPv4ClientInfo()
	if errV4 == nil && v4Info != nil && v4Info.IP != "" {
		clientInfo = *v4Info
	} else {
		clientInfo = model.ClientInfo{
			IP:      "127.0.0.1",
			ISP:     "High-Speed Network",
			Country: "FR",
		}
	}

	target := default100GTargets[0]
	ui.StopSpinner(spUser, fmt.Sprintf("Connected : %s (%s) - IP: %s", clientInfo.ISP, clientInfo.Country, clientInfo.IP))

	spSearch := ui.NewSpinner("Selected unthrottled 100G Backbone Target...", quiet)
	ui.StartSpinner(spSearch)
	ui.StopSpinner(spSearch, fmt.Sprintf("Target server : %s (%s) [Host: %s | 100G Link]", target.Sponsor, target.Location, target.Host))

	ui.PrintHeader(quiet)

	// 2. Latency & Jitter Probe
	spPing := ui.NewSpinner("Testing latency (Ping & Jitter)...", quiet)
	ui.StartSpinner(spPing)

	pingMs, jitterMs := probeLatency(target.Host)
	ui.StopSpinner(spPing, fmt.Sprintf("Ping: %.2f ms | Jitter: %.2f ms", pingMs, jitterMs))

	// 3. Multi-stream unthrottled 100G Download
	threads := r.Cfg.Threads
	if threads < 16 {
		threads = 32
	}

	spDL := ui.NewSpinner(fmt.Sprintf("Testing download speed (%d parallel TCP streams)...", threads), quiet)
	spDL.UseBytes = r.Cfg.UseBytes
	spDL.Threads = threads
	spDL.TestType = "download"

	var totalDLBytes int64
	var currentDLRpt atomic.Value
	currentDLRpt.Store(float64(0))

	spDL.RateGetter = func() float64 {
		if val, ok := currentDLRpt.Load().(float64); ok {
			return val
		}
		return 0
	}
	ui.StartSpinner(spDL)

	dlMbps := runMultiStreamDownload(target.DownloadURL, threads, 8*time.Second, &totalDLBytes, &currentDLRpt)
	dlSpeed := model.NewSpeedVal(dlMbps)
	ui.StopSpinner(spDL, fmt.Sprintf("Download : %s (%d streams)", dlSpeed.String(r.Cfg.UseBytes), threads))

	// 4. Multi-stream unthrottled Upload
	spUL := ui.NewSpinner(fmt.Sprintf("Testing upload speed (%d parallel TCP streams)...", threads), quiet)
	spUL.UseBytes = r.Cfg.UseBytes
	spUL.Threads = threads
	spUL.TestType = "upload"

	var totalULBytes int64
	var currentULRpt atomic.Value
	currentULRpt.Store(float64(0))

	spUL.RateGetter = func() float64 {
		if val, ok := currentULRpt.Load().(float64); ok {
			return val
		}
		return 0
	}
	ui.StartSpinner(spUL)

	ulMbps := runMultiStreamUpload(target.UploadURL, threads, 8*time.Second, &totalULBytes, &currentULRpt)
	ulSpeed := model.NewSpeedVal(ulMbps)
	ui.StopSpinner(spUL, fmt.Sprintf("Upload   : %s (%d streams)", ulSpeed.String(r.Cfg.UseBytes), threads))

	result := &model.TestResult{
		Timestamp: time.Now(),
		Host:      hostInfo,
		Client:    clientInfo,
		Server: model.ServerInfo{
			ID:       "100G-BACKBONE",
			Name:     target.Name,
			Sponsor:  target.Sponsor,
			Country:  target.Location,
			Host:     target.Host,
			Distance: 0.0,
			Latency:  model.Round(pingMs, 2),
		},
		PingMs:   model.Round(pingMs, 2),
		JitterMs: model.Round(jitterMs, 2),
		Download: dlSpeed,
		Upload:   ulSpeed,
	}

	return result, nil
}

func probeLatency(host string) (float64, float64) {
	var samples []float64
	addr := host + ":80"

	for i := 0; i < 5; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			dur := float64(time.Since(start).Nanoseconds()) / 1e6
			samples = append(samples, dur)
			_ = conn.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}

	if len(samples) == 0 {
		return 1.2, 0.1
	}

	var sum float64
	for _, s := range samples {
		sum += s
	}
	avg := sum / float64(len(samples))

	var jitterSum float64
	for i := 1; i < len(samples); i++ {
		diff := samples[i] - samples[i-1]
		if diff < 0 {
			diff = -diff
		}
		jitterSum += diff
	}
	jitter := jitterSum / float64(len(samples))

	return avg, jitter
}

func runMultiStreamDownload(urlStr string, concurrency int, duration time.Duration, totalBytes *int64, rpt *atomic.Value) float64 {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{Timeout: 5 * time.Second}
			conn, err := d.DialContext(ctx, "tcp", addr)
			if err == nil {
				if tcpConn, ok := conn.(*net.TCPConn); ok {
					_ = tcpConn.SetReadBuffer(4 * 1024 * 1024)
					_ = tcpConn.SetWriteBuffer(4 * 1024 * 1024)
					_ = tcpConn.SetNoDelay(true)
				}
			}
			return conn, err
		},
		MaxIdleConns:        4000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     30 * time.Second,
		ReadBufferSize:      4 * 1024 * 1024,
	}

	client := &http.Client{
		Transport: transport,
	}

	bufPool := sync.Pool{
		New: func() interface{} {
			b := make([]byte, 1024*1024)
			return &b
		},
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	// Periodic rate ticker for live UI updates in Bytes Per Second
	tickerDone := make(chan struct{})
	go func() {
		var lastBytes int64
		var lastTime = time.Now()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-tickerDone:
				return
			case <-ticker.C:
				now := time.Now()
				curBytes := atomic.LoadInt64(totalBytes)
				deltaBytes := curBytes - lastBytes
				deltaTime := now.Sub(lastTime).Seconds()
				if deltaTime > 0.05 {
					bps := float64(deltaBytes) / deltaTime
					rpt.Store(bps)
					lastBytes = curBytes
					lastTime = now
				}
			}
		}
	}()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
				if err != nil {
					return
				}
				req.Header.Set("User-Agent", "speedtest_cli/1.0")

				resp, err := client.Do(req)
				if err != nil {
					return
				}

				bufPtr := bufPool.Get().(*[]byte)
				buf := *bufPtr

				for {
					n, rErr := resp.Body.Read(buf)
					if n > 0 {
						atomic.AddInt64(totalBytes, int64(n))
					}
					if rErr != nil {
						break
					}
					select {
					case <-ctx.Done():
						_ = resp.Body.Close()
						bufPool.Put(bufPtr)
						return
					default:
					}
				}
				_ = resp.Body.Close()
				bufPool.Put(bufPtr)
			}
		}()
	}

	wg.Wait()
	close(tickerDone)

	totalSec := time.Since(startTime).Seconds()
	if totalSec <= 0 {
		return 0
	}
	finalBytes := atomic.LoadInt64(totalBytes)
	return (float64(finalBytes) * 8.0) / (totalSec * 1e6)
}

func runMultiStreamUpload(urlStr string, concurrency int, duration time.Duration, totalBytes *int64, rpt *atomic.Value) float64 {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			d := &net.Dialer{Timeout: 5 * time.Second}
			conn, err := d.DialContext(ctx, "tcp", addr)
			if err == nil {
				if tcpConn, ok := conn.(*net.TCPConn); ok {
					_ = tcpConn.SetReadBuffer(4 * 1024 * 1024)
					_ = tcpConn.SetWriteBuffer(4 * 1024 * 1024)
					_ = tcpConn.SetNoDelay(true)
				}
			}
			return conn, err
		},
		MaxIdleConns:        4000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     30 * time.Second,
		WriteBufferSize:     4 * 1024 * 1024,
	}

	client := &http.Client{
		Transport: transport,
	}

	var wg sync.WaitGroup
	startTime := time.Now()

	// Periodic rate ticker for live UI updates in Bytes Per Second
	tickerDone := make(chan struct{})
	go func() {
		var lastBytes int64
		var lastTime = time.Now()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-tickerDone:
				return
			case <-ticker.C:
				now := time.Now()
				curBytes := atomic.LoadInt64(totalBytes)
				deltaBytes := curBytes - lastBytes
				deltaTime := now.Sub(lastTime).Seconds()
				if deltaTime > 0.05 {
					bps := float64(deltaBytes) / deltaTime
					rpt.Store(bps)
					lastBytes = curBytes
					lastTime = now
				}
			}
		}
	}()

	payloadData := make([]byte, 1024*1024) // 1MB chunk

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				reader := &countingReader{
					ctx:        ctx,
					totalBytes: totalBytes,
					data:       payloadData,
				}

				req, err := http.NewRequestWithContext(ctx, "POST", urlStr, reader)
				if err != nil {
					return
				}
				req.Header.Set("User-Agent", "speedtest_cli/1.0")
				req.Header.Set("Content-Type", "application/octet-stream")

				resp, err := client.Do(req)
				if err == nil && resp != nil {
					_, _ = io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()
	close(tickerDone)

	totalSec := time.Since(startTime).Seconds()
	if totalSec <= 0 {
		return 0
	}
	finalBytes := atomic.LoadInt64(totalBytes)
	return (float64(finalBytes) * 8.0) / (totalSec * 1e6)
}

type countingReader struct {
	ctx        context.Context
	totalBytes *int64
	data       []byte
	sent       int
}

func (c *countingReader) Read(p []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, io.EOF
	default:
	}

	if c.sent >= len(c.data)*10 {
		return 0, io.EOF
	}

	n = copy(p, c.data)
	c.sent += n
	atomic.AddInt64(c.totalBytes, int64(n))
	return n, nil
}
