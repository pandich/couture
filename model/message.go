package model

// Message ...
type (
	// Message a message.
	Message string

	// Error a stack trace.
	Error Message
)

// Matches determines if an event Matches the filters criteria.
func (msg Message) Matches(filters *[]Filter) FilterKind {
	var hasIncludes = false
	for i := range *filters {
		filter := (*filters)[i]
		switch filter.Kind {
		case none:
			return Include
		case Exclude:
			if filter.Pattern.MatchString(string(msg)) {
				return Exclude
			}
		case Include:
			hasIncludes = true
			if filter.Pattern.MatchString(string(msg)) {
				return Include
			}
		case AlertOnce:
			if filter.Pattern.MatchString(string(msg)) {
				(*filters)[i].Kind = none
				return AlertOnce
			}
		}
	}
	if hasIncludes {
		return Exclude
	}
	return Include
}

// String ...
func (msg Message) String() string {
	return string(msg)
}
