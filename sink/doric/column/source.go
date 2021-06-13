package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/source"
	"github.com/pandich/couture/theme/color"
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
	style color.HexPair,
	layout layout.ColumnLayout,
	src source.Source,
) {
	layout.Sigil = string(src.Sigil())
	registerStyle(sourcePseudoColumn+src.URL().HashString(), style, layout)
}
