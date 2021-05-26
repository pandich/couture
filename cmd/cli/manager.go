package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"github.com/muesli/termenv"
	errors2 "github.com/pkg/errors"
	"gopkg.in/multierror.v1"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run ...
func Run() {
	// load config
	err := loadAliasConfig()
	parser.FatalIfErrorf(err)

	// expand aliases, etc.
	args, err := expandAliases()
	parser.FatalIfErrorf(err)

	_, err = parser.Parse(args)
	parser.FatalIfErrorf(err)

	// get cli managerOptions and args
	mgrOptions, err := managerOptions()
	parser.FatalIfErrorf(err)

	// create the manager and start it
	mgr, err := manager.New(mgrOptions...)
	parser.FatalIfErrorf(err)

	err = (*mgr).Start()
	parser.FatalIfErrorf(err)

	trapSignals(mgr)

	(*mgr).Wait()
}

func managerOptions() ([]interface{}, error) {
	sourceArgs := func() ([]interface{}, error) {
		var sources []interface{}
		var violations []error
		for _, u := range cli.Source {
			sourceURL := model.SourceURL(u)
			src, err := manager.GetSource(sourceURL)
			if len(err) > 0 {
				violations = append(violations, err...)
			} else {
				sources = append(sources, src...)
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
				Columns:    cli.Column,
				Multiline:  cli.Multiline,
				Theme:      theme.Registry[cli.Theme],
				TimeFormat: string(cli.TimeFormat),
				Width:      cli.Width,
				Wrap:       cli.Wrap,
			}), nil
		default:
			return nil, errors2.Errorf("unknown output format: %s\n", cli.OutputFormat)
		}
	}

	var options = []interface{}{
		manager.LogLevelOption(cli.Level),
		manager.FilterOption(cli.Include, cli.Exclude),
		manager.SinceOption(cli.Since),
	}

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

func trapSignals(mgr *model.Manager) {
	const stopGracePeriod = 2 * time.Second

	die := func() {
		println()
		termenv.Reset()
		os.Exit(0)
	}

	interruption := make(chan os.Signal)
	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruption
		(*mgr).Stop()
		go func() { time.Sleep(stopGracePeriod); die() }()
		(*mgr).Wait()
		die()
	}()
}
