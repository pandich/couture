package column

import (
	"couture/internal/pkg/model"
	layout2 "couture/internal/pkg/sink/layout"
)

type (
	column interface {
		render(event model.SinkEvent) string
		name() string
		layout() layout2.ColumnLayout
	}
	baseColumn struct {
		columnName string
		colLayout  layout2.ColumnLayout
	}
)

// GetName ...
func (col baseColumn) name() string {
	return col.columnName
}

func (col baseColumn) layout() layout2.ColumnLayout {
	return col.colLayout
}

func (col baseColumn) format() string {
	return col.layout().Format(col.columnName)
}

func (col baseColumn) formatWithSuffix(suffix string) string {
	return col.layout().Format(col.columnName + suffix)
}
