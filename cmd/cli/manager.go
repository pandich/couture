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
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run ...
func Run() {
	const consistentRandomSeed = 0xfeed

	// keep consistent random colors, etc between runs
	rand.Seed(consistentRandomSeed)

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

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(interrupt)

		const stopGracePeriod = 250 * time.Millisecond

		cleanup := func() { termenv.Reset(); os.Exit(0) }

		<-interrupt
		(*mgr).Stop()

		go func() { time.Sleep(stopGracePeriod); cleanup() }()
		(*mgr).Wait()
		cleanup()
	}()
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
