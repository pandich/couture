package cli

import (
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	errors2 "github.com/pkg/errors"
)

var defaultColumns = []string{
	"timestamp",
	"application",
	"thread",
	"caller",
	"level",
	"message",
	"error",
}

func sinkFlag() (interface{}, error) {
	var columns = cli.Column
	if len(columns) == 0 {
		columns = defaultColumns
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
