package model

// Config holds runtime configuration options for the CLI
type Config struct {
	ShowHelp    bool
	ShowList    bool
	ServerID    string
	CustomHost  string
	Threads     int
	UseBytes    bool
	UseJSON     bool
	UseSimple   bool
	ShowVersion bool
}
