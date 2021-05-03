package manager

import (
	"couture/pkg/model"
	"regexp"
)

const (
	verbosityLevelError = iota
	verbosityLevelWarn
	verbosityLevelInfo
	verbosityLevelDebug
	verbosityLevelTrace
)

// VerboseDisplayOption ...
func VerboseDisplayOption(verbosity uint) interface{} {
	return baseOption{applier: func(mgr *managerOptions) {
		switch verbosity {
		case verbosityLevelTrace:
			mgr.level = model.LevelTrace
		case verbosityLevelDebug:
			mgr.level = model.LevelDebug
		case verbosityLevelInfo:
			mgr.level = model.LevelInfo
		case verbosityLevelWarn:
			mgr.level = model.LevelWarn
		case verbosityLevelError:
		default:
			mgr.level = model.LevelError
		}
	}}
}

// FilterOption ...
func FilterOption(includeFilters []*regexp.Regexp, excludeFilters []*regexp.Regexp) interface{} {
	return baseOption{applier: func(options *managerOptions) {
		options.includeFilters = includeFilters
		options.excludeFilters = excludeFilters
	}}
}

// LogLevelOption ...
func LogLevelOption(level model.Level) interface{} {
	return baseOption{applier: func(options *managerOptions) {
		options.level = level
	}}
}

type (
	// managerOptions
	managerOptions struct {
		level          model.Level
		includeFilters []*regexp.Regexp
		excludeFilters []*regexp.Regexp
	}

	// option is an entity capable of mutating the state of a managerOptions struct.
	option interface {
		Apply(manager *managerOptions)
	}

	baseOption struct {
		applier func(*managerOptions)
	}
)

// Apply ...
func (opt baseOption) Apply(mgrOptions *managerOptions) {
	opt.applier(mgrOptions)
}
