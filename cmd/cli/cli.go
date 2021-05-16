package cli

import (
	"couture/internal/pkg/source"
	"couture/internal/pkg/source/aws/cloudformation"
	"couture/internal/pkg/source/aws/cloudwatch"
	"couture/internal/pkg/source/elasticsearch"
	"couture/internal/pkg/source/fake"
	"couture/internal/pkg/source/ssh"
	"couture/internal/pkg/source/tail"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

const (
	excludeFilterFlag = "exclude"
	includeFilterFlag = "include"
	levelFlag         = "level"
	rateLimitFlag     = "rate-limit"
	sinceFlag         = "since"
	verboseFlag       = "verbose"
	wrapFlag          = "wrap"
)

var couture = &cobra.Command{
	Use:   "couture [flags] source_url ...",
	Short: "Tails one or more event sources.\n",
	Long: "Description:\n\nTails one or more event sources.\n" +
		"When providing a CloudFormation stack, resources are recursively analyzed until all loggable entities are found. " +
		"This includes the stack events of the stack itself, as well as any log groups " +
		"its entities contain.\n",
	Example: strings.Join([]string{"\n  couture " + strings.Join(sourceMetadata.ExampleURLs(), "\n  couture ")}, ""),
	Args:    cobra.MinimumNArgs(1),
	RunE:    runner,
}

func setupFlags(persistent *pflag.FlagSet) {
	const noWrap = 0
	const defaultRateLimit = 1_000
	persistent.CountP(verboseFlag, "v", "Display additional diagnostic data.")
	persistent.UintP(rateLimitFlag, "r", defaultRateLimit, "Max events per second to process.")
	persistent.UintP(wrapFlag, "w", noWrap, "Display no diagnostic data.")
	persistent.StringP(levelFlag, "l", "info", "Minimum log level to display (trace, debug, info warn, error.")
	persistent.StringP(sinceFlag, "s", "5m", "How far back in time to search for events.")
	persistent.StringSliceP(includeFilterFlag, "i", []string{}, "Include filter regular expressions. Performed before excludes.")
	persistent.StringSliceP(excludeFilterFlag, "e", []string{}, "Exclude filter regular expressions. Performed after includes.")
	couture.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		viper.SetConfigName(".couture")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		err := viper.ReadInConfig()
		target := &viper.ConfigFileNotFoundError{}
		if err != nil && !errors.As(err, &target) {
			return errors.Errorf("fatal error config file: %s\n", err)
		}
		return nil
	}
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
