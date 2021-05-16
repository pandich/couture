package json

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"os"
)

// TODO need to fix marshalling of model objects

var (
	lexer     = lexers.Get("json")
	formatter = formatters.TTY256
	style     = styles.Dracula
)

type (
	sourceEvent struct {
		Source string
		Event  model.Event
	}
	jsonSink struct{}
)

// New ...
func New() *sink.Sink {
	jsonSink := jsonSink{}
	var snk sink.Sink = jsonSink
	return &snk
}

// Accept ...
func (j jsonSink) Accept(src source.Source, event model.Event) {
	sourceEvent := sourceEvent{
		Source: src.URL().String(),
		Event:  event,
	}
	contents, err := json.MarshalIndent(sourceEvent, "", "  ")
	if err != nil {
		// TODO handle this
		return
	}
	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		// TODO handle this
		return
	}
	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		// TODO handle this
		return
	}
	fmt.Println("")
}
