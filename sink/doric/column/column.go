package column

import (
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink/layout"
)

type (
	column interface {
		render(event model.SinkEvent) string
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
