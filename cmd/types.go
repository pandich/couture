package cmd

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var timeFormatNames = []string{
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

type (
	autoResize       bool
	color            bool
	columns          []string
	consistentColors bool
	expand           bool
	levelLike        level.Level
	highlight        bool
	multiline        bool
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
)

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v autoResize) AfterApply() error { prettyConfig.AutoResize = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v color) AfterApply() error { prettyConfig.Color = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v columns) AfterApply() error { prettyConfig.Columns = v; return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v consistentColors) AfterApply() error { prettyConfig.ConsistentColors = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v expand) AfterApply() error { prettyConfig.Expand = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v highlight) AfterApply() error { prettyConfig.Highlight = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v multiline) AfterApply() error { prettyConfig.Multiline = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v tty) AfterApply() error { prettyConfig.TTY = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v width) AfterApply() error { prettyConfig.Width = uint(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v wrap) AfterApply() error { prettyConfig.Wrap = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v dumpMetrics) AfterApply() error { managerConfig.DumpMetrics = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v dumpUnknown) AfterApply() error { managerConfig.DumpUnknown = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v showSchema) AfterApply() error { prettyConfig.ShowSchema = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v rateLimit) AfterApply() error { managerConfig.RateLimit = uint(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v levelLike) AfterApply() error { managerConfig.Level = level.Level(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (f filterLike) AfterApply() (err error) {
	managerConfig.Filters, err = f.asFilters()
	return
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v themeName) AfterApply() error {
	thm, ok := theme.Registry[string(v)]
	if !ok {
		return errors2.Errorf("unknown theme: %s", v)
	}
	prettyConfig.Theme = thm
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (t timeFormat) AfterApply() error {
	format := strings.ToLower(string(t))
	switch format {
	case "c":
		prettyConfig.TimeFormat = time.ANSIC
	case "unix":
		prettyConfig.TimeFormat = time.UnixDate
	case "ruby":
		prettyConfig.TimeFormat = time.RubyDate
	case "rfc822":
		prettyConfig.TimeFormat = time.RFC822
	case "rfc822-utc":
		prettyConfig.TimeFormat = time.RFC822Z
	case "rfc850":
		prettyConfig.TimeFormat = time.RFC850
	case "rfc1123":
		prettyConfig.TimeFormat = time.RFC1123
	case "rfc1123-utc":
		prettyConfig.TimeFormat = time.RFC1123Z
	case "rfc3339", "iso8601":
		prettyConfig.TimeFormat = time.RFC3339
	case "rfc3339-nanos", "iso8601-nanos":
		prettyConfig.TimeFormat = time.RFC3339Nano
	case "kitchen":
		prettyConfig.TimeFormat = time.Kitchen
	case "stamp":
		prettyConfig.TimeFormat = time.Stamp
	case "stamp-millis":
		prettyConfig.TimeFormat = time.StampMilli
	case "stamp-micros":
		prettyConfig.TimeFormat = time.StampMicro
	case "stamp-nanos":
		prettyConfig.TimeFormat = time.StampNano
	default:
		prettyConfig.TimeFormat = "stamp"
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
