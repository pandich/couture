package schema

import (
	"github.com/tidwall/gjson"
	"strings"
)

const textRootField = "_"

// GuessSchema ..
func GuessSchema(s string, schemasToCheck ...Schema) *Schema {
	for _, schema := range schemasToCheck {
		if schema.canHandle(s) {
			return &schema
		}
	}
	return nil
}

func (schema *Schema) canHandleJSON(s string) bool {
	values := gjson.GetMany(s, schema.predicateFields...)
	for i := range schema.predicateFields {
		field := schema.predicateFields[i]

		pattern := schema.predicatePatternsByField[field]
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

func (schema *Schema) canHandleText(s string) bool {
	return schema.predicatePatternsByField[textRootField].MatchString(strings.TrimRight(s, "\n"))
}
