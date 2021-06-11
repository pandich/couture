package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/theme"
	errors2 "github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var timeFormatNames = []string{
	model.HumanTimeFormat,
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
)

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *autoResize) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.AutoResize = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *color) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.Color = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v columns) AfterApply() error { doricConfig.Columns = v; return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *consistentColors) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.ConsistentColors = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *expand) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.Expand = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *highlight) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.Highlight = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *multiLine) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.MultiLine = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v tty) AfterApply() error { doricConfig.TTY = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *width) AfterApply() error {
	if v == nil {
		return nil
	}
	ui := uint(*v)
	doricConfig.Width = &ui
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *wrap) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.Wrap = &b
	return nil
}

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v dumpMetrics) AfterApply() error { managerConfig.DumpMetrics = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v dumpUnknown) AfterApply() error { managerConfig.DumpUnknown = bool(v); return nil }

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (v *showSchema) AfterApply() error {
	if v == nil {
		return nil
	}
	b := bool(*v)
	doricConfig.ShowSchema = &b
	return nil
}

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
func (v *themeName) AfterApply() error {
	if v == nil {
		return nil
	}
	thm, ok := theme.Registry[string(*v)]
	if !ok {
		return errors2.Errorf("unknown theme: %s", *v)
	}
	doricConfig.Theme = &thm
	return nil
}

// AfterApply ...
//nolint:funlen
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (t *timeFormat) AfterApply() error {
	if t == nil {
		return nil
	}
	format := strings.ToLower(string(*t))
	switch format {
	case model.HumanTimeFormat:
		doricConfig.TimeFormat = &format
	case "c":
		s := time.ANSIC
		doricConfig.TimeFormat = &s
	case "unix":
		s := time.UnixDate
		doricConfig.TimeFormat = &s
	case "ruby":
		s := time.RubyDate
		doricConfig.TimeFormat = &s
	case "rfc822":
		s := time.RFC822
		doricConfig.TimeFormat = &s
	case "rfc822-utc":
		s := time.RFC822Z
		doricConfig.TimeFormat = &s
	case "rfc850":
		s := time.RFC850
		doricConfig.TimeFormat = &s
	case "rfc1123":
		s := time.RFC1123
		doricConfig.TimeFormat = &s
	case "rfc1123-utc":
		s := time.RFC1123Z
		doricConfig.TimeFormat = &s
	case "rfc3339", "iso8601":
		s := time.RFC3339
		doricConfig.TimeFormat = &s
	case "rfc3339-nanos", "iso8601-nanos":
		s := time.RFC3339Nano
		doricConfig.TimeFormat = &s
	case "kitchen":
		s := time.Kitchen
		doricConfig.TimeFormat = &s
	case "stamp":
		s := time.Stamp
		doricConfig.TimeFormat = &s
	case "stamp-millis":
		s := time.StampMilli
		doricConfig.TimeFormat = &s
	case "stamp-micros":
		s := time.StampMicro
		doricConfig.TimeFormat = &s
	case "stamp-nanos":
		s := time.StampNano
		doricConfig.TimeFormat = &s
	default:
		s := "stamp"
		doricConfig.TimeFormat = &s
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
