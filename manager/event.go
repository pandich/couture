package manager

import (
	"bytes"
	"github.com/araddon/dateparse"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/event/level"
	"github.com/gagglepanda/couture/schema"
	"github.com/tidwall/gjson"
	"html/template"
	"strings"
	"time"
)

func unmarshallEvent(sch *schema.Schema, s string) *event.Event {
	var evt *event.Event
	if sch != nil {
		switch sch.Format {
		case schema.JSON:
			evt = unmarshallJSONEvent(sch, s)
		case schema.Text:
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

func unmarshallJSONEvent(sch *schema.Schema, s string) *event.Event {
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

func unmarshallTextEvent(sch *schema.Schema, s string) *event.Event {
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
	switch schema.Column(col) {
	case schema.Timestamp:
		s := value
		if s != "" {
			t, _ := dateparse.ParseAny(s)
			evt.Timestamp = event.Timestamp(t)
		}
	case schema.Application:
		evt.Application = event.Application(value)
	case schema.Context:
		evt.Context = event.Context(value)
	case schema.Entity:
		evt.Entity = event.Entity(value)
	case schema.Action:
		evt.Action = event.Action(value)
	case schema.Line:
		if rawValue.Exists() {
			evt.Line = event.Line(rawValue.Int())
		}
	case schema.Level:
		s := value
		const defaultLevel = level.Info
		if s != "" {
			evt.Level = level.ByName(s, defaultLevel)
		} else {
			evt.Level = defaultLevel
		}
	case schema.Message:
		evt.Message = event.Message(value)
	case schema.Error:
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
