package manager

import (
	"fmt"
	"github.com/gagglepanda/couture/couture"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/model"
	"github.com/gagglepanda/couture/source"
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

var (
	errChanMeter = metrics.NewMeter()
	snkChanMeter = metrics.NewMeter()
	srcChanMeter = metrics.NewMeter()
)

func init() {
	metrics.GetOrRegister(errChanMeterame, errChanMeter)
	metrics.GetOrRegister(snkChanMeterName, snkChanMeter)
	metrics.GetOrRegister(srcChanMeterName, srcChanMeter)
}

func (mgr *busManager) createChannels() (chan source.Event, chan event.SinkEvent, chan source.Error) {
	errChan := mgr.makeErrChan()
	alertChan := mgr.makeAlertChan(errChan)
	unknownChan := mgr.makeUnknownChan()
	snkChan := mgr.makeSnkChan(errChan)
	srcChan := mgr.makeSrcChan(snkChan, alertChan, unknownChan)
	return srcChan, snkChan, errChan
}

func (mgr *busManager) makeUnknownChan() chan string {
	unknownChan := make(chan string)
	go func() {
		for {
			s := <-unknownChan
			if mgr.config.DumpUnknown {
				_, err := fmt.Fprintln(os.Stderr, s)
				if err != nil {
					panic(err)
				}
			}
		}
	}()
	return unknownChan
}

func (mgr *busManager) makeAlertChan(errChan chan source.Error) chan event.SinkEvent {
	alertChan := make(chan event.SinkEvent)
	go func() {
		for {
			alert := <-alertChan
			title := fmt.Sprintf("%s: %s (%s)", couture.Name, alert.Application, alert.SourceURL.ShortForm())
			message := fmt.Sprintf("[%s] %s", alert.Level, alert.Message)
			if err := notifyOS(title, message); err != nil {
				errChan <- source.Error{SourceURL: alert.SourceURL, Error: err}
			}
		}
	}()
	return alertChan
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

func (mgr *busManager) makeSrcChan(
	snkChan chan event.SinkEvent,
	alertChan chan event.SinkEvent,
	unknownChan chan string,
) chan source.Event {
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
			sch := mapping.GuessMapping(sourceEvent.Event, mgr.config.Mappings...)
			if sch == nil {
				unknownChan <- sourceEvent.Event
			}
			modelEvent := unmarshallEvent(sch, sourceEvent.Event)
			filterKind := mgr.filter(modelEvent)
			modelEvent.AsCodeLocation().Mark(string(modelEvent.Level))
			evt := event.SinkEvent{
				SourceURL: sourceURL,
				Event:     *modelEvent,
				Filters:   mgr.config.Filters,
				Mapping:   sch,
			}
			switch filterKind {
			case model.Exclude:
				// do nothing
			case model.Include:
				snkChan <- evt
			case model.AlertOnce:
				snkChan <- evt
				alertChan <- evt
			}
		}
	}()
	return srcChan
}

func (mgr *busManager) makeSnkChan(errChan chan source.Error) chan event.SinkEvent {
	snkChan := make(chan event.SinkEvent)
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
