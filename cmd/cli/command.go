package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/pretty"
	"github.com/pkg/errors"
	"github.com/riywo/loginshell"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
)

func handleDocCommand(cmd *cobra.Command) error {
	format := os.Args[2]
	switch format {
	case "man":
		return doc.GenMan(cmd, &doc.GenManHeader{
			Title:   "Couture",
			Section: "5",
			Source:  "",
			Manual:  "",
		}, os.Stdout)
	case "md", "markdown":
		return doc.GenMarkdown(cmd, os.Stdout)
	default:
		return errors.Errorf("invalid documentation format: %s", format)
	}
}

func handleCompleteCommand(cmd *cobra.Command) error {
	const shellNameArgIndex = 2
	var shellName string
	if len(os.Args) > shellNameArgIndex {
		shellName = os.Args[shellNameArgIndex]
	} else {
		shellBinary, err := loginshell.Shell()
		if err != nil {
			return err
		}
		shellName = path.Base(shellBinary)
	}
	switch shellName {
	case "bash":
		return cmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.GenFishCompletion(os.Stdout, true)
	case "powershell", "powershell.exe":
		return cmd.GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return errors.Errorf("invalid shell: %s", shellName)
	}
}

func handleLogCommand(cmd *cobra.Command) error {
	return cmd.Execute()
}

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

	var options = []interface{}{
		verbosityOption,
		filterOption,
		levelOption,
		wrapOption,
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
