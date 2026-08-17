package ui

import (
	"github.com/fatih/color"
)

// UI Color palette helpers
var (
	ColorTitle   = color.New(color.FgHiCyan, color.Bold).SprintFunc()
	ColorHeader  = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	ColorSuccess = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	ColorWarning = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	ColorError   = color.New(color.FgHiRed, color.Bold).SprintFunc()
	ColorMuted   = color.New(color.FgHiBlack).SprintFunc()
	ColorVal     = color.New(color.FgWhite, color.Bold).SprintFunc()
	ColorAccent  = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
)
