package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/tty"
)

// Config ...
type Config struct {
	Wrap        bool
	Width       uint
	MultiLine   bool
	TimeFormat  string
	Theme       theme.Theme
	ClearScreen bool
	ShowSigils  bool
	Columns     []string
}

// WrapWidth ...
func (cfg Config) WrapWidth() int {
	if cfg.Width > 0 {
		return int(cfg.Width)
	}
	if cfg.Wrap {
		return tty.TerminalWidth()
	}
	return 0
}
