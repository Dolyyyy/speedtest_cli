package main

import (
	"fmt"
	"os"

	"github.com/Dolyyyy/speedtest_cli/pkg/config"
	"github.com/Dolyyyy/speedtest_cli/pkg/engine"
	"github.com/Dolyyyy/speedtest_cli/pkg/history"
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

	if cfg.ClearHistory {
		if err := history.ClearHistory(); err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing history: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✔ Speedtest history cleared successfully.")
		return
	}

	if cfg.ShowHistory {
		items, err := history.LoadHistory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading history: %v\n", err)
			os.Exit(1)
		}
		ui.PrintHistoryTable(items, cfg.UseBytes)
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

	// Save result to local history
	_ = history.SaveResult(res)

	_ = printer.PrintResult(res, cfg)
}
