package sink

import (
	"couture/pkg/couture/model"
	"log"
)

var (
	Diagnostic Sink = diagnosticSink{}
)

type diagnosticSink struct {
}

func (l diagnosticSink) ConsumeEvent(event *model.Event) {
	log.Printf("%#v\n", *event)
}
