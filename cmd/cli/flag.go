package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	errors2 "github.com/pkg/errors"
)

func flags() ([]interface{}, error) {
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
	case "pretty":
		return pretty.New(config.Config{
			ClearScreen: cli.ClearScreen,
			Columns:     cli.Column,
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
