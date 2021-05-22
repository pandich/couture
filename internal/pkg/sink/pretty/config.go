package pretty

import (
	"couture/internal/pkg/tty"
)

// noWrap ...
const noWrap = 0

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

func (c Config) wrapWidth() int {
	if c.Width > noWrap {
		return int(c.Width)
	}
	if c.Wrap {
		return tty.TerminalWidth()
	}
	return noWrap
}
