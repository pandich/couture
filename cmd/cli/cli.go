package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/internal/pkg/source/ssh"
	"couture/internal/pkg/source/tail"
	"github.com/pkg/errors"
	"github.com/riywo/loginshell"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

const applicationName = "couture"

var (
	errConfigNotFound = &viper.ConfigFileNotFoundError{}
)

const (
	generateShellCompletionCommand = "complete"
	generateDocumentationCommand   = "doc"
)

var defaultFilters []string

const (
	defaultSince     = "5m"
	defaultSink      = "pretty"
	defaultPaginator = ""
	defaultLevel     = "info"
)

var rootCmd = &cobra.Command{
	Use:   applicationName + " [flags] source_url ...",
	Short: "Tails one or more event sources.\n",
	Long: "Description:\n\nTails one or more event sources.\n" +
		"When providing a CloudFormation stack, resources are recursively analyzed until all loggable entities are found. " +
		"This includes the stack events of the stack itself, as well as any log groups " +
		"its entities contain.\n",
	Example: strings.Join([]string{
		"\n  " + applicationName + " " + strings.Join(
			sourceMetadata.ExampleURLs(),
			"\n  "+applicationName+" "),
	}, ""),
	Args: cobra.MinimumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlag(paginatorEnvKey, cmd.PersistentFlags().Lookup(paginatorFlag)); err != nil {
			return err
		}

		viper.SetConfigName("." + applicationName)
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil && !errors.As(err, &errConfigNotFound) {
			return errors.Errorf("fatal error config file: %s\n", err)
		}
		viper.AutomaticEnv()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.PersistentFlags()
		resources, err := allOptions(flags, args)
		if err != nil {
			return err
		}

		mgr, err := manager.New(resources...)
		if err != nil {
			return err
		}

		return (*mgr).Start()
	},
}

// sourceMetadata is a list of sourceMetadata sourceMetadata.
var sourceMetadata = source.MetadataGroup{
	fake.Metadata(),
	cloudwatch.Metadata(),
	cloudformation.Metadata(),
	elasticsearch.Metadata(),
	tail.Metadata(),
	ssh.Metadata(),
}

// Execute ...
func Execute() error {
	flags := rootCmd.PersistentFlags()
	flags.StringP(outputFormatFlag, "o", defaultSink, "The output format. [pretty | json]")
	flags.CountP(verboseFlag, "v", "Display additional diagnostic data.")
	flags.StringP(paginatorFlag, "p", defaultPaginator, "Paginate output.")
	flags.BoolP(noPaginatorFlag, "P", false, "Do not paginate output.")
	flags.UintP(wrapFlag, "w", manager.NoWrap, "Display no diagnostic data.")
	flags.StringP(levelFlag, "l", defaultLevel, "Minimum log level to display (trace, debug, info warn, error.")
	flags.StringP(sinceFlag, "s", defaultSince, "How far back in time to search for events.")
	flags.StringSliceP(includeFilterFlag, "i", defaultFilters, "Include filter regular expressions. Performed before excludes.")
	flags.StringSliceP(excludeFilterFlag, "e", defaultFilters, "Exclude filter regular expressions. Performed after includes.")

	if (len(os.Args) == 2 || len(os.Args) == 3) && os.Args[1] == generateShellCompletionCommand {
		return handleCompleteCommand(rootCmd)
	}
	if len(os.Args) == 3 && os.Args[1] == generateDocumentationCommand {
		return handleDocCommand(rootCmd)
	}
	return rootCmd.Execute()
}

func handleDocCommand(cmd *cobra.Command) error {
	docFormat := os.Args[2]
	switch docFormat {
	case "man":
		return doc.GenMan(cmd, &doc.GenManHeader{Title: applicationName, Section: "1"}, os.Stdout)
	case "markdown":
		return doc.GenMarkdown(cmd, os.Stdout)
	case "yaml":
		return doc.GenYaml(cmd, os.Stdout)
	default:
		return errors.Errorf("invalid documentation format: %s - must be (man | markdown | yaml)\n", docFormat)
	}
}

func handleCompleteCommand(cmd *cobra.Command) error {
	// FIXME file is generated, but completions don't work
	shellName, err := shellName()
	if err != nil {
		return err
	}

	writer := os.Stdout
	switch shellName {
	case "bash":
		return cmd.GenBashCompletion(writer)
	case "zsh":
		return cmd.GenZshCompletion(writer)
	case "fish":
		return cmd.GenFishCompletion(writer, true)
	case "powershell", "powershell.exe":
		return cmd.GenPowerShellCompletionWithDesc(writer)
	default:
		return errors.Errorf("invalid shell: %s - must be (bash | fish | zsh | powershell(.exe))\n", shellName)
	}
}

func shellName() (string, error) {
	const shellNameArgIndex = 2
	var shellName string
	if len(os.Args) > shellNameArgIndex {
		shellName = os.Args[shellNameArgIndex]
	} else {
		shellBinary, err := loginshell.Shell()
		if err != nil {
			return "", err
		}
		shellName = path.Base(shellBinary)
	}
	return shellName, nil
}
