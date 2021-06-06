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
	// Exclude ...
	Exclude FilterKind = iota - 1
	// Alert ...
	Alert
	// Include ...
	Include
)

// Filter ...
type Filter struct {
	Pattern regexp.Regexp
	Kind    FilterKind
}
