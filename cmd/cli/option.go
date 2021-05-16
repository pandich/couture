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
	"net/url"
	"regexp"
	"strings"
	"time"
)

const noWrap = 0

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

func pagerOption(flags *pflag.FlagSet) (*interface{}, error) {
	pager, err := flags.GetString(pagerFlag)
	if err != nil {
		return nil, err
	}
	pager = strings.Trim(pager, " \t")
	if pager == "" {
		return nil, nil
	}
	option := manager.PagerOption(pager)
	return &option, nil
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

func sourceOptions(sourceStrings []string) ([]interface{}, error) {
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

func getOptions(flags *pflag.FlagSet) ([]interface{}, error) {
	verbosityOption, err := verbosityOption(flags)
	if err != nil {
		return nil, err
	}
	wrapOption, err := wrapOption(flags)
	if err != nil {
		return nil, err
	}
	filterOption, err := filterOption(flags)
	if err != nil {
		return nil, err
	}
	levelOption, err := levelOption(flags)
	if err != nil {
		return nil, err
	}
	sinceOption, err := sinceOption(flags)
	if err != nil {
		return nil, err
	}
	pagerOption, err := pagerOption(flags)
	if err != nil {
		return nil, err
	}
	outputFormat, err := flags.GetString(outputFormatFlag)
	if err != nil {
		return nil, err
	}

	var options = []interface{}{
		verbosityOption,
		filterOption,
		levelOption,
		wrapOption,
	}

	if pagerOption != nil {
		options = append(options, *pagerOption)
	}

	if sinceOption != nil {
		options = append(options, sinceOption)
	}
	switch outputFormat {
	case "json":
		options = append(options, json.New())
	case "pretty":
		options = append(options, pretty.New())
	default:
		return nil, errors2.Errorf("unknown output format: %s", outputFormat)
	}
	return options, nil
}
