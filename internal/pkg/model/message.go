package model

// Message ...
type (
	// Message a message.
	Message string

	// Exception a stack trace.
	Exception Message
)

// Matches determines if an event Matches the filters criteria.
func (msg Message) Matches(filters []Filter) bool {
	var hasIncludes = false
	for _, filter := range filters {
		if filter.ShouldInclude {
			hasIncludes = true
		}
		if filter.Pattern.MatchString(string(msg)) {
			return filter.ShouldInclude
		}
	}
	return !hasIncludes
}

// String ...
func (msg Message) String() string {
	return string(msg)
}
