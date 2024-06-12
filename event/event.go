// Package event provides a generalized construct for representing input from
// a source.Source.
package event

import (
	"fmt"
	"github.com/pandich/couture/event/level"
	"github.com/pandich/couture/mapping"
	"github.com/pandich/couture/model"
	"github.com/rcrowley/go-metrics"
	"math"
)

type (
	// Event is a generalized log event.
	Event struct {
		// Timestamp the timestamp. This field is required, and should default to time.Now() if not present.
		Timestamp Timestamp `json:"timestamp" regroup:"timestamp"`
		// Level the level. This field is required, and should default to Info if not present.
		Level level.Level `json:"level" regroup:"level"`
		// Message the message. This field is required.
		Message Message `json:"message" regroup:"message"`
		// Application is the name of the application that generated this event. This field is optional.
		Application Application `json:"application" regroup:"application"`
		// Action the action name. This field is optional.
		Action Action `json:"action" regroup:"action"`
		// Line the line number. This field is optional.
		Line Line `json:"line" regroup:"line"`
		// Context the context name. This field is optional.
		Context Context `json:"context" regroup:"context"`
		// Entity the entity name. This field is optional.
		Entity Entity `json:"entity" regroup:"entity"`
		// Error the error. This field is optional.
		Error Error `json:"error" regroup:"error"`
	}

	// SinkEvent wraps an Event with additional application metadata.
	SinkEvent struct {
		Event
		SourceURL SourceURL
		Filters   model.Filters
		Mapping   *mapping.Mapping
	}

	// CodeLocation is where the log event happened.
	CodeLocation string

	// Bucket is the
	Bucket uint8
)

// CodeLocation returns the location where the event happened.
func (event Event) CodeLocation() CodeLocation {
	return CodeLocation(
		fmt.Sprintf(
			"%s.%s.%d",
			event.Entity,
			event.Action,
			event.Line,
		),
	)
}

// Mark manages a per-code-location meter. These meters can be used to generate
// metrics, as well as rate limit specific locations.
func (cl CodeLocation) Mark(category string) {
	cl.meter(category).Mark(1)
}

// meter returns the metrics meter for tje spcofied category.
func (cl CodeLocation) meter(category string) metrics.Meter {
	meterName := fmt.Sprintf("%s.%s.meter", cl, category)
	return metrics.GetOrRegister(meterName, metrics.NewMeter()).(metrics.Meter)
}

// eventsPerMinute of the code location.
func (cl CodeLocation) eventsPerMinute(lvl level.Level) float64 {
	const secondsPerMinute = 60.0
	meter := cl.meter(string(lvl))
	eventsPersSecond := meter.Rate1()
	return eventsPersSecond * secondsPerMinute
}

const (
	// Bucket1 is histogram Bucket 1.
	Bucket1 = iota
	// Bucket2 is histogram Bucket 2.
	Bucket2
	// Bucket3 is histogram Bucket 3.
	Bucket3
	// Bucket4 is histogram Bucket 4.
	Bucket4
	// Bucket5 is histogram Bucket 5.
	Bucket5
	// Bucket6 is histogram Bucket 6.
	Bucket6
	// Bucket7 is histogram Bucket 7.
	Bucket7
	// Bucket8 is histogram Bucket 8.
	Bucket8

	// BucketUpperBound the highest Bucket.
	BucketUpperBound = iota - 1
)

// LevelMeterBucket returns the
func (event SinkEvent) LevelMeterBucket() Bucket {
	eventsPerMinute := event.Event.
		CodeLocation().
		eventsPerMinute(event.Level)

	val := Bucket(math.Log2(eventsPerMinute))
	if val > BucketUpperBound {
		val = BucketUpperBound
	}

	return val
}
