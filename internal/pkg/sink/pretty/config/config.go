package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/tty"
)

// noWrap ...
const noWrap = 0

// Config ...
type Config struct {
	Wrap        bool
	Width       uint
	MultiLine   bool
	TimeFormat  string
	Theme       theme.Theme
	ClearScreen bool
	Columns     []string
}

// WrapWidth ...
func (c Config) WrapWidth() int {
	if c.Width > noWrap {
		return int(c.Width)
	}
	if c.Wrap {
		return tty.TerminalWidth()
	}
	return noWrap
}
