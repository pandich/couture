package schema

import (
	"github.com/oriser/regroup"
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
	// Action function, method, etc.
	Action = "action"
	// Line ...
	Line = "line"
	// Context thread, session id, or some other execution context.
	Context = "context"
	// Entity class, struct, etc.
	Entity = "entity"
	// Error ...
	Error = "error"
)

const (
	// JSON ...
	JSON = "json"
	// Text ...
	Text = "text"
)

// Schema ...
type Schema struct {
	Name                     string            `yaml:"-"`
	Format                   string            `yaml:"format,omitempty"`
	Priority                 uint8             `yaml:"priority,omitempty"`
	PredicatesByField        map[string]string `yaml:"predicates,omitempty"`
	FieldByColumn            map[string]string `yaml:"mapping,omitempty"`
	TemplateByColumn         map[string]string `yaml:"display,omitempty"`
	Fields                   []string
	TextPattern              *regroup.ReGroup
	canHandle                func(s string) bool
	predicatePatternsByField map[string]*regexp.Regexp
	predicateFields          []string
}
