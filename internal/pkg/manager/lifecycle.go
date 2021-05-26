package manager

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
)

// Start the Manager. This starts all source.PushingSource instances, and begins polling all polling.Pushable instances.
// Waits until it has been stopped.
func (mgr *publishingManager) Start() error {
	mgr.running = true

	for _, snk := range mgr.sinks {
		(*snk).Init(mgr.sources)
	}

	out := make(chan source.Event)
	go func() {
		defer close(out)
		for {
			mgr.publishEvent(<-out)
		}
	}()

	for _, src := range mgr.sources {
		mgr.wg.Add(1)
		err := (*src).Start(mgr.wg, func() bool { return mgr.running }, out)
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

func (mgr *publishingManager) publishEvent(evt source.Event) {
	if !evt.Level.IsAtLeast(mgr.options.level) {
		return
	}
	if evt.Message.Matches(mgr.options.includeFilters, mgr.options.excludeFilters) {
		mgr.out <- sink.Event{
			Event:   evt,
			Filters: mgr.options.includeFilters,
		}
	}
}
