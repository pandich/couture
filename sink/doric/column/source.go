package column

import (
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"github.com/gagglepanda/couture/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

const sourcePseudoColumn = "source"

type sourceColumn struct {
	baseColumn
}

func newSourceColumn(layout layout.ColumnLayout) column {
	return sourceColumn{baseColumn{columnName: sourcePseudoColumn, colLayout: layout}}
}

func (col sourceColumn) render(event model.SinkEvent) string {
	return cfmt.Sprintf(col.formatWithSuffix(event.SourceURL.HashString()), event.SourceURL.ShortForm())
}

// RegisterSourceStyle ...
func RegisterSourceStyle(
	style color.FgBgTuple,
	layout layout.ColumnLayout,
	src source.Source,
) {
	layout.Sigil = string(src.Sigil())
	registerStyle(sourcePseudoColumn+src.URL().HashString(), style, layout)
}
