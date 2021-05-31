package manager

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/joomcode/errorx"
	"github.com/muesli/termenv"
	"github.com/tidwall/gjson"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Pushable instances.
// Waits until it has been stopped.
func (mgr *publishingManager) Start() error {
	mgr.running = true
	for _, snk := range mgr.sinks {
		(*snk).Init(mgr.sources)
	}
	srcChan, snkChan, errChan := mgr.createChannels()
	for _, src := range mgr.sources {
		mgr.wg.Add(1)
		err := (*src).Start(mgr.wg, func() bool { return mgr.running }, srcChan, snkChan, errChan)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop the Manager. This stops all source.PushingSource instances, and stops polling all polling.Pushable instances.
func (mgr *publishingManager) Stop() {
	mgr.running = false
}

// Wait ...
func (mgr *publishingManager) Wait() {
	mgr.wg.Wait()
}

func (mgr *publishingManager) createChannels() (chan source.Event, chan model.SinkEvent, chan source.Error) {
	errChan := mgr.makeErrChan()
	snkChan := mgr.makeSnkChan(errChan)
	srcChan := mgr.makeSrcChan(errChan, snkChan)
	return srcChan, snkChan, errChan
}

func unmarshallEvent(sch schema.Schema, s string) *model.Event {
	values := gjson.GetMany(s, sch.Fields()...)
	event := model.Event{}
	for i, field := range sch.Fields() {
		if col, ok := sch.Column(field); ok {
			value := values[i]
			updateEvent(&event, col, value)
		}
	}
	return &event
}

func updateEvent(event *model.Event, col string, v gjson.Result) {
	switch col {
	case schema.Timestamp:
		if v.Exists() {
			t, _ := dateparse.ParseAny(v.String())
			event.Timestamp = model.Timestamp(t)
		}
	case schema.Level:
		const defaultLevel = level.Info
		if v.Exists() {
			event.Level = level.ByName(v.String(), defaultLevel)
		} else {
			event.Level = defaultLevel
		}
	case schema.Message:
		if v.Exists() {
			event.Message = model.Message(model.PrettyJSON(v.String()))
		}
	case schema.Application:
		if v.Exists() {
			event.Application = model.Application(v.String())
		}
	case schema.Method:
		if v.Exists() {
			event.Method = model.Method(v.String())
		}
	case schema.Line:
		if v.Exists() {
			event.Line = model.Line(v.Int())
		}
	case schema.Thread:
		if v.Exists() {
			event.Thread = model.Thread(v.String())
		}
	case schema.Class:
		if v.Exists() {
			event.Class = model.Class(v.String())
		}
	case schema.Exception:
		if v.Exists() {
			stackTrace := model.PrettyJSON(v.String())
			event.Exception = model.Exception(stackTrace)
		}
	}
}

func (mgr *publishingManager) makeErrChan() chan source.Error {
	errChan := make(chan source.Error)

	go func() {
		defer close(errChan)
		for {
			incoming := <-errChan
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

func (mgr *publishingManager) makeSrcChan(_ chan source.Error, snkChan chan model.SinkEvent) chan source.Event {
	srcChan := make(chan source.Event)
	go func() {
		defer close(srcChan)
		for {
			sourceEvent := <-srcChan
			sch := schema.Guess(sourceEvent.Event, mgr.config.Schemas...)
			if sch == nil {
				snkChan <- model.SinkEvent{
					Event: model.Event{
						Timestamp:   model.Timestamp(time.Now()),
						Level:       level.Warn,
						Message:     model.Message(sourceEvent.Event),
						Application: couture.Name,
						Method:      "-",
						Line:        0,
						Thread:      "-",
						Class:       "-",
						Exception:   "Warning: entry is in an unknown log format.",
					},
					SourceURL: sourceEvent.Source.URL(),
				}
			} else {
				modelEvent := unmarshallEvent(*sch, sourceEvent.Event)
				if mgr.shouldInclude(modelEvent) {
					snkChan <- model.SinkEvent{
						SourceURL: sourceEvent.Source.URL(),
						Event:     *modelEvent,
						Filters:   mgr.config.IncludeFilters,
					}
				}
			}
		}
	}()
	return srcChan
}

func (mgr *publishingManager) makeSnkChan(errChan chan source.Error) chan model.SinkEvent {
	snkChan := make(chan model.SinkEvent)
	go func() {
		defer close(snkChan)
		for {
			event := <-snkChan
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

// TrapSignals ...
func (mgr *publishingManager) TrapSignals() {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		const stopGracePeriod = 250 * time.Millisecond
		defer close(interrupt)

		cleanup := func() { termenv.Reset(); os.Exit(0) }

		<-interrupt
		(*mgr).Stop()

		go func() { time.Sleep(stopGracePeriod); cleanup() }()
		(*mgr).Wait()
		cleanup()
	}()
}

// Run ...
func (mgr *publishingManager) Run() error {
	mgr.TrapSignals()
	err := mgr.Start()
	if err != nil {
		return err
	}
	// wait for it to die
	(*mgr).Wait()
	return nil
}
