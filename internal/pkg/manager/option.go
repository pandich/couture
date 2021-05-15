package manager

import (
	"couture/pkg/model/level"
	"regexp"
	"time"
)

// SinceOption ...
func SinceOption(t time.Time) interface{} {
	return baseOption{applier: func(mgr *managerOptions) {
		mgr.since = &t
	}}
}

// VerboseDisplayOption ...
func VerboseDisplayOption(level level.Level) interface{} {
	return baseOption{applier: func(mgr *managerOptions) {
		mgr.level = level
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
func LogLevelOption(level level.Level) interface{} {
	return baseOption{applier: func(options *managerOptions) {
		options.level = level
	}}
}

// WrapOption ...
func WrapOption(width int) interface{} {
	return baseOption{applier: func(options *managerOptions) {
		if width > 0 {
			options.wrap = &width
		}
	}}
}

type (
	// managerOptions
	managerOptions struct {
		level          level.Level
		wrap           *int
		since          *time.Time
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
