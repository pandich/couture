package sink

import (
	"couture/pkg/couture/model"
)

type (

	// Sink of events. Responsible for consuming an event.
	Sink interface {
		// ConsumeEvent consumes an event, typically for display.
		ConsumeEvent(event *model.Event)
	}
)

/*

Implementations go here. Each implementation struct should be unexported and exposed with a var.

Example:

	var Stdout Sink = stdoutSink{}

	type stdoutSink struct {}

	func (s stdoutSink) ConsumeEvent(event *model.Event) { fmt.Printf("%#v\n", *event) }

*/
