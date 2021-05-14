package cli

import (
	"couture/internal/pkg/manager"
	"couture/pkg/model"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/pflag"
	"gopkg.in/multierror.v1"
	"net/url"
	"regexp"
	"time"
)

func filterOption(persistent *pflag.FlagSet) (interface{}, error) {
	filters := func(persistent *pflag.FlagSet, key string) ([]*regexp.Regexp, error) {
		filterStrings, err := persistent.GetStringSlice(key)
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
	includes, err := filters(persistent, includeFilterFlag)
	if err != nil {
		return nil, err
	}
	excludes, err := filters(persistent, excludeFilterFlag)
	if err != nil {
		return nil, err
	}
	return manager.FilterOption(includes, excludes), nil
}

func sinceOption(persistent *pflag.FlagSet) (interface{}, error) {
	var err error
	var d time.Duration
	sinceString, err := persistent.GetString(sinceFlag)
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
	return nil, errors2.Errorf("invalid timestamp or duration: %s", sinceString)
}

func verbosityOption(persistent *pflag.FlagSet) (interface{}, error) {
	verbose, err := persistent.GetBool(verboseFlag)
	if err != nil {
		return nil, err
	}
	quiet, err := persistent.GetBool(quietFlag)
	if err != nil {
		return nil, err
	}
	if verbose && quiet {
		return nil, errors2.Errorf("verbose and quiet are mutually exclusive")
	}
	var verbosity = 1
	if verbose {
		verbosity++
	}
	if quiet {
		verbosity--
	}
	return manager.VerboseDisplayOption(uint(verbosity)), nil
}

func wrapOption(persistent *pflag.FlagSet) (interface{}, error) {
	wrap, err := persistent.GetInt(wrapFlag)
	if err != nil {
		return nil, err
	}
	if wrap < 0 {
		return nil, errors2.Errorf("bad wrap width: %d", wrap)
	}
	return manager.WrapOption(wrap), nil
}

func levelOption(persistent *pflag.FlagSet) (interface{}, error) {
	level, err := persistent.GetString(levelFlag)
	if err != nil {
		return nil, err
	}

	if !model.IsValidLevel(level) {
		return nil, errors2.Errorf("invalid level: %s", level)
	}
	return manager.LogLevelOption(model.Level(level)), nil
}

func sourceOptions(sourceStrings []string) ([]interface{}, error) {
	if len(sourceStrings) == 0 {
		return nil, errors2.Errorf("no source URLs provided")
	}
	var violations []error
	var configuredSources []interface{}
	for _, sourceArgs := range sourceStrings {
		u, err := url.Parse(sourceArgs)
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
			violations = append(violations, errors2.Errorf("invalid source URL: %+v", sourceURL))
		}
	}
	if len(violations) > 0 {
		return nil, multierror.New(violations)
	}
	return configuredSources, nil
}
