package sink

import (
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/ts"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/theme"
	"os"
)

// Config ...
type Config struct {
	AutoResize       *bool          `yaml:"auto_resize,omitempty"`
	ShowSchema       *bool          `yaml:"show_schema,omitempty"`
	Color            *bool          `yaml:"color,omitempty"`
	Columns          []string       `yaml:"columns,omitempty"`
	ConsistentColors *bool          `yaml:"consistent_colors,omitempty"`
	Expand           *bool          `yaml:"expand,omitempty"`
	Highlight        *bool          `yaml:"highlight,omitempty"`
	MultiLine        *bool          `yaml:"multi_line,omitempty"`
	Theme            *theme.Theme   `yaml:"theme,omitempty"`
	Layout           *layout.Layout `yaml:"-"`
	TimeFormat       *string        `yaml:"time_format,omitempty"`
	TTY              bool           `yaml:"-"`
	Width            *uint          `yaml:"width,omitempty"`
	Wrap             *bool          `yaml:"wrap,omitempty"`
	Out              *os.File       `yaml:"-"`
}

// EffectiveTerminalWidth ...
func (config Config) EffectiveTerminalWidth() uint {
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
func (config Config) EffectiveIsTTY() bool {
	return isatty.IsTerminal(config.Out.Fd()) || config.TTY
}

// PopulateMissing ...
func (config *Config) PopulateMissing(other Config) *Config {
	if config.AutoResize == nil {
		config.AutoResize = other.AutoResize
	}
	if config.ShowSchema == nil {
		config.ShowSchema = other.ShowSchema
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
	if config.MultiLine == nil {
		config.MultiLine = other.MultiLine
	}
	if config.ShowSchema == nil {
		config.ShowSchema = other.ShowSchema
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
