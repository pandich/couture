package model

// Message ...
type (
	// Message a message.
	Message string

	// Exception a stack trace.
	Exception Message
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
				// TODO cleanup alert/include logic
				// this is an ugly hack to allow for an alert to fire once and then turn into a normal include
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
