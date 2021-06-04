package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"github.com/araddon/dateparse"
	"github.com/tidwall/gjson"
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
	values := gjson.GetMany(s, (*sch).Fields()...)
	event := model.Event{}
	for i, field := range (*sch).Fields() {
		if col, ok := (*sch).Column(field); ok {
			value := values[i]
			updateEvent(&event, col, value)
		}
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
		Timestamp:   model.Timestamp(time.Now()),
		Level:       level.Warn,
		Message:     model.Message(msg),
		Application: "",
		Method:      "",
		Line:        0,
		Thread:      "",
		Class:       "",
		Exception:   "Warning: entry is in an unknown log format.",
	}
}

func updateEvent(event *model.Event, col string, v gjson.Result) {
	switch col {
	case schema.Timestamp:
		if v.Exists() {
			t, _ := dateparse.ParseAny(v.String())
			event.Timestamp = model.Timestamp(t)
		}
	case schema.Level:
		const defaultLevel = level.Info
		if v.Exists() {
			event.Level = level.ByName(v.String(), defaultLevel)
		} else {
			event.Level = defaultLevel
		}
	case schema.Message:
		if v.Exists() {
			event.Message = model.Message(v.String())
		}
	case schema.Application:
		if v.Exists() {
			event.Application = model.Application(v.String())
		}
	case schema.Method:
		if v.Exists() {
			event.Method = model.Method(v.String())
		}
	case schema.Line:
		if v.Exists() {
			event.Line = model.Line(v.Int())
		}
	case schema.Thread:
		if v.Exists() {
			event.Thread = model.Thread(v.String())
		}
	case schema.Class:
		if v.Exists() {
			event.Class = model.Class(v.String())
		}
	case schema.Exception:
		if v.Exists() {
			stackTrace := v.String()
			event.Exception = model.Exception(stackTrace)
		}
	}
}
