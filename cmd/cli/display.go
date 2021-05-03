package cli

import (
	"couture/internal/pkg/manager"
)

// displayOptions returns sources manager.option values defined in displayCLI.
func displayOptions() []interface{} {
	var verbosity = cli.Log.Verbosity
	if cli.Log.Quiet && cli.Log.Verbosity > 0 {
		verbosity--
	}
	return []interface{}{
		manager.VerboseDisplayOption(verbosity),
	}
}
