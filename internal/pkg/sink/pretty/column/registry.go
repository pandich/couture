package column

import (
	"couture/internal/pkg/sink/pretty/theme"
)

func init() {
	// build the registry
	for _, col := range columns {
		ByName[col.name()] = col
	}
}

// Names all available column names.
func Names() []string {
	var columnNames []string
	for _, col := range columns {
		columnNames = append(columnNames, col.name())
	}
	return columnNames
}

// ByName ...
var ByName = registry{}

type registry map[string]column

// Init initializes a theme with the registry.
func (registry registry) Init(theme theme.Theme) {
	for _, col := range registry {
		col.RegisterStyles(theme)
	}
}
