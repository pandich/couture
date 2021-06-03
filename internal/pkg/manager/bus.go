package manager

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/rcrowley/go-metrics"
	"os"
)

var errChanMeter = metrics.NewMeter()
var snkChanMeter = metrics.NewMeter()
var srcChanMeter = metrics.NewMeter()

func init() {
	metrics.GetOrRegister("manager.errChan.in", errChanMeter)
	metrics.GetOrRegister("manager.snkChan.in", snkChanMeter)
	metrics.GetOrRegister("manager.srcChan.in", srcChanMeter)
}

func (mgr *busManager) createChannels() (chan source.Event, chan model.SinkEvent, chan source.Error) {
	errChan := mgr.makeErrChan()
	snkChan := mgr.makeSnkChan(errChan)
	srcChan := mgr.makeSrcChan(errChan, snkChan)
	return srcChan, snkChan, errChan
}

func (mgr *busManager) makeErrChan() chan source.Error {
	errChan := make(chan source.Error)

	go func() {
		defer close(errChan)
		for {
			incoming := <-errChan
			errChanMeter.Mark(1)
			if incoming.Error == nil {
				continue
			}
			var sourceName = incoming.SourceURL.String()
			if sourceName == "" {
				sourceName = couture.Name
			}
			outgoing := errorx.Decorate(incoming.Error, "source: %s", sourceName)
			_, err := fmt.Fprintf(os.Stderr, "\nError: %+v\n", outgoing)
			if err != nil {
				panic(err)
			}
		}
	}()

	return errChan
}

func (mgr *busManager) makeSrcChan(_ chan source.Error, snkChan chan model.SinkEvent) chan source.Event {
	srcChan := make(chan source.Event)
	go func() {
		defer close(srcChan)
		for {
			sourceEvent := <-srcChan
			srcChanMeter.Mark(1)
			sch := schema.Guess(sourceEvent.Event, mgr.config.Schemas...)
			modelEvent := unmarshallEvent(sch, sourceEvent.Event)
			if mgr.shouldInclude(modelEvent) {
				snkChan <- model.SinkEvent{
					SourceURL: sourceEvent.Source.URL(),
					Event:     *modelEvent,
					Filters:   mgr.config.Filters,
				}
			}
		}
	}()
	return srcChan
}

func (mgr *busManager) makeSnkChan(errChan chan source.Error) chan model.SinkEvent {
	snkChan := make(chan model.SinkEvent)
	go func() {
		defer close(snkChan)
		for {
			event := <-snkChan
			snkChanMeter.Mark(1)
			for _, snk := range mgr.sinks {
				err := (*snk).Accept(event)
				if err != nil {
					errChan <- source.Error{SourceURL: event.SourceURL, Error: err}
				}
			}
		}
	}()
	return snkChan
}
