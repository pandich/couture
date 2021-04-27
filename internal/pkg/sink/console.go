package sink

import (
	"couture/pkg/couture/model"
	"fmt"
)

var Console Sink = consoleSink{}

type consoleSink struct {
}

func (l consoleSink) ConsumeEvent(event *model.Event) {
	fmt.Printf("%#v\n", *event)
}
