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
		InputFields() []string
		Mapping() map[string]string
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

// InputFields ...
func (b baseSchema) InputFields() []string {
	return b.inputFields
}

// Mapping ...
func (b baseSchema) Mapping() map[string]string {
	return b.mapping
}

// NewSchema ...
func NewSchema(definition Definition) Schema {
	var inputFields []string
	for inputField := range definition.Mapping {
		inputFields = append(inputFields, inputField)
	}
	var predicateFields []string
	for f := range definition.Predicates {
		predicateFields = append(predicateFields, f)
	}
	predicate := func(s string) bool {
		if !gjson.Valid(s) {
			return false
		}
		values := gjson.GetMany(s, predicateFields...)
		for i := range predicateFields {
			predicateField := predicateFields[i]
			value := values[i]
			test := regexp.MustCompile(definition.Predicates[predicateField])
			if !test.MatchString(value.String()) {
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
