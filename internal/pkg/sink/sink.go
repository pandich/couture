package sink

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/source"
	"github.com/rcrowley/go-metrics"
	"io"
)

// Sink of events. Responsible for consuming an event.
type Sink interface {
	// Init called prior to the beginning of logging.
	Init(sources []*source.Source)
	// Accept consumes an event, typically for display.
	Accept(event model.SinkEvent) error
}

// NewOut ...
func NewOut(name string, writer io.WriteCloser) chan string {
	const bufferSize = 32 * 1_024
	eventMeter := metrics.NewMeter()
	metrics.GetOrRegister(name+".outChan.events", eventMeter)
	byteMeter := metrics.NewMeter()
	metrics.GetOrRegister(name+".outChan.bytes", byteMeter)

	out := make(chan string)
	go func() {
		defer close(out)
		writer := bufio.NewWriterSize(writer, bufferSize)
		for {
			message := <-out
			eventMeter.Mark(1)
			byteMeter.Mark(int64(len(message)))
			_, err := writer.WriteString(message + "\n")
			if err != nil {
				panic(err)
			}
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
		}
	}()

	return out
}
