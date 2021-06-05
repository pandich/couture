package schema

import (
	"github.com/oriser/regroup"
	"regexp"
)

const textRootField = "_"

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
	canHandle func(s string) bool

	priority = uint8

	format string

	// Schema ...
	Schema struct {
		Name                     string            `yaml:"-"`
		Format                   format            `yaml:"format,omitempty"`
		Priority                 priority          `yaml:"priority,omitempty"`
		PredicatesByField        map[string]string `yaml:"predicates,omitempty"`
		FieldByColumn            map[string]string `yaml:"mapping,omitempty"`
		TemplateByColumn         map[string]string `yaml:"display,omitempty"`
		Fields                   []string
		canHandle                canHandle
		predicatePatternsByField map[string]*regexp.Regexp
		predicateFields          []string
		TextPattern              *regroup.ReGroup
	}
)
