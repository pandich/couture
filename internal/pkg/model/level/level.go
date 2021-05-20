package level

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
