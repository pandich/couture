package manager

import (
	"bytes"
	"github.com/araddon/dateparse"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/event/level"
	"github.com/gagglepanda/couture/mapping"
	"github.com/tidwall/gjson"
	"html/template"
	"strings"
	"time"
)

func unmarshallEvent(sch *mapping.Mapping, s string) *event.Event {
	var evt *event.Event
	if sch != nil {
		switch sch.Format {
		case mapping.JSON:
			evt = unmarshallJSONEvent(sch, s)
		case mapping.Text:
			fallthrough
		default:
			evt = unmarshallTextEvent(sch, s)
		}
	}
	if evt == nil {
		evt = unmarshallUnknown(s)
	}
	return evt
}

func unmarshallJSONEvent(sch *mapping.Mapping, s string) *event.Event {
	values := map[string]gjson.Result{}
	for i, value := range gjson.GetMany(s, sch.Fields...) {
		field := sch.Fields[i]
		values[field] = value
	}

	evt := event.Event{}
	for col, field := range sch.FieldByColumn {
		updateEvent(
			&evt,
			col,
			field,
			values,
			sch.TemplateByColumn[col],
		)
	}
	return &evt
}

func unmarshallTextEvent(sch *mapping.Mapping, s string) *event.Event {
	if sch.TextPattern == nil {
		return nil
	}

	evt := event.Event{}
	err := sch.TextPattern.MatchToTarget(strings.TrimRight(s, "\n"), &evt)
	if err != nil {
		return nil
	}
	return &evt
}

func unmarshallUnknown(msg string) *event.Event {
	return &event.Event{
		Timestamp: event.Timestamp(time.Now()),
		Level:     level.Info,
		Message:   event.Message(msg),
	}
}

func updateEvent(evt *event.Event, col string, field string, values map[string]gjson.Result, tmpl string) {
	rawValue := values[field]
	value := getValue(tmpl, values, rawValue)
	switch mapping.Column(col) {
	case mapping.Timestamp:
		s := value
		if s != "" {
			t, _ := dateparse.ParseAny(s)
			evt.Timestamp = event.Timestamp(t)
		}
	case mapping.Application:
		evt.Application = event.Application(value)
	case mapping.Context:
		evt.Context = event.Context(value)
	case mapping.Entity:
		evt.Entity = event.Entity(value)
	case mapping.Action:
		evt.Action = event.Action(value)
	case mapping.Line:
		if rawValue.Exists() {
			evt.Line = event.Line(rawValue.Int())
		}
	case mapping.Level:
		s := value
		const defaultLevel = level.Info
		if s != "" {
			evt.Level = level.ByName(s, defaultLevel)
		} else {
			evt.Level = defaultLevel
		}
	case mapping.Message:
		evt.Message = event.Message(value)
	case mapping.Error:
		evt.Error = event.Error(value)
	}
}

func getValue(tmpl string, data interface{}, defaultValue gjson.Result) string {
	if tmpl == "" {
		if defaultValue.Exists() {
			return defaultValue.String()
		}
		return ""
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "%%error:parse%%"
	}

	var txt bytes.Buffer
	err = t.Execute(&txt, data)
	if err != nil {
		return "%%error:execute%%"
	}

	return txt.String()
}
