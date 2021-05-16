package cli

import (
	"couture/internal/pkg/manager"
	"couture/pkg/model"
	"couture/pkg/model/level"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/pflag"
	"gopkg.in/multierror.v1"
	"net/url"
	"regexp"
	"time"
)

func filterOption(persistent *pflag.FlagSet) (interface{}, error) {
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

func filters(persistent *pflag.FlagSet, flagName string) ([]*regexp.Regexp, error) {
	if !isFlagSet(persistent, flagName) {
		return []*regexp.Regexp{}, nil
	}
	filterStrings, err := persistent.GetStringSlice(flagName)
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

func isFlagSet(persistent *pflag.FlagSet, key string) bool {
	var found = false
	persistent.VisitAll(func(f *pflag.Flag) {
		if f.Name == key {
			found = true
		}
	})
	return found
}

func sinceOption(persistent *pflag.FlagSet) (interface{}, error) {
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
	verboseFlag, err := persistent.GetCount(verboseFlag)
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

func wrapOption(persistent *pflag.FlagSet) (interface{}, error) {
	wrap, err := persistent.GetUint(wrapFlag)
	if err != nil {
		return nil, err
	}
	return manager.WrapOption(wrap), nil
}

func rateLimitOption(persistent *pflag.FlagSet) (interface{}, error) {
	const minRateLimit = 100
	const maxRateLimit = 10_000
	rateLimit, err := persistent.GetUint(rateLimitFlag)
	if err != nil {
		return nil, err
	}
	if rateLimit < minRateLimit || rateLimit > maxRateLimit {
		return nil, errors2.Errorf("bad rate limit: %d - must be in (%d, %d)", rateLimit, minRateLimit, maxRateLimit)
	}
	return manager.RateLimitOption(rateLimit), nil
}

func levelOption(persistent *pflag.FlagSet) (interface{}, error) {
	levelName, err := persistent.GetString(levelFlag)
	if err != nil {
		return nil, err
	}
	if lvl, ok := level.New(levelName); ok {
		return manager.LogLevelOption(lvl), nil
	}
	return nil, errors2.Errorf("invalid levelName: %s", levelName)
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
