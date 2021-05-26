package manager

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
)

// eventTopic is the topic for all registry and sinks to communicate over.
const eventTopic = "topic:event"

func (mgr *publishingManager) publishEvent(src source.Source, event model.Event) {
	if !event.Level.IsAtLeast(mgr.options.level) {
		return
	}
	if event.Message.Matches(mgr.options.includeFilters, mgr.options.excludeFilters) {
		mgr.bus.Publish(eventTopic, src, sink.Event{
			Event:   event,
			Filters: mgr.options.includeFilters,
		})
	}
}
