package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
	"couture/pkg/model"
	"couture/pkg/model/level"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/multierror.v1"
	"io"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	excludeFilterFlag = "exclude"
	includeFilterFlag = "include"
	levelFlag         = "level"
	outputFormatFlag  = "format"
	paginatorFlag     = "pager"
	noPaginatorFlag   = "no-pager"
	sinceFlag         = "since"
	verboseFlag       = "verbose"
	wrapFlag          = "wrap"
)

const (
	paginatorEnvKey    = "COUTURE_PAGER"
	paginatorConfigKey = "paginator"
)

var optionCreators = []func(*pflag.FlagSet) (interface{}, error){
	verbosityOption,
	wrapOption,
	filterOption,
	levelOption,
	sinceOption,
	sinkOption,
}

func allOptions(flags *pflag.FlagSet, args []string) ([]interface{}, error) {
	options, err := processFlags(flags)
	if err != nil {
		return nil, err
	}
	sources, err := processArgs(args)
	if err != nil {
		return nil, err
	}
	options = append(options, sources...)
	return options, nil
}

func sinkOption(flags *pflag.FlagSet) (interface{}, error) {
	out, err := getStdout(flags)
	if err != nil {
		return nil, err
	}
	outputFormat, err := flags.GetString(outputFormatFlag)
	if err != nil {
		return nil, err
	}
	switch outputFormat {
	case "json":
		return json.New(out), nil
	case "pretty":
		return pretty.New(out), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s", outputFormat)
	}
}

func filterOption(flags *pflag.FlagSet) (interface{}, error) {
	includes, err := filters(flags, includeFilterFlag)
	if err != nil {
		return nil, err
	}
	excludes, err := filters(flags, excludeFilterFlag)
	if err != nil {
		return nil, err
	}
	return manager.FilterOption(includes, excludes), nil
}

func sinceOption(flags *pflag.FlagSet) (interface{}, error) {
	var d time.Duration
	sinceString, err := flags.GetString(sinceFlag)
	if err != nil {
		return nil, err
	}
	if sinceString == "" {
		return nil, nil
	}
	d, err = time.ParseDuration(sinceString)
	if err == nil {
		t := time.Now().Add(-d)
		return manager.SinceOption(t), nil
	}

	var t time.Time
	t, err = dateparse.ParseAny(sinceString)
	if err == nil {
		return manager.SinceOption(t), nil
	}
	return nil, errors2.Errorf("invalid timestamp or duration: %s\n", sinceString)
}

func verbosityOption(flags *pflag.FlagSet) (interface{}, error) {
	verboseFlag, err := flags.GetCount(verboseFlag)
	if err != nil {
		return nil, err
	}
	var verbosity = level.Warn.Priority() + verboseFlag
	if verbosity > level.Trace.Priority() {
		verbosity = level.Trace.Priority()
	}
	lvl := level.ByPriority(verbosity)
	return manager.VerboseDisplayOption(lvl), nil
}

func wrapOption(flags *pflag.FlagSet) (interface{}, error) {
	wrap, err := flags.GetUint(wrapFlag)
	if err != nil {
		return nil, err
	}
	return manager.WrapOption(wrap), nil
}

func levelOption(flags *pflag.FlagSet) (interface{}, error) {
	levelName, err := flags.GetString(levelFlag)
	if err != nil {
		return nil, err
	}
	if lvl, ok := level.New(levelName); ok {
		return manager.LogLevelOption(lvl), nil
	}
	return nil, errors2.Errorf("invalid levelName: %s\n", levelName)
}

func processFlags(flags *pflag.FlagSet) ([]interface{}, error) {
	var options []interface{}
	for _, creator := range optionCreators {
		option, err := creator(flags)
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}
	}
	return options, nil
}

func processArgs(sourceStrings []string) ([]interface{}, error) {
	const aliasesKey = "aliases"
	aliases := viper.GetStringMapString(aliasesKey)

	if len(sourceStrings) == 0 {
		return nil, errors2.Errorf("no source URLs provided\n")
	}
	var violations []error
	var configuredSources []interface{}
	var sourceString string
	for _, sourceString = range sourceStrings {
		if alias, ok := aliases[sourceString]; ok {
			sourceString = alias
		}
		u, err := url.Parse(sourceString)
		if err != nil {
			return nil, err
		}
		sourceURL := model.SourceURL(*u)
		var handled bool
		for _, metadata := range sourceMetadata {
			if !metadata.CanHandle(sourceURL) {
				continue
			}
			handled = true
			configuredSource, err := metadata.Creator(sourceURL)
			if err != nil {
				violations = append(violations, err)
			} else {
				configuredSources = append(configuredSources, *configuredSource)
			}
		}
		if !handled {
			violations = append(violations, errors2.Errorf("invalid source URL: %+v\n", sourceURL))
		}
	}
	if len(violations) > 0 {
		return nil, multierror.New(violations)
	}
	return configuredSources, nil
}

func getStdout(flags *pflag.FlagSet) (io.Writer, error) {
	defaultOut := os.Stdout

	noPaginate, err := flags.GetBool(noPaginatorFlag)
	if err != nil {
		return nil, err
	}
	if noPaginate {
		return defaultOut, nil
	}

	var pager = viper.GetString(paginatorEnvKey)
	if pager == "" {
		pager = viper.GetString(paginatorConfigKey)
	}
	if pager == "" {
		var err error
		pager, err = flags.GetString(paginatorFlag)
		if err != nil {
			return nil, err
		}
	}
	if pager == "" {
		return defaultOut, nil
	}

	pagerArgs := strings.Split(pager, " ")
	//nolint:gosec
	cmd := exec.Command(pagerArgs[0], pagerArgs[1:]...)

	// I/O
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	out, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}
	return out, nil
}

func filters(flags *pflag.FlagSet, flagName string) ([]*regexp.Regexp, error) {
	if !isFlagSet(flags, flagName) {
		return []*regexp.Regexp{}, nil
	}
	filterStrings, err := flags.GetStringSlice(flagName)
	if err != nil {
		return nil, err
	}
	var filters []*regexp.Regexp
	for _, filterString := range filterStrings {
		filter, err := regexp.Compile(filterString)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func isFlagSet(flags *pflag.FlagSet, key string) bool {
	var found = false
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Name == key {
			found = true
		}
	})
	return found
}