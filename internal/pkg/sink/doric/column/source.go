package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/sink"
	"github.com/pandich/couture/internal/pkg/sink/layout"
	"github.com/pandich/couture/internal/pkg/source"
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
	style sink.Style,
	layout layout.ColumnLayout,
	src source.Source,
) {
	layout.Sigil = string(src.Sigil())
	registerStyle(sourcePseudoColumn+src.URL().HashString(), style, layout)
}
