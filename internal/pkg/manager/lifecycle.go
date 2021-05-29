package manager

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/model/schema"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/joomcode/errorx"
	errors2 "github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"os"
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

//nolint:funlen
func unmarshallEvent(sch schema.Schema, s string) (*model.Event, error) {
	if !gjson.Valid(s) {
		return nil, errors2.Errorf("invalid JSON: %s", s)
	}
	values := gjson.GetMany(s, sch.Fields...)
	event := model.Event{}
	for i := 0; i < len(sch.Fields); i++ {
		v := values[i]
		f := schema.Field(sch.Fields[i])
		c := sch.Mapping[f]
		switch c {
		case schema.TimestampCol:
			if v.Exists() {
				t, _ := dateparse.ParseAny(v.String())
				event.Timestamp = model.Timestamp(t)
			}
		case schema.LevelCol:
			const defaultLevel = level.Info
			if v.Exists() {
				event.Level = level.ByName(v.String(), defaultLevel)
			} else {
				event.Level = defaultLevel
			}
		case schema.MessageCol:
			if v.Exists() {
				event.Message = model.Message(model.PrettyJSON(v.String()))
			}
		case schema.ApplicationCol:
			if v.Exists() {
				name := model.ApplicationName(v.String())
				event.ApplicationName = &name
			}
		case schema.MethodCol:
			if v.Exists() {
				event.MethodName = model.MethodName(v.String())
			}
		case schema.LineCol:
			if v.Exists() {
				event.LineNumber = model.LineNumber(v.Int())
			}
		case schema.ThreadCol:
			if v.Exists() {
				name := model.ThreadName(v.String())
				event.ThreadName = &name
			}
		case schema.ClassCol:
			if v.Exists() {
				event.ClassName = model.ClassName(v.String())
			}
		case schema.ExceptionCol:
			if v.Exists() {
				stackTrace := model.PrettyJSON(v.String())
				event.Exception = &model.Exception{StackTrace: model.StackTrace(stackTrace)}
			}
		}
	}
	return &event, nil
}

func (mgr *publishingManager) makeErrChan() chan source.Error {
	errChan := make(chan source.Error)

	go func() {
		defer close(errChan)
		for {
			incoming := <-errChan
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

func (mgr *publishingManager) makeSrcChan(errChan chan source.Error, snkChan chan model.SinkEvent) chan source.Event {
	srcChan := make(chan source.Event)
	go func() {
		defer close(srcChan)
		for {
			sourceEvent := <-srcChan
			sch, ok := mgr.config.Schemas[sourceEvent.Schema]
			if !ok {
				errChan <- source.Error{
					SourceURL: sourceEvent.Source.URL(),
					Error:     errors2.Errorf("unknown schema: %s", sourceEvent.Schema),
				}
			} else {
				modelEvent, err := unmarshallEvent(sch, sourceEvent.Event)
				if err != nil {
					errChan <- source.Error{
						SourceURL: sourceEvent.Source.URL(),
						Error:     err,
					}
				} else if mgr.shouldInclude(*modelEvent) {
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
