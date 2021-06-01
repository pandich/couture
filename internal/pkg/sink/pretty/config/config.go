package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/olekukonko/ts"
	"os"
)

// Config ...
type Config struct {
	AutoResize       bool        `json:"auto_resize"`
	Banner           bool        `json:"banner"`
	Color            bool        `json:"color"`
	Columns          []string    `json:"columns"`
	ConsistentColors bool        `json:"consistent_colors"`
	ExpandJSON       bool        `json:"expand_json"`
	Highlight        bool        `json:"highlight"`
	Multiline        bool        `json:"multiline"`
	Theme            theme.Theme `json:"theme"`
	TimeFormat       string      `json:"time_format"`
	TTY              bool        `json:"-"`
	Width            uint        `json:"width"`
	Wrap             bool        `json:"wrap"`
	Out              *os.File    `json:"-"`
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
