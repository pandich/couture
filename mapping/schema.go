package mapping

import (
	"github.com/oriser/regroup"
	"regexp"
)

// Column ...
type Column string

const (
	// Timestamp ...
	Timestamp Column = "timestamp"
	// Level ...
	Level Column = "level"
	// Message ...
	Message Column = "message"
	// Application ...
	Application Column = "application"
	// Action function, method, etc.
	Action Column = "action"
	// Line ...
	Line Column = "line"
	// Context thread, session id, or some other execution context.
	Context Column = "context"
	// Entity class, struct, etc.
	Entity Column = "entity"
	// Error ...
	Error Column = "error"
)

// Names ...
func Names() []string {
	var names []string
	for _, v := range []Column{
		Timestamp,
		Level,
		Message,
		Application,
		Action,
		Line,
		Context,
		Entity,
		Error,
	} {
		names = append(names, string(v))
	}
	return names
}

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
