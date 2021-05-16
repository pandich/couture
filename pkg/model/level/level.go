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

var priorities = map[Level]int{
	Trace: tracePriority,
	Debug: debugPriority,
	Info:  infoPriority,
	Warn:  warnPriority,
	Error: errorPriority,
}

var priorityLevels = map[int]Level{
	tracePriority: Trace,
	debugPriority: Debug,
	infoPriority:  Info,
	warnPriority:  Warn,
	errorPriority: Error,
}

// New ...
func New(s string) (Level, bool) {
	lvl := Level(strings.ToUpper(s))
	if _, ok := priorities[lvl]; ok {
		return lvl, true
	}
	return "", false
}

// Level a log level.
type Level string

// IsAtLeast determines if the specified level is at least as high as this level.
func (level Level) IsAtLeast(l Level) bool {
	return level.Priority() >= l.Priority()
}

// Priority is the relative Priority of this level.
func (level Level) Priority() int {
	return priorities[level]
}

// ByPriority ...
func ByPriority(priority int) Level {
	level, ok := priorityLevels[priority]
	if ok {
		return level
	}
	return Error
}
