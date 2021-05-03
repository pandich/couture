package cli

import (
	"couture/internal/pkg/manager"
)

// filterOptions returns sources manager.option values defined in filterCLI.
func filterOptions() []interface{} {
	return []interface{}{
		manager.FilterOption(cli.Log.IncludeFilters, cli.Log.ExcludeFilters),
		manager.LogLevelOption(cli.Log.Level.Normalize()),
	}
}
