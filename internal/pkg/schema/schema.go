package schema

import (
	"github.com/tidwall/gjson"
	"regexp"
)

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

type (
	// priority ...
	priority = uint8

	// Schema ...
	Schema interface {
		Name() string
		Priority() priority
		Fields() []string
		Column(field string) (string, bool)
		CanHandle(s string) bool
	}

	baseSchema struct {
		name        string
		priority    priority
		mapping     map[string]string
		inputFields []string
		predicate   func(s string) bool
	}
)

// Priority ...
func (schema baseSchema) Priority() priority {
	return schema.priority
}

func newSchema(name string, definition definition) Schema {
	var predicateFields []string
	predicatePatterns := map[string]*regexp.Regexp{}
	for fieldName, pattern := range definition.Predicates {
		if pattern != "" {
			predicatePatterns[fieldName] = regexp.MustCompile(pattern)
		} else {
			predicatePatterns[fieldName] = nil
		}
		predicateFields = append(predicateFields, fieldName)
	}
	predicate := func(possibleJSON string) bool {
		values := gjson.GetMany(possibleJSON, predicateFields...)
		for i := range predicateFields {
			value := values[i]
			field := predicateFields[i]
			pattern := predicatePatterns[field]
			if pattern == nil {
				if !value.Exists() {
					return false
				}
			} else {
				stringValue := value.String()
				if !pattern.MatchString(stringValue) {
					return false
				}
			}
		}
		return true
	}

	var inputFields []string
	inverseMapping := map[string]string{}
	for k, v := range definition.Mapping {
		inverseMapping[v] = k
		inputFields = append(inputFields, v)
	}
	return baseSchema{
		name:        name,
		priority:    definition.Priority,
		mapping:     inverseMapping,
		inputFields: inputFields,
		predicate:   predicate,
	}
}

// Name ...
func (schema baseSchema) Name() string {
	return schema.name
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
