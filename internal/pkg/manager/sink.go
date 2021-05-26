package manager

import (
	"couture/internal/pkg/sink"
)

func sinkWriter(sinks []*sink.Sink) chan sink.Event {
	writer := make(chan sink.Event)
	go func() {
		for {
			event := <-writer
			for _, snk := range sinks {
				err := (*snk).Accept(event)
				if err != nil {
					panic(err)
				}
			}
		}
	}()
	return writer
}
