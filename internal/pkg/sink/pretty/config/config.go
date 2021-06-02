package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/ts"
	"os"
)

// Config ...
type Config struct {
	AutoResize       bool        `json:"auto_resize"`
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

// EffectiveIsTTY ...
func (cfg Config) EffectiveIsTTY() bool {
	return isatty.IsTerminal(cfg.Out.Fd()) || cfg.TTY
}
