package config

import (
	"flag"
	"os"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

// ParseFlags parses os.Args and constructs model.Config
func ParseFlags() *model.Config {
	return ParseFlagsArgs(os.Args[1:])
}

// ParseFlagsArgs parses explicit argument slice and constructs model.Config
func ParseFlagsArgs(args []string) *model.Config {
	cfg := &model.Config{}
	fs := flag.NewFlagSet("speedtest", flag.ContinueOnError)

	fs.BoolVar(&cfg.ShowHelp, "h", false, "Display help guide")
	fs.BoolVar(&cfg.ShowHelp, "help", false, "Display help guide")

	fs.BoolVar(&cfg.ShowList, "l", false, "List nearest available speedtest servers")
	fs.BoolVar(&cfg.ShowList, "list", false, "List nearest available speedtest servers")

	fs.BoolVar(&cfg.ShowHistory, "H", false, "Display local speedtest benchmark history")
	fs.BoolVar(&cfg.ShowHistory, "history", false, "Display local speedtest benchmark history")

	fs.BoolVar(&cfg.ClearHistory, "clear-history", false, "Clear all saved local speedtest history records")

	fs.StringVar(&cfg.ServerID, "s", "", "Specify a target server ID")
	fs.StringVar(&cfg.ServerID, "server", "", "Specify a target server ID")

	fs.StringVar(&cfg.CustomHost, "custom", "", "Specify a custom speedtest server host (e.g. host:port)")

	fs.IntVar(&cfg.Threads, "t", 16, "Number of parallel TCP connections for 10G/50G/100G link saturation (1-128)")
	fs.IntVar(&cfg.Threads, "threads", 16, "Number of parallel TCP connections for 10G/50G/100G link saturation (1-128)")

	fs.BoolVar(&cfg.UseBytes, "bytes", false, "Display speed in MB/s instead of Mbps")

	fs.BoolVar(&cfg.UseJSON, "json", false, "Output results in JSON format")

	fs.BoolVar(&cfg.UseSimple, "simple", false, "Output minimalist 3-line format")

	fs.BoolVar(&cfg.ShowVersion, "v", false, "Show CLI version")
	fs.BoolVar(&cfg.ShowVersion, "version", false, "Show CLI version")

	_ = fs.Parse(args)

	if cfg.Threads < 1 {
		cfg.Threads = 1
	} else if cfg.Threads > 128 {
		cfg.Threads = 128
	}

	return cfg
}
