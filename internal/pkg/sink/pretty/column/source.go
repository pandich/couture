package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/source"
	"crypto/sha256"
	"encoding/hex"
	"github.com/i582/cfmt/cmd/cfmt"
	"strings"
)

type sourceColumn struct {
	baseColumn
}

func newSourceColumn(cfg config.Config) column {
	layout := cfg.Layout.Source
	return sourceColumn{baseColumn{
		columnName: "source",
		widthMode:  weighted,
		colLayout:  layout,
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

// Render ...
func (col sourceColumn) Render(_ config.Config, event model.SinkEvent) string {
	return cfmt.Sprintf(formatStyleOfWidth(sourceID(event.SourceURL), col.layout().Width), event.SourceURL.ShortForm())
}

func sourceID(sourceURL model.SourceURL) string {
	hasher := sha256.New()
	hasher.Write([]byte(sourceURL.String()))
	return strings.ReplaceAll(hex.EncodeToString(hasher.Sum(nil)), "-", "")
}
