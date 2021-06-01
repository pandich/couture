package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/olekukonko/ts"
)

// Config ...
type Config struct {
	AutoResize       bool
	Banner           bool
	Columns          []string
	ConsistentColors bool
	ExpandJSON       bool
	Highlight        bool
	Multiline        bool
	Theme            theme.Theme
	TimeFormat       string
	TTY              bool
	Width            uint
	Wrap             bool
}

// EffectiveTerminalWidth ...
func (cfg Config) EffectiveTerminalWidth() uint {
	if cfg.Width > 0 {
		return cfg.Width
	}
	if cfg.Wrap {
		return uint(TerminalWidth())
	}
	return 0
}

// TerminalWidth ...
func TerminalWidth() int {
	var terminalWidth = 0
	if size, err := ts.GetSize(); err == nil {
		terminalWidth = size.Col()
	}
	return terminalWidth
}
