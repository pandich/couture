package model

// Manager manages the lifecycle of registry, and the routing of their events to the sinks.
type Manager interface {
	// Start the Manager.
	Start() error
	// Stop the Manager.
	Stop()
	// RegisterOptions one or more sinks or registry.
	RegisterOptions(opts ...interface{}) error
}
