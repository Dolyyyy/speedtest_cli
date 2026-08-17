package model

// Config holds runtime configuration options for speedtest execution
type Config struct {
	ShowHelp     bool
	ShowList     bool
	ShowHistory  bool
	ClearHistory bool
	ServerID     string
	CustomHost   string
	Threads      int
	UseBytes     bool
	UseJSON      bool
	UseSimple    bool
	ShowVersion  bool
	Mode100G     bool
}

// IsQuiet returns true if formatted UI output should be suppressed (JSON or Simple mode)
func (c *Config) IsQuiet() bool {
	return c.UseJSON || c.UseSimple
}
