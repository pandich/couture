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
		case Exclude:
			if filter.Pattern.MatchString(string(msg)) {
				return filter.Kind
			}
		case Include, Alert:
			hasIncludes = true
			if filter.Pattern.MatchString(string(msg)) {
				kind := filter.Kind
				// downgrade the alert to an include after it fires - this feels a bit hacky
				(*filters)[i].Kind = Include
				return kind
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
