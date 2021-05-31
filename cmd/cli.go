package cmd

import (
	"couture/internal/pkg/model/level"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// FEATURE shell completion
// FEATURE use the kong config integration for all the non-alias config

//nolint:lll
var cli struct {
	OutputFormat     string `group:"display" hidden:"true" help:"The output format: ${enum}." enum:"${outputFormats}" default:"${defaultOutputFormat}" placeholder:"format" short:"f" required:"true" env:"COUTURE_DEFAULT_FORMAT"`
	Wrap             bool   `group:"display" help:"Wrap the output tp the terminal width, or that specified by --width." short:"w" default:"true" negatable:"true"`
	Width            uint   `group:"display" help:"Wrap width." placeholder:"width" short:"W"`
	Theme            string `group:"display" help:"Specify the core Theme color: ${enum}." placeholder:"Theme" default:"${defaultTheme}" enum:"${themeNames}" env:"COUTURE_THEME"`
	Multiline        bool   `group:"display" help:"Display each log event in multiline format." negatable:"true" default:"false"`
	Highlight        bool   `group:"display" help:"Highlight matches from the patterns specified in --include." negatable:"true" default:"true"`
	AutoResize       bool   `group:"display" help:"Auto-resize columns when the terminal resizes." negatable:"true" default:"true"`
	ConsistentColors bool   `group:"display" help:"Maintain consistent source URL colors between runs." negatable:"true" default:"true"`

	Column     []string   `group:"content" help:"Specify one or more columns to display: ${enum}." placeholder:"column" enum:"${columnNames}" env:"COUTURE_DEFAULT_COLUMN_NAMES"`
	TimeFormat timeFormat `group:"content" help:"Go-standard time format string or a named format: ${timeFormatNames}." short:"t" default:"stamp" env:"COUTURE_DEFAULT_TIME_FORMAT"`
	ExpandJSON bool       `group:"content" help:"Example JSON message bodies. Warning: has a significant performance impact." negatable:"true" default:"false"`

	Level   level.Level     `group:"filter" help:"The minimum log level to display: ${enum}." default:"${defaultLogLevel}" placeholder:"level" short:"l" enum:"${logLevels}" env:"COUTURE_DEFAULT_LEVEL"`
	Since   time.Time       `group:"filter" help:"How far back to look for events. Parses most time and duration formats including human friendly." placeholder:"(time|duration)" short:"s" default:"15m" env:"COUTURE_DEFAULT_SINCE"`
	Include []regexp.Regexp `group:"filter" help:"Include filter regular expressions; they are performed before excludes." placeholder:"regex" short:"i" sep:"|"`
	Exclude []regexp.Regexp `group:"filter" help:"Exclude filter regular expressions; they are performed after includes." placeholder:"regex" short:"x" sep:"|"`

	Source []url.URL `arg:"true" help:"Log event source URLs." name:"source_url" required:"true"`
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
