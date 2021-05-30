package cli

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"github.com/araddon/dateparse"
	errors2 "github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const helpSummary = "Tails one or more event sources."

var parser = kong.Must(&cli,
	kong.Name(couture.Name),
	kong.Description(helpDescription()),
	kong.UsageOnError(),
	kong.ConfigureHelp(kong.HelpOptions{
		Summary:   true,
		FlagsLast: true,
	}),
	kong.TypeMapper(reflect.TypeOf(regexp.Regexp{}), regexpDecoder()),
	kong.TypeMapper(reflect.TypeOf(time.Time{}), timeLikeDecoder()),
	kong.Vars{
		"timeFormatNames": strings.Join([]string{
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
		}, ","),
		"columnNames":         strings.Join(column.Names(), ","),
		"themeNames":          strings.Join(theme.Names(), ","),
		"defaultTheme":        theme.Default,
		"logLevels":           strings.Join(level.Lower(), ","),
		"defaultLogLevel":     level.Info.Lower(),
		"outputFormats":       strings.Join([]string{pretty.Name}, ","),
		"defaultOutputFormat": pretty.Name,
	},
	kong.Groups{
		"display": "Display Config",
		"content": "Content Config",
		"filter":  "Filter Config",
	},
	kong.PostBuild(completionsHook()),
)

var maybeDie = parser.FatalIfErrorf

func helpDescription() string {
	var lines = []string{
		helpSummary,
		"",
		"Examples Sources:",
		"",
	}
	for _, src := range manager.AvailableSources {
		if len(src.ExampleURLs) > 0 {
			lines = append(lines, "  "+src.Name+":")
			for _, u := range src.ExampleURLs {
				lines = append(lines, "    "+u)
			}
			lines = append(lines, "")
		}
	}
	return strings.Join(lines, "\n")
}

func regexpDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("regex", &value); err != nil {
			return err
		}
		r, err := regexp.Compile(value)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(*r))
		return nil
	}
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
			t = time.Now().Add(-d)
		} else {
			t, err = dateparse.ParseAny(value)
			if err != nil {
				return errors2.Errorf("expected duration but got %q: %s", value, err)
			}
		}
		target.Set(reflect.ValueOf(t))
		return nil
	}
}
