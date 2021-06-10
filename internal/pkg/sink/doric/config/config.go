package config

import (
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/ts"
	"os"
)

// Config ...
type Config struct {
	AutoResize       *bool           `yaml:"auto_resize,omitempty"`
	ShowSchema       *bool           `yaml:"show_schema,omitempty"`
	Color            *bool           `yaml:"color,omitempty"`
	Columns          []string        `yaml:"columns,omitempty"`
	ConsistentColors *bool           `yaml:"consistent_colors,omitempty"`
	Expand           *bool           `yaml:"expand,omitempty"`
	Highlight        *bool           `yaml:"highlight,omitempty"`
	Multiline        *bool           `yaml:"multiline,omitempty"`
	Theme            *theme2.Theme   `yaml:"theme,omitempty"`
	Layout           *layout2.Layout `yaml:"-"`
	TimeFormat       *string         `yaml:"time_format,omitempty"`
	TTY              bool            `yaml:"-"`
	Width            *uint           `yaml:"width,omitempty"`
	Wrap             *bool           `yaml:"wrap,omitempty"`
	Out              *os.File        `yaml:"-"`
}

// EffectiveTerminalWidth ...
func (cfg Config) EffectiveTerminalWidth() uint {
	if cfg.Width != nil && *cfg.Width > 0 {
		return *cfg.Width
	}
	if cfg.Wrap != nil && *cfg.Wrap {
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

// FillMissing ...
func (cfg *Config) FillMissing(other Config) {
	if cfg.AutoResize == nil {
		cfg.AutoResize = other.AutoResize
	}
	if cfg.ShowSchema == nil {
		cfg.ShowSchema = other.ShowSchema
	}
	if cfg.Color == nil {
		cfg.Color = other.Color
	}
	if cfg.Columns == nil {
		cfg.Columns = other.Columns
	}
	if cfg.ConsistentColors == nil {
		cfg.ConsistentColors = other.ConsistentColors
	}
	if cfg.Expand == nil {
		cfg.Expand = other.Expand
	}
	if cfg.Highlight == nil {
		cfg.Highlight = other.Highlight
	}
	if cfg.Layout == nil {
		cfg.Layout = other.Layout
	}
	if cfg.Multiline == nil {
		cfg.Multiline = other.Multiline
	}
	if cfg.Theme == nil {
		cfg.Theme = other.Theme
	}
	if cfg.TimeFormat == nil {
		cfg.TimeFormat = other.TimeFormat
	}
	if cfg.Width == nil {
		cfg.Width = other.Width
	}
	if cfg.Wrap == nil {
		cfg.Wrap = other.Wrap
	}
	if cfg.Out == nil {
		cfg.Out = other.Out
	}
}
