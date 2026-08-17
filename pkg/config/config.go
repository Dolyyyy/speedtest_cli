package config

import (
	"flag"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

// ParseFlags parses command-line flags and constructs model.Config
func ParseFlags() *model.Config {
	cfg := &model.Config{}

	flag.BoolVar(&cfg.ShowHelp, "h", false, "Display help guide")
	flag.BoolVar(&cfg.ShowHelp, "help", false, "Display help guide")

	flag.BoolVar(&cfg.ShowList, "l", false, "List nearest available speedtest servers")
	flag.BoolVar(&cfg.ShowList, "list", false, "List nearest available speedtest servers")

	flag.StringVar(&cfg.ServerID, "s", "", "Specify a target server ID")
	flag.StringVar(&cfg.ServerID, "server", "", "Specify a target server ID")

	flag.StringVar(&cfg.CustomHost, "custom", "", "Specify a custom speedtest server host (e.g. host:port)")

	flag.IntVar(&cfg.Threads, "t", 16, "Number of parallel TCP connections for high-speed link saturation (1-64)")
	flag.IntVar(&cfg.Threads, "threads", 16, "Number of parallel TCP connections for high-speed link saturation (1-64)")

	flag.BoolVar(&cfg.UseBytes, "bytes", false, "Display speed in MB/s instead of Mbps")

	flag.BoolVar(&cfg.UseJSON, "json", false, "Output results in JSON format")

	flag.BoolVar(&cfg.UseSimple, "simple", false, "Output minimalist 3-line format")

	flag.BoolVar(&cfg.ShowVersion, "v", false, "Show CLI version")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Show CLI version")

	flag.Parse()

	if cfg.Threads < 1 {
		cfg.Threads = 1
	} else if cfg.Threads > 64 {
		cfg.Threads = 64
	}

	return cfg
}
