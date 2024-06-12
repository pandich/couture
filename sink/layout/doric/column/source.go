package column

import (
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/source"
)

const sourcePseudoColumn = "source"

type sourceColumn struct {
	baseColumn
}

func newSourceColumn(layout layout.ColumnLayout) column {
	return sourceColumn{baseColumn{columnName: sourcePseudoColumn, colLayout: layout}}
}

func (col sourceColumn) render(event event.SinkEvent) string {
	path := event.SourceURL.Path
	if uint(len(event.SourceURL.Path)) > col.colLayout.Width {
		path = event.SourceURL.Path[len(event.SourceURL.Path)-int(col.colLayout.Width):]
	}

	return cfmt.Sprintf(col.formatWithSuffix(event.SourceURL.HashString()), path)
}

// RegisterSourceStyle ...
func RegisterSourceStyle(
	style color.FgBgTuple,
	layout layout.ColumnLayout,
	src source.Source,
) {
	layout.Sigil = string(src.Sigil())
	url := src.URL()
	registerStyle(sourcePseudoColumn+url.HashString(), style, layout)
}
