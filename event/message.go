package event

import "github.com/pandich/couture/model"

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
				// explicit exclude
				return model.Exclude
			}

		case model.Include:
			hasIncludes = true

			if filter.Pattern.MatchString(string(msg)) {
				// exlplicit include
				return model.Include
			}

		case model.AlertOnce:
			if filter.Pattern.MatchString(string(msg)) {
				// explicit alert
				(*filters)[i].Kind = model.None
				return model.AlertOnce
			}
		}
	}

	// if it has includes, but never hit a matching include, it is an implicit exclusion.
	if hasIncludes {
		return model.Exclude
	}

	// otherwise it is an implicit include
	return model.Include
}

// String ...
func (msg Message) String() string {
	return string(msg)
}
