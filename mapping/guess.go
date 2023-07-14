package mapping

import (
	"github.com/tidwall/gjson"
	"strings"
)

const textRootField = "_"

// GuessMapping ..
func GuessMapping(s string, mappingsToCheck ...Mapping) *Mapping {
	for _, mapping := range mappingsToCheck {
		if mapping.canHandle(s) {
			return &mapping
		}
	}
	return nil
}

func (mapping *Mapping) canHandleJSON(s string) bool {
	values := gjson.GetMany(s, mapping.predicateFields...)
	for i := range mapping.predicateFields {
		field := mapping.predicateFields[i]

		pattern := mapping.predicatePatternsByField[field]
		if pattern == nil {
			return false
		}

		value := values[i]
		if !value.Exists() {
			return false
		}
		stringValue := value.String()
		if !pattern.MatchString(stringValue) {
			return false
		}
	}
	return true
}

func (mapping *Mapping) canHandleText(s string) bool {
	return mapping.predicatePatternsByField[textRootField].MatchString(strings.TrimRight(s, "\n"))
}
