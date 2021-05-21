package json

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"encoding/json"
	"fmt"
)

type (
	sourceEvent struct {
		SourceURL string `json:"source_url"`
		model.Event
	}
	jsonSink struct {
	}
)

// New ...
func New() *sink.Sink {
	jsonSink := jsonSink{}
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
	fmt.Println(string(contents))
	return nil
}
