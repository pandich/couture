package json

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
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

// Init ...
func (snk jsonSink) Init(_ []source.Source) {
}

// Accept ...
func (snk jsonSink) Accept(src source.Source, event model.Event) error {
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
