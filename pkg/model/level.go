package model

import (
	"strings"
)

const (
	// TraceLevel log level for tracing
	TraceLevel Level = "TRACE"
	// DebugLevel log level for debugging
	DebugLevel Level = "DEBUG"
	// InfoLevel log level for information
	InfoLevel Level = "INFO"
	// WarnLevel log level for warnings
	WarnLevel Level = "WARN"
	// ErrorLevel log level for errors
	ErrorLevel Level = "ERROR"
	// disabledLevel disabled log level.
	disabledLevel Level = ""
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
	TraceLevel:    levelTracePriority,
	DebugLevel:    levelDebugPriority,
	InfoLevel:     levelInfoPriority,
	WarnLevel:     levelWarnPriority,
	ErrorLevel:    levelErrorPriority,
	disabledLevel: levelDisabledPriority,
}

// IsValidLevel ...
func IsValidLevel(s string) bool {
	_, ok := levelPriorities[Level(strings.ToUpper(s))]
	return ok
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
