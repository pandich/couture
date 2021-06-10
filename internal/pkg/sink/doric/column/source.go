package column

import (
	"couture/internal/pkg/model"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

const sourcePseudoColumn = "source"

type sourceColumn struct {
	baseColumn
}

func newSourceColumn(layout layout2.ColumnLayout) column {
	return sourceColumn{baseColumn{columnName: sourcePseudoColumn, colLayout: layout}}
}

func (col sourceColumn) render(event model.SinkEvent) string {
	return cfmt.Sprintf(col.formatWithSuffix(event.SourceURL.HashString()), event.SourceURL.ShortForm())
}

// RegisterSourceStyle ...
func RegisterSourceStyle(
	style theme2.Style,
	layout layout2.ColumnLayout,
	src source.Source,
) {
	layout.Sigil = string(src.Sigil())
	registerStyle(sourcePseudoColumn+src.URL().HashString(), style, layout)
}
