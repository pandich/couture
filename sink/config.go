package sink

import (
	"github.com/gagglepanda/couture/sink/layout"
	theme2 "github.com/gagglepanda/couture/sink/theme"
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/ts"
	"os"
	"time"
)

var (
	// Enabled specifies that a feature is enabled.
	Enabled = true
	// Disabled specifies that a feature is disabled.
	Disabled = false

	// DefaultTimeFormat is the default time format used by the sink.
	DefaultTimeFormat = time.Stamp
)

// DefaultConfig returns default configuration for the sink.
func DefaultConfig() Config {
	return Config{
		AutoResize:       &Enabled,
		Color:            &Enabled,
		ConsistentColors: &Enabled,
		Expand:           &Disabled,
		Highlight:        &Disabled,
		MultiLine:        &Disabled,
		Wrap:             &Disabled,
		Layout:           &layout.Default,
		Out:              os.Stdout,
		Theme:            nil,
		TimeFormat:       &DefaultTimeFormat,
	}
}

// Config ...
type Config struct {
	AutoResize       *bool          `yaml:"auto_resize,omitempty"`
	Color            *bool          `yaml:"color,omitempty"`
	Columns          []string       `yaml:"columns,omitempty"`
	ConsistentColors *bool          `yaml:"consistent_colors,omitempty"`
	Expand           *bool          `yaml:"expand,omitempty"`
	Highlight        *bool          `yaml:"highlight,omitempty"`
	MultiLine        *bool          `yaml:"multi_line,omitempty"`
	LevelMeter       *bool          `yaml:"level_meter,omitempty"`
	Theme            *theme2.Theme  `yaml:"theme,omitempty"`
	Layout           *layout.Layout `yaml:"-"`
	TimeFormat       *string        `yaml:"time_format,omitempty"`
	TTY              bool           `yaml:"-"`
	Width            *uint          `yaml:"width,omitempty"`
	Wrap             *bool          `yaml:"wrap,omitempty"`
	Out              *os.File       `yaml:"-"`
}

// EffectiveTerminalWidth ...
func (config *Config) EffectiveTerminalWidth() uint {
	if config.Width != nil && *config.Width > 0 {
		return *config.Width
	}
	if config.Wrap != nil && *config.Wrap {
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
func (config *Config) EffectiveIsTTY() bool {
	return isatty.IsTerminal(config.Out.Fd()) || config.TTY
}

// PopulateMissing ...
func (config *Config) PopulateMissing(other Config) *Config {
	if config.AutoResize == nil {
		config.AutoResize = other.AutoResize
	}
	if config.Color == nil {
		config.Color = other.Color
	}
	if config.Columns == nil {
		config.Columns = other.Columns
	}
	if config.ConsistentColors == nil {
		config.ConsistentColors = other.ConsistentColors
	}
	if config.Expand == nil {
		config.Expand = other.Expand
	}
	if config.Highlight == nil {
		config.Highlight = other.Highlight
	}
	if config.Layout == nil {
		config.Layout = other.Layout
	}
	if config.LevelMeter == nil {
		config.LevelMeter = other.LevelMeter
	}
	if config.MultiLine == nil {
		config.MultiLine = other.MultiLine
	}
	if config.Theme == nil {
		config.Theme = other.Theme
	}
	if config.TimeFormat == nil {
		config.TimeFormat = other.TimeFormat
	}
	if config.Width == nil {
		config.Width = other.Width
	}
	if config.Wrap == nil {
		config.Wrap = other.Wrap
	}
	if config.Out == nil {
		config.Out = other.Out
	}
	return config
}
