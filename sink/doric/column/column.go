package column

import (
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/schema"
	"github.com/gagglepanda/couture/sink/layout"
)

type (
	column interface {
		render(event event.SinkEvent) string
		name() schema.Column
		layout() layout.ColumnLayout
	}
	baseColumn struct {
		columnName schema.Column
		colLayout  layout.ColumnLayout
	}
)

// GetName ...
func (col baseColumn) name() schema.Column {
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
