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

// Schema ...
type (
	// Schema ...
	Schema interface {
		Name() string
		Fields() []string
		Column(field string) (string, bool)
		CanHandle(s string) bool
	}

	baseSchema struct {
		name        string
		mapping     map[string]string
		inputFields []string
		predicate   func(s string) bool
	}

	definition struct {
		Predicates map[string]string `json:"predicates"`
		Mapping    map[string]string `json:"mapping"`
	}
)

func newSchema(name string, definition definition) Schema {
	var inputFields []string
	for inputField := range definition.Mapping {
		inputFields = append(inputFields, inputField)
	}
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
	return baseSchema{
		name:        name,
		mapping:     definition.Mapping,
		inputFields: inputFields,
		predicate:   predicate,
	}
}

// Name ...
func (b baseSchema) Name() string {
	return b.name
}

// Fields ...
func (b baseSchema) Fields() []string {
	return b.inputFields
}

// Column ...
func (b baseSchema) Column(field string) (string, bool) {
	s, ok := b.mapping[field]
	return s, ok
}

// CanHandle ...
func (b baseSchema) CanHandle(s string) bool {
	return b.predicate(s)
}
