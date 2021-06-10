package level

import (
	"strings"
)

const (
	// Trace log level for tracing
	Trace Level = "TRACE"
	// Debug log level for debugging
	Debug Level = "DEBUG"
	// Info log level for information
	Info Level = "INFO"
	// Warn log level for warnings
	Warn Level = "WARN"
	// Error log level for errors
	Error Level = "ERROR"
)

const (
	tracePriority int = iota
	debugPriority
	infoPriority
	warnPriority
	errorPriority
)

// Levels ...
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

// Default ...
var Default = Info

// LowerNames ...
func LowerNames() []string {
	var logLevels []string
	for _, level := range Levels {
		logLevels = append(logLevels, level.LowerName())
	}
	return logLevels
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

// LowerName returns a name that is more simple for the user to enter.
func (level Level) LowerName() string {
	return strings.ToLower(string(level))
}

// ByName ...
func ByName(s string, defaultLevel Level) Level {
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
		return defaultLevel
	}
}
