package cli

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/internal/pkg/source/ssh"
	"couture/internal/pkg/source/tail"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	version = "0.1.0"
)

const (
	docCommand      = "__doc"
	completeCommand = "__complete"
)

const (
	excludeFilterFlag = "exclude"
	includeFilterFlag = "include"
	levelFlag         = "level"
	quietFlag         = "quiet"
	sinceFlag         = "since"
	verboseFlag       = "verbose"
	wrapFlag          = "wrap"
)

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
	coutureCmd := &cobra.Command{
		Version: version,
		Use:     "coutureCmd [flags] source_url ...",
		Short:   "Tails one or more event sourceMetadata.\n",
		Long: "Description:\n\nTails one or more event sources.\n" +
			"When providing a CloudFormation stack, resources are recursively analyzed until all loggable entities are found. " +
			"This includes the stack events of the stack itself, as well as any log groups " +
			"its entities contain.\n",
		Example: strings.Join([]string{"\n  " + strings.Join(sourceMetadata.ExampleURLs(), "\n  ")}, ""),
		Args:    cobra.MinimumNArgs(1),
		RunE:    runner,
	}
	persistent := coutureCmd.PersistentFlags()
	persistent.BoolP(verboseFlag, "v", false, "Display additional diagnostic data.")
	persistent.BoolP(quietFlag, "q", false, "Display no diagnostic data.")
	persistent.IntP(wrapFlag, "w", 72, "Display no diagnostic data.")
	persistent.StringP(levelFlag, "l", "info", "Minimum log level to display (trace, debug, info warn, error.")
	persistent.StringP(sinceFlag, "s", "5m", "How far back in time to search for events.")
	persistent.StringSliceP(excludeFilterFlag, "e", []string{}, "Exclude filter regular expressions. Always performed after includes.")

	if (len(os.Args) == 2 || len(os.Args) == 3) && os.Args[1] == completeCommand {
		return handleCompleteCommand(coutureCmd)
	}
	if len(os.Args) == 3 && os.Args[1] == docCommand {
		return handleDocCommand(coutureCmd)
	}
	return handleLogCommand(coutureCmd)
}
