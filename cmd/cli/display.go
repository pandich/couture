package cli

import (
	"couture/internal/pkg/manager"
)

// displayOptions returns sources manager.option values defined in displayCLI.
func displayOptions() []interface{} {
	var verbosity = cli.Verbosity
	if cli.Quiet && cli.Verbosity > 0 {
		verbosity--
	}
	return []interface{}{
		manager.VerboseDisplayOption(verbosity),
	}
}
