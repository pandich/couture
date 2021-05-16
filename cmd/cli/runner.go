package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/pretty"
	"github.com/spf13/cobra"
)

func runner(cmd *cobra.Command, args []string) error {
	persistent := cmd.PersistentFlags()
	verbosityOption, err := verbosityOption(persistent)
	if err != nil {
		return err
	}
	wrapOption, err := wrapOption(persistent)
	if err != nil {
		return err
	}
	filterOption, err := filterOption(persistent)
	if err != nil {
		return err
	}
	levelOption, err := levelOption(persistent)
	if err != nil {
		return err
	}
	sinceOption, err := sinceOption(persistent)
	if err != nil {
		return err
	}
	rateLimitOption, err := rateLimitOption(persistent)
	if err != nil {
		return err
	}

	var options = []interface{}{
		verbosityOption,
		filterOption,
		levelOption,
		wrapOption,
		rateLimitOption,
		pretty.New(),
	}
	if sinceOption != nil {
		options = append(options, sinceOption)
	}
	sourcesOptions, err := sourceOptions(args)
	if err != nil {
		return err
	}
	options = append(options, sourcesOptions...)

	mgr, err := manager.New(options...)
	if err != nil {
		return err
	}

	return (*mgr).Start()
}
