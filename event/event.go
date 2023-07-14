package event

import (
	"fmt"
	"github.com/gagglepanda/couture/event/level"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/model"
	"github.com/rcrowley/go-metrics"
	"math"
)

// Event a log event
type Event struct {
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

// SinkEvent ...
type SinkEvent struct {
	Event
	SourceURL SourceURL
	Filters   model.Filters
	Mapping   *mapping.Mapping
}

// CodeLocation ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
type CodeLocation string

// AsCodeLocation ...
func (event Event) AsCodeLocation() CodeLocation {
	return CodeLocation(fmt.Sprintf(
		"%s.%s.%d",
		event.Entity,
		event.Action,
		event.Line,
	))
}

// Mark ...
func (cl CodeLocation) Mark(category string) {
	cl.meter(category).Mark(1)
}

func (cl CodeLocation) meter(category string) metrics.Meter {
	meterName := fmt.Sprintf("%s.%s.meter", cl, category)
	return metrics.GetOrRegister(meterName, metrics.NewMeter()).(metrics.Meter)
}

func (cl CodeLocation) eventsPerMinute(lvl level.Level) float64 {
	const secondsPerMinute = 60.0
	meter := cl.meter(string(lvl))
	eventsPersSecond := meter.Rate1()
	return eventsPersSecond * secondsPerMinute
}

// BucketMax ...
const BucketMax = Bucket8

const (
	// Bucket1 ...
	Bucket1 = iota
	// Bucket2 ...
	Bucket2
	// Bucket3 ...
	Bucket3
	// Bucket4 ...
	Bucket4
	// Bucket5 ...
	Bucket5
	// Bucket6 ...
	Bucket6
	// Bucket7 ...
	Bucket7
	// Bucket8 ...
	Bucket8
)

// LevelMeterBucket ...
func (event SinkEvent) LevelMeterBucket() uint8 {
	eventsPerMinute := event.Event.AsCodeLocation().eventsPerMinute(event.Level)
	var bucket = uint8(math.Log2(eventsPerMinute))
	if bucket > BucketMax {
		bucket = BucketMax
	}
	return bucket
}
