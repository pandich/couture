package column

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
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
func RegisterSource(theme theme.Theme, consistentColors bool, src source.Source) {
	bgColor := theme.SourceColor(consistentColors, src)
	fgColor := tty.Contrast(bgColor)
	// TODO sigil should stand out
	cfmt.RegisterStyle(src.ID(), func(s string) string {
		return cfmt.Sprintf("{{"+string(src.Sigil())+" %s }}::"+fgColor+"|bg"+bgColor, s)
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
