package manager

import (
	"couture/internal/pkg/model/level"
	"regexp"
	"time"
)

// SinceOption ...
func SinceOption(t time.Time) interface{} {
	return baseOption{applier: func(options *managerOptions) {
		options.since = &t
	}}
}

// FilterOption ...
func FilterOption(includeFilters []regexp.Regexp, excludeFilters []regexp.Regexp) interface{} {
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

type (
	// managerOptions
	managerOptions struct {
		level          level.Level
		since          *time.Time
		includeFilters []regexp.Regexp
		excludeFilters []regexp.Regexp
	}

	// option is an entity capable of mutating the state of a managerOptions struct.
	option interface {
		Apply(options *managerOptions)
	}

	baseOption struct {
		applier func(*managerOptions)
	}
)

// Apply ...
func (opt baseOption) Apply(mgrOptions *managerOptions) {
	opt.applier(mgrOptions)
}
