package model

import (
	"strings"
)

// LevelTrace ...
const (
	// LevelTrace log level for tracing
	LevelTrace Level = "TRACE"
	// LevelDebug log level for debugging
	LevelDebug Level = "DEBUG"
	// LevelInfo log level for information
	LevelInfo Level = "INFO"
	// LevelWarn log level for warnings
	LevelWarn Level = "WARN"
	// LevelError log level for errors
	LevelError Level = "ERROR"
	// levelDisabled disabled log level.
	levelDisabled Level = ""
)

const (
	levelDisabledPriority levelPriority = iota
	levelTracePriority
	levelDebugPriority
	levelInfoPriority
	levelWarnPriority
	levelErrorPriority
)

// levelPriorities relative level priorities
var levelPriorities = map[Level]levelPriority{
	LevelTrace:    levelTracePriority,
	LevelDebug:    levelDebugPriority,
	LevelInfo:     levelInfoPriority,
	LevelWarn:     levelWarnPriority,
	LevelError:    levelErrorPriority,
	levelDisabled: levelDisabledPriority,
}

type (
	// levelPriority is the relative priority of levels.
	levelPriority uint8
	// Level a log level.
	Level string
)

// IsAtLeast determines if the specified level is at least as high as this level.
func (level Level) IsAtLeast(l Level) bool {
	return l.priority() <= level.priority()
}

// priority is the relative priority of this level.
func (level Level) priority() levelPriority {
	return levelPriorities[level]
}

// Normalize corrects case and anything else to normalize the log level.
func (level Level) Normalize() Level {
	return Level(strings.ToUpper(string(level)))
}
