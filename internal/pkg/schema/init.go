package schema

import (
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
		schema.TextPattern = regroup.MustCompile(schema.predicatePatternsByField[textRootField].String())
		schema.canHandle = schema.canHandleText
	}
}

func (schema *Schema) initFields() {
	for _, field := range schema.FieldByColumn {
		schema.Fields = append(schema.Fields, field)
	}
}
