package column

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

type sourceColumn struct {
	baseColumn
}

func newSourceColumn() sourceColumn {
	const weight = 40
	return sourceColumn{baseColumn{
		columnName:  "source",
		widthMode:   weighted,
		widthWeight: weight,
	}}
}

// RegisterStyles ...
func (col sourceColumn) RegisterStyles(_ theme.Theme) {}

// RegisterSource ...
func RegisterSource(th theme.Theme, consistentColors bool, src source.Source) {
	bgColor := th.SourceColor(consistentColors, src)
	fgColor := contrast(bgColor)
	sigilColor := fgColor
	cfmt.RegisterStyle(src.ID(), func(s string) string {
		return cfmt.Sprintf("{{%s}}::"+sigilColor+"|bg"+bgColor+"{{ %s }}::"+fgColor+"|bg"+bgColor, string(src.Sigil()), s)
	})
}

// Format ...
func (col sourceColumn) Format(width uint, event sink.Event) string {
	return formatStyleOfWidth(event.Source.ID(), width)
}

// Render ...
func (col sourceColumn) Render(_ config.Config, event sink.Event) []interface{} {
	return []interface{}{event.Source.URL().ShortForm()}
}
