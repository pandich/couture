package schema

import (
	"github.com/oriser/regroup"
	"regexp"
)

// Guess ..
func Guess(s string, schemasToCheck ...Schema) *Schema {
	for _, schema := range schemasToCheck {
		if schema.CanHandle(s) {
			return &schema
		}
	}
	return nil
}

const (
	// Timestamp ...
	Timestamp = "timestamp"
	// Level ...
	Level = "level"
	// Message ...
	Message = "message"
	// Application ...
	Application = "application"
	// Method ...
	Method = "method"
	// Line ...
	Line = "line"
	// Thread ...
	Thread = "thread"
	// Class ...
	Class = "class"
	// Exception ...
	Exception = "exception"
)

const (
	// JSON ...
	JSON format = "json"
	// Text ...
	Text format = "text"
)

type (
	predicate func(s string) bool

	priority = uint8

	format string

	// Schema ...
	Schema interface {
		Name() string
		Format() format
		Priority() priority
		Fields() []string
		Column(field string) (string, bool)
		Template(field string) (string, bool)
		CanHandle(s string) bool
		TextPattern() *regroup.ReGroup
	}

	baseSchema struct {
		name              string
		format            format
		priority          priority
		mapping           map[string]string
		templates         map[string]string
		inputFields       []string
		predicate         predicate
		predicatePatterns map[string]*regexp.Regexp
		textPattern       *regroup.ReGroup
	}
)

func newSchema(name string, definition definition) Schema {
	return baseSchema{
		name:              name,
		format:            definition.Format,
		priority:          definition.Priority,
		mapping:           definition.inverseMapping(),
		templates:         definition.Display,
		inputFields:       definition.inputFields(),
		predicate:         definition.canHandlePredicate(),
		predicatePatterns: definition.predicatePatterns(),
		textPattern:       definition.textPattern(),
	}
}

// Template ...
func (schema baseSchema) Template(name string) (string, bool) {
	tmpl, ok := schema.templates[name]
	return tmpl, ok
}

// Priority ...
func (schema baseSchema) Priority() priority {
	return schema.priority
}

// Name ...
func (schema baseSchema) Name() string {
	return schema.name
}

// Format ...
func (schema baseSchema) Format() format {
	return schema.format
}

// Fields ...
func (schema baseSchema) Fields() []string {
	return schema.inputFields
}

// Column ...
func (schema baseSchema) Column(field string) (string, bool) {
	s, ok := schema.mapping[field]
	return s, ok
}

// CanHandle ...
func (schema baseSchema) CanHandle(s string) bool {
	return schema.predicate(s)
}

// TextPattern ...
func (schema baseSchema) TextPattern() *regroup.ReGroup {
	return schema.textPattern
}
