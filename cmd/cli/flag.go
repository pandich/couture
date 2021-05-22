package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
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
	switch cli.OutputFormat {
	case "json":
		return json.New(), nil
	case "pretty":
		return pretty.New(pretty.Config{
			Wrap:       cli.Wrap,
			Width:      cli.Width,
			MultiLine:  cli.MultiLine,
			Theme:      themeFlag(),
			Columns:    columnNamesFlag(),
			TimeFormat: timeFormatFlag(),
		}), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s\n", cli.OutputFormat)
	}
}

func timeFormatFlag() string {
	var timeFormat = cli.TimeFormat
	switch strings.ToLower(timeFormat) {
	case "c":
		timeFormat = time.ANSIC
	case "unix":
		timeFormat = time.UnixDate
	case "ruby":
		timeFormat = time.RubyDate
	case "rfc822":
		timeFormat = time.RFC822
	case "rfc822-utc":
		timeFormat = time.RFC822Z
	case "rfc850":
		timeFormat = time.RFC850
	case "rfc1123":
		timeFormat = time.RFC1123
	case "rfc1123-utc":
		timeFormat = time.RFC1123Z
	case "rfc3339", "iso8601":
		timeFormat = time.RFC3339
	case "rfc3339-nanos", "iso8601-nanos":
		timeFormat = time.RFC3339Nano
	case "kitchen":
		timeFormat = time.Kitchen
	case "stamp":
		timeFormat = time.Stamp
	case "stamp-millis":
		timeFormat = time.StampMilli
	case "stamp-micros":
		timeFormat = time.StampMicro
	case "stamp-nanos":
		timeFormat = time.StampNano
	}
	return timeFormat
}

func themeFlag() pretty.Theme {
	return pretty.ThemeByName[cli.Theme]
}

func columnNamesFlag() []pretty.ColumnName {
	var columnNames []pretty.ColumnName
	for _, n := range cli.Column {
		columnNames = append(columnNames, pretty.ColumnName(n))
	}
	return columnNames
}
