// Package cmd is the entry point for the application.
// It is responsible for parsing command line arguments and launching the application.
package cmd

import (
	"fmt"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gagglepanda/couture/couture"
	"github.com/gagglepanda/couture/event"
	"github.com/gagglepanda/couture/manager"
	"github.com/gagglepanda/couture/mapping"
	"github.com/gagglepanda/couture/sink"
	"github.com/gagglepanda/couture/sink/color"
	"github.com/gagglepanda/couture/sink/layout"
	"github.com/gagglepanda/couture/sink/layout/doric"
	"github.com/gagglepanda/couture/sink/layout/doric/column"
	theme2 "github.com/gagglepanda/couture/sink/theme"
	"github.com/pkg/errors"
	"gopkg.in/multierror.v1"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path"
	"time"
)

var (
	mgrCfg  = manager.Config{}
	sinkCfg = sink.Config{
		AutoResize:       &sink.Enabled,
		Color:            &sink.Enabled,
		ConsistentColors: &sink.Enabled,
		Expand:           &sink.Disabled,
		Highlight:        &sink.Disabled,
		MultiLine:        &sink.Disabled,
		ShowMapping:      &sink.Disabled,
		Wrap:             &sink.Disabled,
		Layout:           &layout.Default,
		Out:              os.Stdout,
		Theme:            nil,
		TimeFormat:       &sink.DefaultTimeFormat,
	}
)

// Run runs the manager using the CLI arguments.
func Run() {
	var err error

	// expand any arguments into a new list of arguments
	aliasCfg, err := loadAliasConfig()
	parser.FatalIfErrorf(err)
	args, err := aliasCfg.expandAliases(os.Args[1:])
	parser.FatalIfErrorf(err)

	// parse
	_, err = parser.Parse(args)
	parser.FatalIfErrorf(err)

	// light/dark mode
	switch cli.ColorMode {
	case colorModeDark:
		color.Mode = color.DarkMode
	case colorModeLight:
		color.Mode = color.LightMode
	case colorModeAuto:
		// leave unchanged
	}

	// layer any user sink configuration over the defaults
	sinkConfig.
		PopulateMissing(loadSinkConfigFile()).
		PopulateMissing(sinkCfg)

	// parse the arguments
	options, err := parseOptions()
	parser.FatalIfErrorf(err)

	// load the (optional) mappings from the user's config file
	mgrCfg.Mappings, err = mapping.LoadMappings()
	parser.FatalIfErrorf(err)

	// instantiate the app
	mgr, err := manager.New(mgrCfg, options...)
	parser.FatalIfErrorf(err)

	// start
	err = (*mgr).Run()
	parser.FatalIfErrorf(err)
}

// parseOptions
func parseOptions() ([]interface{}, error) {
	if len(sinkConfig.Columns) == 0 {
		var defaultColumnNames []string
		for i := range column.DefaultColumns {
			defaultColumnNames = append(defaultColumnNames, string(column.DefaultColumns[i]))
		}
		sinkConfig.Columns = defaultColumnNames
	}

	if sinkConfig.TimeFormat == nil {
		tf := time.Stamp
		sinkConfig.TimeFormat = &tf
	}

	th, err := theme2.GenerateTheme(string(cli.Theme))
	parser.FatalIfErrorf(err)
	sinkConfig.Theme = th
	var options = []interface{}{
		doric.New(sinkConfig),
	}

	// attempt to set up each logging source specified in the cli args
	var sources []interface{}
	var violations []error
	for _, u := range cli.Source {
		sourceURL := event.SourceURL(u)
		src, err := manager.GetSource(cli.Since, sourceURL)
		if len(err) > 0 {
			violations = append(violations, err...)
		} else {
			for _, s := range src {
				sources = append(sources, s)
			}
		}
	}

	// fail on any violation
	if len(violations) > 0 {
		return nil, multierror.New(violations)
	}

	// otherwise add the sources to the options
	options = append(options, sources...)

	return options, nil
}

func loadSinkConfigFile() sink.Config {
	// try to load the config file
	tryLoad := func() (*sink.Config, error) {
		// get the home dir
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		// create the filename
		filename := path.Join(home, ".config", couture.Name, "config.yaml")

		// the file is optional: if it doesn't exist provide an empty config
		if !fileutil.Exist(filename) {
			return &sink.Config{}, nil
		}

		// try to open it
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// read the config file
		text, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		// unmarshall config
		var c sink.Config
		if err = yaml.Unmarshal(text, &c); err != nil {
			return nil, err
		}

		return &c, nil
	}

	// try to load. as no output sinks are defined at this point, errors must go to STDERR directly.
	cfg, err := tryLoad()
	if err != nil {
		// use a default
		cfg = &sink.Config{}

		_, _ = fmt.Fprintf(
			os.Stderr, "%s\n",
			errors.Wrapf(err, "error processing configuration file"),
		)
	}

	// successfully loaded an existing config
	return *cfg
}
