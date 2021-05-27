package manager

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"fmt"
	errors2 "github.com/pkg/errors"
	"os"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Pushable instances.
// Waits until it has been stopped.
func (mgr *publishingManager) Start() error {
	mgr.running = true
	for _, snk := range mgr.sinks {
		(*snk).Init(mgr.sources)
	}
	srcChan, errChan := mgr.createChannels()
	for _, src := range mgr.sources {
		mgr.wg.Add(1)
		err := (*src).Start(mgr.wg, func() bool { return mgr.running }, srcChan, errChan)
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

func (mgr *publishingManager) createChannels() (chan source.Event, chan source.Error) {
	srcChan := make(chan source.Event)
	snkChan := make(chan sink.Event)
	errChan := make(chan source.Error)

	go func() {
		defer close(srcChan)
		for {
			evt := <-srcChan
			if mgr.shouldInclude(evt) {
				snkChan <- sink.Event{Event: evt, Filters: mgr.options.includeFilters}
			}
		}
	}()

	go func() {
		defer close(errChan)
		for {
			err := <-errChan
			var sourceName = ""
			if err.Source != nil {
				sourceName = err.Source.URL().String()
			}
			// TODO proper error logging
			fmt.Fprintln(os.Stderr, errors2.Errorf(sourceName, err.Error))
		}
	}()

	go func() {
		defer close(snkChan)
		for {
			event := <-snkChan
			for _, snk := range mgr.sinks {
				err := (*snk).Accept(event)
				if err != nil {
					errChan <- source.Error{Source: event.Source, Error: err}
				}
			}
		}
	}()

	return srcChan, errChan
}
