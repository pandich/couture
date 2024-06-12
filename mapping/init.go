package mapping

import (
	"fmt"
	"github.com/oriser/regroup"
	"regexp"
)

func (mapping *Mapping) init(name string) {
	mapping.Name = name
	mapping.initPredicatePatterns()
	mapping.initCanHandle()
	mapping.initFields()
}

func (mapping *Mapping) initPredicatePatterns() {
	mapping.predicatePatternsByField = map[string]*regexp.Regexp{}
	for field, pattern := range mapping.PredicatesByField {
		mapping.predicateFields = append(mapping.predicateFields, field)
		if pattern != "" {
			mapping.predicatePatternsByField[field] = regexp.MustCompile(pattern)
		} else {
			mapping.predicatePatternsByField[field] = nil
		}
	}
}

func (mapping *Mapping) initCanHandle() {
	switch mapping.Format {
	case JSON:
		mapping.canHandle = mapping.canHandleJSON
	case Text:
		var pattern = mapping.predicatePatternsByField[textRootField].String()
		re := regexp.MustCompile(pattern)
		names := map[Column]bool{
			Timestamp:   false,
			Level:       false,
			Message:     false,
			Application: false,
			Action:      false,
			Line:        false,
			Context:     false,
			Entity:      false,
			Error:       false,
		}
		for i := 0; i < re.NumSubexp(); i++ {
			names[Column(re.SubexpNames()[i+1])] = true
		}

		for name, ok := range names {
			if !ok {
				pattern += fmt.Sprintf("(?P<%s>)", name)
			}
		}
		mapping.TextPattern = regroup.MustCompile(pattern)
		mapping.canHandle = mapping.canHandleText
	}
}

func (mapping *Mapping) initFields() {
	for _, field := range mapping.FieldByColumn {
		mapping.Fields = append(mapping.Fields, field)
	}
}
