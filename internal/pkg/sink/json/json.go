package json

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"encoding/json"
	"fmt"
	"io"
)

type (
	sourceEvent struct {
		SourceURL string `json:"source_url"`
		model.Event
	}
	jsonSink struct {
		out io.Writer
	}
)

// New ...
func New(out io.Writer) *sink.Sink {
	jsonSink := jsonSink{out: out}
	var snk sink.Sink = jsonSink
	return &snk
}

// Accept ...
func (snk jsonSink) Accept(src source.Pushable, event model.Event) error {
	sourceEvent := sourceEvent{
		SourceURL: src.URL().String(),
		Event:     event,
	}
	contents, err := json.Marshal(sourceEvent)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(snk.out, string(contents))
	if err != nil {
		return err
	}
	return nil
}
