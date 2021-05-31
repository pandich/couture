package cmd

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	errors2 "github.com/pkg/errors"
	"gopkg.in/multierror.v1"
	"os"
)

// Run runs the manager using the CLI arguments.
func Run() {
	var args = os.Args[1:]
	args, err := expandAliases(args)
	maybeDie(err)

	_, err = parser.Parse(args)
	maybeDie(err)

	schemas, err := schema.LoadSchemas()
	maybeDie(err)
	cfg := manager.Config{
		Level:          cli.Level,
		Since:          &cli.Since,
		IncludeFilters: cli.Include,
		ExcludeFilters: cli.Exclude,
		Schemas:        schemas,
	}

	options, err := getOptions()
	maybeDie(err)

	mgr, err := manager.New(cfg, options...)
	maybeDie(err)

	err = (*mgr).Run()
	maybeDie(err)
}

func getOptions() ([]interface{}, error) {
	sourceArgs := func() ([]interface{}, error) {
		var sources []interface{}
		var violations []error
		for _, u := range cli.Source {
			sourceURL := model.SourceURL(u)
			src, err := manager.GetSource(sourceURL)
			if len(err) > 0 {
				violations = append(violations, err...)
			} else {
				for _, s := range src {
					sources = append(sources, s)
				}
			}
		}
		if len(violations) > 0 {
			return nil, multierror.New(violations)
		}
		return sources, nil
	}

	sinkFlag := func() (interface{}, error) {
		switch cli.OutputFormat {
		case "pretty":
			return pretty.New(config.Config{
				AutoResize:       cli.AutoResize,
				Columns:          cli.Column,
				ConsistentColors: cli.ConsistentColors,
				Highlight:        cli.Highlight,
				Multiline:        cli.Multiline,
				Theme:            theme.Registry[cli.Theme],
				TimeFormat:       string(cli.TimeFormat),
				Width:            cli.Width,
				Wrap:             cli.Wrap,
			}), nil
		default:
			return nil, errors2.Errorf("unknown output format: %s\n", cli.OutputFormat)
		}
	}

	var options []interface{}
	snk, err := sinkFlag()
	if err != nil {
		return nil, err
	}
	options = append(options, snk)

	sources, err := sourceArgs()
	if err != nil {
		return nil, err
	}
	options = append(options, sources...)

	return options, nil
}
