package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/ts"
	"os"
)

// Config ...
type Config struct {
	AutoResize       bool        `yaml:"auto_resize,omitempty"`
	ShowSchema       bool        `yaml:"show_schema,omitempty"`
	Color            bool        `yaml:"color,omitempty"`
	Columns          []string    `yaml:"columns,omitempty"`
	ConsistentColors bool        `yaml:"consistent_colors,omitempty"`
	ExpandJSON       bool        `yaml:"expand_json,omitempty"`
	Highlight        bool        `yaml:"highlight,omitempty"`
	Multiline        bool        `yaml:"multiline,omitempty"`
	Theme            theme.Theme `yaml:"theme,omitempty"`
	TimeFormat       string      `yaml:"time_format,omitempty"`
	TTY              bool        `yaml:"-"`
	Width            uint        `yaml:"width,omitempty"`
	Wrap             bool        `yaml:"wrap,omitempty"`
	Out              *os.File    `yaml:"-"`
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
