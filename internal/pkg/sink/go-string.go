package sink

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
)

//NewGoString provides a configured GoString sink.
func NewGoString(options Options, _ string) interface{} {
	return GoString{baseSink{options: options}}
}

//GoString uses the GoStringer interface to display.
type GoString struct {
	baseSink
}

func (sink GoString) Accept(src source.Source, event model.Event) {
	var line = fmt.Sprintf("%s %+v", src, event)
	if sink.options.Wrap() > 0 {
		line = wordwrap.WrapString(line, sink.options.Wrap())
	}
	fmt.Println(line)
}
