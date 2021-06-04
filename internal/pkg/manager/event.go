package manager

import (
	"bytes"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"github.com/araddon/dateparse"
	"github.com/tidwall/gjson"
	"html/template"
	"strings"
	"time"
)

func unmarshallEvent(sch *schema.Schema, s string) *model.Event {
	var evt *model.Event
	if sch != nil {
		switch (*sch).Format() {
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

func unmarshallJSONEvent(sch *schema.Schema, s string) *model.Event {
	values := map[string]gjson.Result{}
	fields := (*sch).Fields()
	for i, value := range gjson.GetMany(s, (*sch).Fields()...) {
		field := fields[i]
		col, _ := (*sch).Column(field)
		values[col] = value
	}

	event := model.Event{}
	for _, field := range fields {
		col, _ := (*sch).Column(field)
		tmpl, _ := (*sch).Template(col)
		updateEvent(&event, col, values, tmpl)
	}
	return &event
}

func unmarshallTextEvent(sch *schema.Schema, s string) *model.Event {
	pattern := (*sch).TextPattern()
	if pattern == nil {
		return nil
	}

	event := model.Event{}
	err := pattern.MatchToTarget(strings.TrimRight(s, "\n"), &event)
	if err != nil {
		return nil
	}
	return &event
}

func unmarshallUnknown(msg string) *model.Event {
	return &model.Event{
		Timestamp: model.Timestamp(time.Now()),
		Level:     level.Info,
		Message:   model.Message(msg),
	}
}

func updateEvent(event *model.Event, col string, values map[string]gjson.Result, tmpl string) {
	rawValue := values[col]
	value := getValue(tmpl, values, rawValue)
	switch col {
	case schema.Timestamp:
		s := value
		if s != "" {
			t, _ := dateparse.ParseAny(s)
			event.Timestamp = model.Timestamp(t)
		}
	case schema.Application:
		event.Application = model.Application(value)
	case schema.Thread:
		event.Thread = model.Thread(value)
	case schema.Class:
		event.Class = model.Class(value)
	case schema.Method:
		event.Method = model.Method(value)
	case schema.Line:
		if rawValue.Exists() {
			event.Line = model.Line(rawValue.Int())
		}
	case schema.Level:
		s := value
		const defaultLevel = level.Info
		if s != "" {
			event.Level = level.ByName(s, defaultLevel)
		} else {
			event.Level = defaultLevel
		}
	case schema.Message:
		event.Message = model.Message(value)
	case schema.Exception:
		event.Exception = model.Exception(value)
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
