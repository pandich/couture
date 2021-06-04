package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/ratelimit"
	"os"
)

const (
	errChanMeterame  = "manager.errChan.in"
	srcChanMeterName = "manager.srcChan.in"
	snkChanMeterName = "manager.snkChan.in"
)

var errChanMeter = metrics.NewMeter()
var snkChanMeter = metrics.NewMeter()
var srcChanMeter = metrics.NewMeter()

func init() {
	metrics.GetOrRegister(errChanMeterame, errChanMeter)
	metrics.GetOrRegister(snkChanMeterName, snkChanMeter)
	metrics.GetOrRegister(srcChanMeterName, srcChanMeter)
}

func (mgr *busManager) createChannels() (chan source.Event, chan model.SinkEvent, chan source.Error) {
	errChan := mgr.makeErrChan()
	snkChan := mgr.makeSnkChan(errChan)
	srcChan := mgr.makeSrcChan(snkChan)
	return srcChan, snkChan, errChan
}

func (mgr *busManager) makeErrChan() chan source.Error {
	errChan := make(chan source.Error)

	go func() {
		defer close(errChan)
		for {
			incoming := <-errChan
			sourceName := incoming.SourceURL.String()
			var errChanSrcMeter = metrics.GetOrRegister(
				errChanMeterame+"."+sourceName,
				metrics.NewMeter(),
			).(metrics.Meter)
			errChanSrcMeter.Mark(1)
			errChanMeter.Mark(1)
			if incoming.Error == nil {
				continue
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

func (mgr *busManager) makeSrcChan(snkChan chan model.SinkEvent) chan source.Event {
	srcChan := make(chan source.Event)
	go func() {
		defer close(srcChan)
		for {
			sourceEvent := <-srcChan
			srcChanMeter.Mark(1)
			sourceURL := sourceEvent.Source.URL()
			var srcChanSrcMeter = metrics.GetOrRegister(
				srcChanMeterName+"."+sourceURL.String(),
				metrics.NewMeter(),
			).(metrics.Meter)
			srcChanSrcMeter.Mark(1)
			sch := schema.Guess(sourceEvent.Event, mgr.config.Schemas...)
			modelEvent := unmarshallEvent(sch, sourceEvent.Event)
			if mgr.shouldInclude(modelEvent) {
				snkChan <- model.SinkEvent{
					SourceURL: sourceURL,
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
	var limiter ratelimit.Limiter
	if mgr.config.RateLimit == 0 {
		limiter = ratelimit.NewUnlimited()
	} else {
		limiter = ratelimit.New(int(mgr.config.RateLimit))
	}
	go func() {
		defer close(snkChan)
		for {
			event := <-snkChan
			limiter.Take()
			snkChanMeter.Mark(1)
			for _, snk := range mgr.sinks {
				if err := (*snk).Accept(event); err != nil {
					errChan <- source.Error{SourceURL: event.SourceURL, Error: err}
				}
			}
		}
	}()
	return snkChan
}
