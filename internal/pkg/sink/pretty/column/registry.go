package column

import (
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
	"couture/internal/pkg/tty"
	"github.com/i582/cfmt/cmd/cfmt"
)

type registry map[string]column

// ByName ...
var ByName = registry{}

func init() {
	for _, col := range columns {
		ByName[col.name] = col
	}
}

// Names ...
func Names() []string {
	var columnNames []string
	for _, col := range columns {
		columnNames = append(columnNames, col.name)
	}
	return columnNames
}

// Init ...
func (registry registry) Init(theme theme.Theme) {
	for _, col := range registry {
		col.register(theme)
	}
}

// RegisterSourceStyle ...
func RegisterSourceStyle(src source.Source, styleColor string) {
	bgColor := styleColor
	fgColor := tty.Contrast(bgColor)
	cfmt.RegisterStyle(src.ID(), func(s string) string { return cfmt.Sprintf("{{%s}}::"+fgColor+"|bg"+bgColor, s) })
}
