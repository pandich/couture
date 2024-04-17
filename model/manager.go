package model

import (
	"regexp"
)

// Manager manages the lifecycle of registry, and the routing of their events to the sinks.
type Manager interface {
	// Run ...
	Run() error
	// Start the Manager.
	Start() error
	// Stop the Manager.
	Stop()
	// Wait for the manager to stop.
	Wait()
	// Register one or more sinks or registry.
	Register(opts ...interface{}) error
	// TrapSignals ...
	TrapSignals()
}

// FilterKind ...
type FilterKind int

const (
	// None ...
	None FilterKind = iota
	// Exclude ...
	Exclude
	// AlertOnce ...
	AlertOnce
	// Include ...
	Include
)

// Filter ...
type Filter struct {
	Pattern regexp.Regexp
	Kind    FilterKind
}

// Filters provides a mechansim to filter events in one of the following ways:
//   - None: Do not filter the event.
//   - Exclude: Filter the event by excluding it.
//   - AlertOnce: Filter the event by alerting it once and only once.
//   - Include: Filter the event by including it.
type Filters []Filter

func (f FilterKind) isHighlighted() bool {
	return f == Include
}

func (f Filter) replaceAllStringFunc(s string, replacer func(string) string) string {
	if f.Kind.isHighlighted() {
		s = f.Pattern.ReplaceAllStringFunc(s, replacer)
	}
	return s
}

// ReplaceAllStringFunc ...
func (fs Filters) ReplaceAllStringFunc(s string, replacer func(string) string) string {
	for _, filter := range fs {
		s = filter.replaceAllStringFunc(s, replacer)
	}
	return s
}
