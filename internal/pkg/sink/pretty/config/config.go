package config

import (
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/tty"
)

// Config ...
type Config struct {
	AutoSize    bool
	ClearScreen bool
	Columns     []string
	MultiLine   bool
	ShowSigils  bool
	Theme       theme.Theme
	TimeFormat  string
	Width       uint
	Wrap        bool
}

// EffectiveTerminalWidth ...
func (cfg Config) EffectiveTerminalWidth() uint {
	if cfg.Width > 0 {
		return cfg.Width
	}
	if cfg.Wrap {
		return uint(tty.TerminalWidth())
	}
	return 0
}