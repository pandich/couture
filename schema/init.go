package schema

import (
	"fmt"
	"github.com/oriser/regroup"
	"regexp"
)

func (schema *Schema) init(name string) {
	schema.Name = name
	schema.initPredicatePatterns()
	schema.initCanHandle()
	schema.initFields()
}

func (schema *Schema) initPredicatePatterns() {
	schema.predicatePatternsByField = map[string]*regexp.Regexp{}
	for field, pattern := range schema.PredicatesByField {
		schema.predicateFields = append(schema.predicateFields, field)
		if pattern != "" {
			schema.predicatePatternsByField[field] = regexp.MustCompile(pattern)
		} else {
			schema.predicatePatternsByField[field] = nil
		}
	}
}

func (schema *Schema) initCanHandle() {
	switch schema.Format {
	case JSON:
		schema.canHandle = schema.canHandleJSON
	case Text:
		var pattern = schema.predicatePatternsByField[textRootField].String()
		re := regexp.MustCompile(pattern)
		names := map[string]bool{
			Timestamp:   false,
			Level:       false,
			Message:     false,
			Application: false,
			Action:      false,
			Line:        false,
			Context:     false,
			Entity:      false,
			Error:       false,
		}
		for i := 0; i < re.NumSubexp(); i++ {
			names[re.SubexpNames()[i+1]] = true
		}

		for name, ok := range names {
			if !ok {
				pattern += fmt.Sprintf("(?P<%s>)", name)
			}
		}
		schema.TextPattern = regroup.MustCompile(pattern)
		schema.canHandle = schema.canHandleText
	}
}

func (schema *Schema) initFields() {
	for _, field := range schema.FieldByColumn {
		schema.Fields = append(schema.Fields, field)
	}
}
