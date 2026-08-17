package history

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Dolyyyy/speedtest_cli/pkg/model"
)

func getHistoryFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".speedtest_history.json")
}

// LoadHistory reads history items from local disk
func LoadHistory() ([]model.HistoryItem, error) {
	filePath := getHistoryFilePath()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []model.HistoryItem{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var items []model.HistoryItem
	if err := json.Unmarshal(data, &items); err != nil {
		return []model.HistoryItem{}, nil
	}

	return items, nil
}

// SaveResult appends a benchmark result to local history
func SaveResult(res *model.TestResult) error {
	items, _ := LoadHistory()

	newItem := model.HistoryItem{
		Timestamp: res.Timestamp.Format("2006-01-02 15:04:05"),
		Server:    res.Server.Sponsor + " (" + res.Server.Name + ")",
		PingMs:    res.PingMs,
		JitterMs:  res.JitterMs,
		Download:  res.Download,
		Upload:    res.Upload,
	}

	items = append([]model.HistoryItem{newItem}, items...)

	// Retain last 50 entries max
	if len(items) > 50 {
		items = items[:50]
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getHistoryFilePath(), data, 0644)
}

// ClearHistory deletes local history file
func ClearHistory() error {
	filePath := getHistoryFilePath()
	if _, err := os.Stat(filePath); err == nil {
		return os.Remove(filePath)
	}
	return nil
}
