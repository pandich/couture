package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
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

func newSourceColumn(layout layout.ColumnLayout) column {
	return sourceColumn{baseColumn{columnName: "source", widthMode: weighted, colLayout: layout}}
}

// RegisterSource ...
func RegisterSource(style theme.Style, src source.Source) {
	cfmt.RegisterStyle(sourceID(src.URL()), func(s string) string {
		return cfmt.Sprintf("{{%s }}::"+style.Fg+"|bg"+style.Bg+"{{ %s }}::"+style.Fg+"|bg"+style.Bg, string(src.Sigil()), s)
	})
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
