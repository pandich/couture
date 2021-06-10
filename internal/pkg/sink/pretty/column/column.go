package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
)

type (
	column interface {
		render(event model.SinkEvent) string
		name() string
		layout() layout.ColumnLayout
	}
	baseColumn struct {
		columnName string
		colLayout  layout.ColumnLayout
	}
)

// GetName ...
func (col baseColumn) name() string {
	return col.columnName
}

func (col baseColumn) layout() layout.ColumnLayout {
	return col.colLayout
}

func (col baseColumn) format() string {
	return col.layout().Format(col.columnName)
}

func (col baseColumn) formatWithSuffix(suffix string) string {
	return col.layout().Format(col.columnName + suffix)
}
