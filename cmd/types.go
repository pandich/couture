package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/event/level"
	"github.com/gagglepanda/couture/model"
	errors2 "github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

import (
	"github.com/gagglepanda/couture/sink"
)

// custom type declarations to add functionakity to kong.

const (
	colorModeAuto  colorMode = "auto"
	colorModeDark  colorMode = "dark"
	colorModeLight colorMode = "light"
)

var (
	sinkConfig = sink.DefaultConfig()

	timeFormatNames = []string{
		event.HumanTimeFormat,
		"c",
		"iso8601",
		"iso8601-nanos",
		"kitchen",
		"rfc1123",
		"rfc1123-utc",
		"rfc3339",
		"rfc3339-nanos",
		"rfc822",
		"rfc822-utc",
		"rfc850",
		"ruby",
		"stamp",
		"stamp-micros",
		"stamp-millis",
		"stamp-nanos",
		"unix",
	}
)

type (
	autoResize       bool
	noColor          bool
	sourceStyle      string
	columns          []string
	consistentColors bool
	expand           bool
	levelLike        level.Level
	highlight        bool
	levelMeter       bool
	multiLine        bool
	themeName        string
	timeFormat       string
	tty              bool
	width            uint
	wrap             bool
	dumpMetrics      bool
	dumpUnknown      bool
	showMapping      bool
	rateLimit        uint
	filterLike       []string
	colorMode        string
)

// sink configuration:

// AfterApply provides additional functionality to CLI processor.
func (v *autoResize) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.AutoResize = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v *noColor) AfterApply() error {
	if v == nil {
		return nil
	}
	b := !bool(*v)
	sinkConfig.Color = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v columns) AfterApply() error { sinkConfig.Columns = v; return nil }

// AfterApply provides additional functionality to CLI processor.
func (v *consistentColors) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.ConsistentColors = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v *expand) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.Expand = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v *highlight) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.Highlight = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v *multiLine) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.MultiLine = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v *levelMeter) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.LevelMeter = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v tty) AfterApply() error { sinkConfig.TTY = bool(v); return nil }

// AfterApply provides additional functionality to CLI processor.
func (v *width) AfterApply() error {
	if v == nil {
		return nil
	}
	ui := uint(*v)
	sinkConfig.Width = &ui
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v *wrap) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.Wrap = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (t *timeFormat) AfterApply() error {
	if t == nil {
		return nil
	}
	format := strings.ToLower(string(*t))
	switch format {
	case event.HumanTimeFormat:
		sinkConfig.TimeFormat = &format
	case "c":
		s := time.ANSIC
		sinkConfig.TimeFormat = &s
	case "unix":
		s := time.UnixDate
		sinkConfig.TimeFormat = &s
	case "ruby":
		s := time.RubyDate
		sinkConfig.TimeFormat = &s
	case "rfc822":
		s := time.RFC822
		sinkConfig.TimeFormat = &s
	case "rfc822-utc":
		s := time.RFC822Z
		sinkConfig.TimeFormat = &s
	case "rfc850":
		s := time.RFC850
		sinkConfig.TimeFormat = &s
	case "rfc1123":
		s := time.RFC1123
		sinkConfig.TimeFormat = &s
	case "rfc1123-utc":
		s := time.RFC1123Z
		sinkConfig.TimeFormat = &s
	case "rfc3339", "iso8601":
		s := time.RFC3339
		sinkConfig.TimeFormat = &s
	case "rfc3339-nanos", "iso8601-nanos":
		s := time.RFC3339Nano
		sinkConfig.TimeFormat = &s
	case "kitchen":
		s := time.Kitchen
		sinkConfig.TimeFormat = &s
	case "stamp":
		s := time.Stamp
		sinkConfig.TimeFormat = &s
	case "stamp-millis":
		s := time.StampMilli
		sinkConfig.TimeFormat = &s
	case "stamp-micros":
		s := time.StampMicro
		sinkConfig.TimeFormat = &s
	case "stamp-nanos":
		s := time.StampNano
		sinkConfig.TimeFormat = &s
	default:
		s := "stamp"
		sinkConfig.TimeFormat = &s
	}
	return nil
}

// logging manager configuration:

// AfterApply provides additional functionality to CLI processor.
func (v *showMapping) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.ShowMapping = &b
	return nil
}

// AfterApply provides additional functionality to CLI processor.
func (v dumpMetrics) AfterApply() error { mgrCfg.DumpMetrics = bool(v); return nil }

// AfterApply provides additional functionality to CLI processor.
func (v dumpUnknown) AfterApply() error { mgrCfg.DumpUnknown = bool(v); return nil }

// AfterApply provides additional functionality to CLI processor.
func (v rateLimit) AfterApply() error { mgrCfg.RateLimit = uint(v); return nil }

// AfterApply provides additional functionality to CLI processor.
func (v levelLike) AfterApply() error { mgrCfg.Level = level.Level(v); return nil }

// AfterApply provides additional functionality to CLI processor.
func (f filterLike) AfterApply() (err error) {
	mgrCfg.Filters, err = f.asFilters()
	return
}

// helpers

// timeLikeDecoder create a kong function to provide flexible time/duration parsing.
func timeLikeDecoder() kong.MapperFunc {
	startupTime := time.Now()

	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var t time.Time

		var value string

		// get the value
		if err := ctx.Scan.PopValueInto("(time|duration)", &value); err != nil {
			return err
		}

		// if this is a valid duration
		d, err := time.ParseDuration(value)
		if err == nil {
			// subtract it from startup time
			t = startupTime.Add(-d)
		} else {
			// otherwise try to parse it as a datetime
			t, err = dateparse.ParseAny(value)
			if err != nil {
				return errors2.Errorf("expected time or duration but got %q: %s", value, err)
			}
		}

		// set the value
		target.Set(reflect.ValueOf(&t))

		// success
		return nil
	}
}

// asFilters converts filterLike values into model.Filter values. filters allow
// inclusion/exclusion of lines, as well a fire-once alert feature.
func (f filterLike) asFilters() ([]model.Filter, error) {
	const (
		// alert indicates that the source definition should stop after the first event received.
		alert = '@'
		// include the events matching this filter
		include = '+'
		// exclude the events matching this filter
		exclude = '-'
	)

	var filters []model.Filter

	// convert a string pattern into a regex filter and add it to the list of filters
	addPattern := func(pattern string, kind model.FilterKind) error {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		filters = append(filters, model.Filter{Pattern: *re, Kind: kind})
		return nil
	}

	var value string
	for i := range f {
		value = f[i]
		flag := value[0]

		switch flag {
		case alert:
			if err := addPattern(value[1:], model.AlertOnce); err != nil {
				return nil, err
			}

		case include:
			if err := addPattern(value[1:], model.Include); err != nil {
				return nil, err
			}

		case exclude:
			if err := addPattern(value[1:], model.Exclude); err != nil {
				return nil, err
			}

		default: //  implicit include
			if err := addPattern(value, model.Include); err != nil {
				return nil, err
			}
		}
	}

	return filters, nil
}
