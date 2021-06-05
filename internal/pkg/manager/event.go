package manager

import (
	"bytes"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/tidwall/gjson"
	"html/template"
	"os"
	"strings"
	"time"
)

func unmarshallEvent(sch *schema.Schema, s string) *model.Event {
	var evt *model.Event
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
		const envKey = "COUTURE_DIE_ON_UNKNOWN"
		const exitCode = 12
		if os.Getenv(envKey) != "" {
			fmt.Printf("unknown: %+v\n", s)
			os.Exit(exitCode)
		}
		evt = unmarshallUnknown(s)
	}
	return evt
}

func unmarshallJSONEvent(sch *schema.Schema, s string) *model.Event {
	values := map[string]gjson.Result{}
	for i, value := range gjson.GetMany(s, sch.Fields...) {
		field := sch.Fields[i]
		values[field] = value
	}

	event := model.Event{}
	for col, field := range sch.FieldByColumn {
		updateEvent(
			&event,
			col,
			field,
			values,
			sch.TemplateByColumn[col],
		)
	}
	return &event
}

func unmarshallTextEvent(sch *schema.Schema, s string) *model.Event {
	if sch.TextPattern == nil {
		return nil
	}

	event := model.Event{}
	err := sch.TextPattern.MatchToTarget(strings.TrimRight(s, "\n"), &event)
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

func updateEvent(event *model.Event, col string, field string, values map[string]gjson.Result, tmpl string) {
	rawValue := values[field]
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
