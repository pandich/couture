package cli

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/alecthomas/kong"
	"strings"
)

func parserVars() kong.Vars {
	return kong.Vars{
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
		}, ", "),
		"columnNames":         strings.Join(column.Names(), ","),
		"themeNames":          strings.Join(theme.Names(), ","),
		"defaultTheme":        theme.Prince,
		"logLevels":           strings.Join(level.SimpleNames(), ","),
		"defaultLogLevel":     level.Info.SimpleName(),
		"outputFormats":       strings.Join([]string{pretty.Name}, ","),
		"defaultOutputFormat": pretty.Name,
	}
}
