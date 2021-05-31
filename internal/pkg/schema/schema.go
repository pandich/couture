package schema

import (
	"github.com/oriser/regroup"
	errors2 "github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
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

	// priority ...
	priority = uint8

	format string

	// Schema ...
	Schema interface {
		Name() string
		Format() format
		Priority() priority
		Fields() []string
		Column(field string) (string, bool)
		CanHandle(s string) bool
		TextPattern() *regroup.ReGroup
	}

	baseSchema struct {
		name              string
		format            format
		priority          priority
		mapping           map[string]string
		inputFields       []string
		predicate         predicate
		predicatePatterns map[string]*regexp.Regexp
		textPattern       *regroup.ReGroup
	}
)

func newSchema(name string, definition definition) (*Schema, error) {
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
	var textPattern *regroup.ReGroup
	var test predicate
	switch definition.Format {
	case JSON:
		test = func(s string) bool {
			values := gjson.GetMany(s, predicateFields...)
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
	case Text:
		pattern := predicatePatterns["text"]
		textPattern = regroup.MustCompile(pattern.String())
		test = func(s string) bool {
			return pattern.MatchString(strings.TrimRight(s, "\n"))
		}
	default:
		return nil, errors2.Errorf("unknown schema format: %s", definition.Format)
	}
	var inputFields []string
	inverseMapping := map[string]string{}
	for k, v := range definition.Mapping {
		inverseMapping[v] = k
		inputFields = append(inputFields, v)
	}
	var schema Schema = baseSchema{
		name:              name,
		format:            definition.Format,
		priority:          definition.Priority,
		mapping:           inverseMapping,
		inputFields:       inputFields,
		predicate:         test,
		predicatePatterns: predicatePatterns,
		textPattern:       textPattern,
	}
	return &schema, nil
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
