package main

import (
	"fmt"
	"os"

	"github.com/Dolyyyy/speedtest_cli/pkg/config"
	"github.com/Dolyyyy/speedtest_cli/pkg/engine"
	"github.com/Dolyyyy/speedtest_cli/pkg/model"
	"github.com/Dolyyyy/speedtest_cli/pkg/printer"
	"github.com/Dolyyyy/speedtest_cli/pkg/ui"
)

func main() {
	cfg := config.ParseFlags()

	if cfg.ShowVersion {
		fmt.Printf("speedtest_cli v%s\n", model.Version)
		return
	}

	if cfg.ShowHelp {
		ui.PrintHelp()
		return
	}

	runner := engine.NewRunner(cfg)

	if cfg.ShowList {
		servers, err := engine.FetchServerItems(runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		ui.PrintServerList(servers)
		return
	}

	res, err := engine.Run(runner)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printer.PrintResult(res, cfg)
}
