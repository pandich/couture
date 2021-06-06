package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"crypto/sha256"
	"encoding/hex"
	"github.com/i582/cfmt/cmd/cfmt"
	"strings"
)

type sourceColumn struct {
	baseColumn
}

func newSourceColumn() column {
	const weight = 40
	return sourceColumn{baseColumn{
		columnName:  "source",
		widthMode:   weighted,
		widthWeight: weight,
	}}
}

// Init ...
func (col sourceColumn) Init(_ theme.Theme) {}

// RegisterSource ...
func RegisterSource(th theme.Theme, consistentColors bool, src source.Source) string {
	fgColor, bgColor := th.SourceColor(consistentColors, src)
	sigilColor := fgColor
	cfmt.RegisterStyle(sourceID(src.URL()), func(s string) string {
		return cfmt.Sprintf("{{%s }}::"+sigilColor+"|bg"+bgColor+"{{ %s }}::"+fgColor+"|bg"+bgColor, string(src.Sigil()), s)
	})
	return bgColor
}

// RenderFormat ...
func (col sourceColumn) RenderFormat(width uint, event model.SinkEvent) string {
	return formatStyleOfWidth(sourceID(event.SourceURL), width)
}

// RenderValue ...
func (col sourceColumn) RenderValue(_ config.Config, event model.SinkEvent) []interface{} {
	return []interface{}{event.SourceURL.ShortForm()}
}

func sourceID(sourceURL model.SourceURL) string {
	hasher := sha256.New()
	hasher.Write([]byte(sourceURL.String()))
	return strings.ReplaceAll(hex.EncodeToString(hasher.Sum(nil)), "-", "")
}
