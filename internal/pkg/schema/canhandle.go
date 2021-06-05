package schema

import (
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

// Guess ..
func Guess(s string, schemasToCheck ...Schema) *Schema {
	for _, schema := range schemasToCheck {
		if schema.canHandle(s) {
			return &schema
		}
	}
	const envKey = "COUTURE_DIE_ON_UNKNOWN"
	const exitCode = 12
	if os.Getenv(envKey) != "" {
		fmt.Printf("unknown: %+v\n", s)
		os.Exit(exitCode)
	}
	return nil
}

func (schema Schema) canHandleJSON(s string) bool {
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

func (schema Schema) canHandleText(s string) bool {
	return schema.predicatePatternsByField[textRootField].MatchString(strings.TrimRight(s, "\n"))
}
