package model

import (
	"regexp"
)

// Message ...
type (
	// Message a message.
	Message string

	// Exception a stack trace.
	Exception Message
)

// Matches determines if an event Matches the filters criteria.
func (msg Message) Matches(include []regexp.Regexp, exclude []regexp.Regexp) bool {
	var shouldInclude = false

	// process the includes returning true on the first match
	for _, filter := range include {
		if filter.MatchString(string(msg)) {
			shouldInclude = true
		}
	}
	// if we made it this far and have include filters, none of them matched, so we return false
	if len(include) > 0 && !shouldInclude {
		return false
	}

	// process the excludes returning false on the first match
	for _, filter := range exclude {
		if filter.MatchString(string(msg)) {
			return false
		}
	}
	return true
}

// String ...
func (msg Message) String() string {
	return string(msg)
}
