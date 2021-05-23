package column

import (
	"couture/internal/pkg/sink/pretty/theme"
)

type registry map[string]column

// ByName ...
var ByName = registry{}

func init() {
	for _, col := range columns {
		ByName[col.name()] = col
	}
}

// Names ...
func Names() []string {
	var columnNames []string
	for _, col := range columns {
		columnNames = append(columnNames, col.name())
	}
	return columnNames
}

// Init ...
func (registry registry) Init(theme theme.Theme) {
	for _, col := range registry {
		col.RegisterStyles(theme)
	}
}
