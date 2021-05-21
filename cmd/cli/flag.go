package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
	"github.com/muesli/gamut"
	errors2 "github.com/pkg/errors"
)

const (
	blackAndWhite = ""
	purpleRain    = "#ae99bf"
	merlot        = "#a01010"
	ocean         = "#5198eb"
)

// TODO theme color tweaks
var themeByName = map[string]pretty.Theme{
	"none":     {BaseColor: blackAndWhite, SourceColors: gamut.PastelGenerator{}},
	"prince":   {BaseColor: purpleRain, SourceColors: gamut.PastelGenerator{}},
	"brougham": {BaseColor: merlot, SourceColors: gamut.WarmGenerator{}},
	"ocean":    {BaseColor: ocean, SourceColors: gamut.HappyGenerator{}},
}

func getFlags() ([]interface{}, error) {
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
	var columnNames []pretty.ColumnName
	for _, n := range cli.Column {
		columnNames = append(columnNames, pretty.ColumnName(n))
	}

	switch cli.OutputFormat {
	case "json":
		return json.New(), nil
	case "pretty":
		theme, err := themeFlag()
		if err != nil {
			return nil, err
		}
		return pretty.New(pretty.Config{
			Wrap:       cli.Wrap,
			Width:      cli.Width,
			MultiLine:  cli.MultiLine,
			Theme:      theme,
			Columns:    columnNames,
			TimeFormat: cli.TimeFormat,
		}), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s\n", cli.OutputFormat)
	}
}

func themeFlag() (pretty.Theme, error) {
	t, ok := themeByName[cli.Theme]
	if !ok {
		return t, errors2.Errorf("unknown Theme: %s", cli.Theme)
	}
	return t, nil
}
