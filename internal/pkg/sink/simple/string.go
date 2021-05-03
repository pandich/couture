package simple

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"couture/pkg/model"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
)

// Sink uses the fmt.GoStringer interface to display.
type Sink struct {
	sink.Base
	full bool
}

// New ...
func New(options sink.Options, config string) interface{} {
	return Sink{Base: sink.New(options), full: config == "full"}
}

// Accept ...
func (sink Sink) Accept(src source.Source, event model.Event) {
	var line string
	if sink.full {
		line = fmt.Sprintf("%s %+v", src, event)
	} else {
		line = fmt.Sprintf("%s %#v", src, event)
	}
	if sink.Options().Wrap() > 0 {
		line = wordwrap.WrapString(line, sink.Options().Wrap())
	}
	fmt.Println(line)
}
