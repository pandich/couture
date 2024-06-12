// Package level provides an event-source agnostic log level system. Each source.Source is required to map its
// implementation-specific levels into a Level.
package level

const (
	// Trace log level.
	Trace Level = "TRACE"
	// Debug log level.
	Debug Level = "DEBUG"
	// Info log level.
	Info Level = "INFO"
	// Warn log level.
	Warn Level = "WARN"
	// Error log level.
	Error Level = "ERROR"
	// Default log level.
	Default = Info
)

const (
	tracePriority int = iota
	debugPriority
	infoPriority
	warnPriority
	errorPriority
)

// Levels in order of increasing significance.
var Levels = []Level{
	Trace,
	Debug,
	Info,
	Warn,
	Error,
}

var priorities = map[Level]int{
	Trace: tracePriority,
	Debug: debugPriority,
	Info:  infoPriority,
	Warn:  warnPriority,
	Error: errorPriority,
}

// Level a log level.
type Level string

// IsAtLeast determines if the specified level is at least as high as this level.
func (level Level) IsAtLeast(l Level) bool {
	return level.priority() >= l.priority()
}

// priority is the relative priority of this level.
func (level Level) priority() int {
	return priorities[level]
}

// ByName looks up the log level.
func ByName(s string) Level {
	switch Level(s) {
	case Trace:
		return Trace
	case Debug:
		return Debug
	case Info:
		return Info
	case Warn:
		return Warn
	case Error:
		return Error
	default:
		return Default
	}
}
