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
	showSchema       bool
	rateLimit        uint
	filterLike       []string
	colorMode        string
)

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *autoResize) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.AutoResize = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *noColor) AfterApply() error {
	if v == nil {
		return nil
	}
	b := !bool(*v)
	sinkConfig.Color = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v columns) AfterApply() error { sinkConfig.Columns = v; return nil }

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *consistentColors) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.ConsistentColors = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *expand) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.Expand = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *highlight) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.Highlight = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *multiLine) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.MultiLine = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *levelMeter) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.LevelMeter = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v tty) AfterApply() error { sinkConfig.TTY = bool(v); return nil }

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *width) AfterApply() error {
	if v == nil {
		return nil
	}
	ui := uint(*v)
	sinkConfig.Width = &ui
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *wrap) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.Wrap = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v dumpMetrics) AfterApply() error { managerConfig.DumpMetrics = bool(v); return nil }

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v dumpUnknown) AfterApply() error { managerConfig.DumpUnknown = bool(v); return nil }

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *showSchema) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	sinkConfig.ShowSchema = &b
	return nil
}

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v rateLimit) AfterApply() error { managerConfig.RateLimit = uint(v); return nil }

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v levelLike) AfterApply() error { managerConfig.Level = level.Level(v); return nil }

// AfterApply ...
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (f filterLike) AfterApply() (err error) {
	managerConfig.Filters, err = f.asFilters()
	return
}

// AfterApply ...
// nolint: funlen
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
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

func timeLikeDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("(time|duration)", &value); err != nil {
			return err
		}
		var t time.Time
		d, err := time.ParseDuration(value)
		if err == nil {
			t = now.Add(-d)
		} else {
			t, err = dateparse.ParseAny(value)
			if err != nil {
				return errors2.Errorf("expected duration but got %q: %s", value, err)
			}
		}
		v := reflect.ValueOf(&t)
		target.Set(v)
		return nil
	}
}

func (f filterLike) asFilters() ([]model.Filter, error) {
	const alert = '@'
	const include = '+'
	const exclude = '-'

	var filters []model.Filter
	addPattern := func(pattern string, kind model.FilterKind) error {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		filters = append(filters, model.Filter{Pattern: *re, Kind: kind})
		return nil
	}

	for _, value := range f {
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
		default:
			if err := addPattern(value, model.Include); err != nil {
				return nil, err
			}
		}
	}

	return filters, nil
}
