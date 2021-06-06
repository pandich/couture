package manager

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/joomcode/errorx"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/ratelimit"
	"os"
	"time"
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
	alertChan := mgr.makeAlertChan(errChan)
	snkChan := mgr.makeSnkChan(errChan)
	srcChan := mgr.makeSrcChan(snkChan, alertChan)
	return srcChan, snkChan, errChan
}

func (mgr *busManager) makeAlertChan(errChan chan source.Error) chan model.SinkEvent {
	alertChan := make(chan model.SinkEvent)
	go func() {
		const maxNotificationsPerMinute = 10
		const noIcon = ""

		limiter := ratelimit.New(maxNotificationsPerMinute, ratelimit.Per(time.Minute))

		for {
			alert := <-alertChan
			limiter.Take()
			title := fmt.Sprintf("%s: %s (%s)", couture.Name, alert.Application, alert.SourceURL.ShortForm())
			message := fmt.Sprintf("[%s] %s", alert.Level, alert.Message)
			if err := beeep.Notify(title, message, noIcon); err != nil {
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

func (mgr *busManager) makeSrcChan(snkChan chan model.SinkEvent, alertChan chan model.SinkEvent) chan source.Event {
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
			filterKind := mgr.filter(modelEvent)
			evt := model.SinkEvent{
				SourceURL: sourceURL,
				Event:     *modelEvent,
				Filters:   mgr.config.Filters,
				Schema:    sch,
			}
			switch filterKind {
			case model.Exclude:
				// do nothing
			case model.Include:
				snkChan <- evt
			case model.Alert:
				alertChan <- evt
				snkChan <- evt
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
