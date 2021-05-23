package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	errors2 "github.com/pkg/errors"
	"strings"
	"time"
)

func managerOptionFlags() ([]interface{}, error) {
	snk, err := sinkFlag()
	if err != nil {
		return nil, err
	}
	return []interface{}{
		manager.LogLevelOption(cli.Level),
		manager.FilterOption(cli.Include, cli.Exclude),
		manager.SinceOption(cli.Since),
		snk,
	}, nil
}

func sinkFlag() (interface{}, error) {
	var columns = cli.Column
	if len(columns) == 0 {
		columns = []string{
			"timestamp",
			"application",
			"thread",
			"caller",
			"level",
			"message",
			"error",
		}
	}
	switch cli.OutputFormat {
	case "pretty":
		return pretty.New(config.Config{
			ClearScreen: cli.ClearScreen,
			Columns:     columns,
			MultiLine:   cli.MultiLine,
			ShowSigils:  cli.Sigil,
			Theme:       theme.Registry[cli.Theme],
			TimeFormat:  string(cli.TimeFormat),
			Width:       cli.Width,
			Wrap:        cli.Wrap,
		}), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s\n", cli.OutputFormat)
	}
}

type timeFormat string

// AfterApply ...
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (t *timeFormat) AfterApply() error {
	switch strings.ToLower(string(*t)) {
	case "c":
		*t = time.ANSIC
	case "unix":
		*t = time.UnixDate
	case "ruby":
		*t = time.RubyDate
	case "rfc822":
		*t = time.RFC822
	case "rfc822-utc":
		*t = time.RFC822Z
	case "rfc850":
		*t = time.RFC850
	case "rfc1123":
		*t = time.RFC1123
	case "rfc1123-utc":
		*t = time.RFC1123Z
	case "rfc3339", "iso8601":
		*t = time.RFC3339
	case "rfc3339-nanos", "iso8601-nanos":
		*t = time.RFC3339Nano
	case "kitchen":
		*t = time.Kitchen
	case "stamp":
		*t = time.Stamp
	case "stamp-millis":
		*t = time.StampMilli
	case "stamp-micros":
		*t = time.StampMicro
	case "stamp-nanos":
		*t = time.StampNano
	}
	return nil
}
