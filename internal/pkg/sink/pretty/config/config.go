package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/olekukonko/ts"
)

// Config ...
type Config struct {
	AutoResize       bool
	Columns          []string
	ConsistentColors bool
	Highlight        bool
	Multiline        bool
	Theme            theme.Theme
	TimeFormat       string
	Width            uint
	Wrap             bool
}

// EffectiveTerminalWidth ...
func (cfg Config) EffectiveTerminalWidth() uint {
	if cfg.Width > 0 {
		return cfg.Width
	}
	if cfg.Wrap {
		return uint(terminalWidth())
	}
	return 0
}

func terminalWidth() int {
	var terminalWidth = 0
	if size, err := ts.GetSize(); err == nil {
		terminalWidth = size.Col()
	}
	return terminalWidth
}
