package pretty

import (
	"couture/internal/pkg/tty"
	"github.com/muesli/gamut"
)

// noWrap ...
const noWrap = 0

const (
	errorColor = "#ff0000"
	traceColor = "#868686"
	debugColor = "#f6f6f6"
	infoColor  = "#00ff00"
	warnColor  = "#ffff00"
)

// Config ...
type Config struct {
	Wrap       bool
	Width      uint
	MultiLine  bool
	TimeFormat string
	Theme      Theme
	Columns    []ColumnName
}

func (c Config) effectiveColumns() []ColumnName {
	var columnOrder = defaultColumnOrder
	if len(c.Columns) > 0 {
		columnOrder = c.Columns
	}
	columnOrder = append([]ColumnName{sourceColumn}, columnOrder...)
	return columnOrder
}

func (c Config) palette() chan string {
	return tty.NewColorCycle(c.Theme.SourceColors)
}

func (c Config) terminalWidth() int {
	if c.Width > noWrap {
		return int(c.Width)
	}
	if c.Wrap {
		return tty.TerminalWidth()
	}
	return noWrap
}

// Theme ...
type Theme struct {
	BaseColor    string
	SourceColors gamut.ColorGenerator
}

func (t Theme) tinted(hex string) string {
	return tty.Tinted(t.BaseColor, hex)
}
