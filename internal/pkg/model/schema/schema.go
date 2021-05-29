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
		Fields() []string
		Column(field string) (string, bool)
		Test(s string) bool
	}

	baseSchema struct {
		mapping     map[string]string
		inputFields []string
		predicate   func(s string) bool
	}

	// Definition ...
	Definition struct {
		Predicates map[string]string `json:"predicates"`
		Mapping    map[string]string `json:"mapping"`
	}
)

// NewSchema ...
func NewSchema(definition Definition) Schema {
	var inputFields []string
	for inputField := range definition.Mapping {
		inputFields = append(inputFields, inputField)
	}
	var predicateFields []string
	predicatePatterns := map[string]*regexp.Regexp{}
	for fieldName, pattern := range definition.Predicates {
		predicatePatterns[fieldName] = regexp.MustCompile(pattern)
		predicateFields = append(predicateFields, fieldName)
	}
	predicate := func(s string) bool {
		if !gjson.Valid(s) {
			return false
		}
		values := gjson.GetMany(s, predicateFields...)
		for i := range predicateFields {
			value := values[i]
			field := predicateFields[i]
			pattern := predicatePatterns[field]
			if !pattern.MatchString(value.String()) {
				return false
			}
		}
		return true
	}
	return baseSchema{
		mapping:     definition.Mapping,
		inputFields: inputFields,
		predicate:   predicate,
	}
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

// Test ...
func (b baseSchema) Test(s string) bool {
	return b.predicate(s)
}
