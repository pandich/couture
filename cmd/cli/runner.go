package cli

import (
	"couture/internal/pkg/manager"
	"github.com/spf13/cobra"
)

func runner(cmd *cobra.Command, args []string) error {
	flags := cmd.PersistentFlags()
	options, err := getOptions(flags)
	if err != nil {
		return err
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
