package column

import (
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/sink/layout"
)

type (
	column interface {
		render(event event.SinkEvent) string
		name() mapping.Column
		layout() layout.ColumnLayout
	}
	baseColumn struct {
		columnName mapping.Column
		colLayout  layout.ColumnLayout
	}
)

// GetName ...
func (col baseColumn) name() mapping.Column {
	return col.columnName
}

func (col baseColumn) layout() layout.ColumnLayout {
	return col.colLayout
}

func (col baseColumn) format() string {
	return col.layout().Format(string(col.columnName))
}

func (col baseColumn) formatWithSuffix(suffix string) string {
	return col.layout().Format(string(col.columnName) + suffix)
}
