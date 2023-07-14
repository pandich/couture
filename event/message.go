package event

import "github.com/gagglepanda/couture/model"

// Message ...
type (
	// Message a message.
	Message string

	// Error a stack trace.
	Error Message
)

// Matches determines if an event Matches the filters criteria.
func (msg Message) Matches(filters *[]model.Filter) model.FilterKind {
	var hasIncludes = false
	for i := range *filters {
		filter := (*filters)[i]
		switch filter.Kind {
		case model.None:
			return model.Include
		case model.Exclude:
			if filter.Pattern.MatchString(string(msg)) {
				return model.Exclude
			}
		case model.Include:
			hasIncludes = true
			if filter.Pattern.MatchString(string(msg)) {
				return model.Include
			}
		case model.AlertOnce:
			if filter.Pattern.MatchString(string(msg)) {
				(*filters)[i].Kind = model.None
				return model.AlertOnce
			}
		}
	}
	if hasIncludes {
		return model.Exclude
	}
	return model.Include
}

// String ...
func (msg Message) String() string {
	return string(msg)
}
