package cli

import (
	"couture/internal/pkg/manager"
)

// filterOptions returns sources manager.option values defined in filterCLI.
func filterOptions() []interface{} {
	return []interface{}{
		manager.FilterOption(cli.IncludeFilters, cli.ExcludeFilters),
		manager.LogLevelOption(cli.Level.Normalize()),
	}
}
