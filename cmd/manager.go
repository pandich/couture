package cmd

import (
	"fmt"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/pandich/couture/couture"
	"github.com/pandich/couture/manager"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink"
	"github.com/pandich/couture/sink/color"
	"github.com/pandich/couture/sink/doric"
	"github.com/pandich/couture/sink/doric/column"
	"github.com/pandich/couture/sink/layout"
	theme2 "github.com/pandich/couture/sink/theme"
	"github.com/pkg/errors"
	"gopkg.in/multierror.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var (
	managerConfig      = manager.Config{}
	defaultDoricConfig = sink.Config{
		AutoResize:       &enabled,
		Color:            &enabled,
		ConsistentColors: &enabled,
		Expand:           &disabled,
		Highlight:        &disabled,
		MultiLine:        &disabled,
		ShowSchema:       &disabled,
		Wrap:             &disabled,
		Layout:           &layout.Default,
		Out:              os.Stdout,
		Theme:            nil,
		TimeFormat:       &defaultTimeFormat,
	}
)

// Run runs the manager using the CLI arguments.
func Run() {
	var err error

	args, err := expandAliases(os.Args[1:])
	parser.FatalIfErrorf(err)

	_, err = parser.Parse(args)
	parser.FatalIfErrorf(err)

	setColorMode()

	cliDoricConfig.
		PopulateMissing(loadDoricConfigFile()).
		PopulateMissing(defaultDoricConfig)

	options, err := parseOptions()
	parser.FatalIfErrorf(err)

	managerConfig.Schemas, err = schema.LoadSchemas()
	parser.FatalIfErrorf(err)

	mgr, err := manager.New(managerConfig, options...)
	parser.FatalIfErrorf(err)

	err = (*mgr).Run()
	parser.FatalIfErrorf(err)
}

func parseOptions() ([]interface{}, error) {
	sourceArgs := func() ([]interface{}, error) {
		var sources []interface{}
		var violations []error
		for _, u := range cli.Source {
			sourceURL := model.SourceURL(u)
			src, err := manager.GetSource(cli.Since, sourceURL)
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

	if len(cliDoricConfig.Columns) == 0 {
		var defaultColumnNames []string
		for i := range column.DefaultColumns {
			defaultColumnNames = append(defaultColumnNames, string(column.DefaultColumns[i]))
		}
		cliDoricConfig.Columns = defaultColumnNames
	}
	if cliDoricConfig.TimeFormat == nil {
		tf := time.Stamp
		cliDoricConfig.TimeFormat = &tf
	}

	th, err := theme2.GenerateTheme(string(cli.Theme))
	parser.FatalIfErrorf(err)
	cliDoricConfig.Theme = th
	var options = []interface{}{
		doric.New(cliDoricConfig),
	}
	sources, err := sourceArgs()
	if err != nil {
		return nil, err
	}
	options = append(options, sources...)
	if len(options) == 0 {
		return nil, nil
	}
	return options, nil
}

func loadDoricConfigFile() sink.Config {
	tryLoad := func() (*sink.Config, error) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		filename := path.Join(home, ".config", couture.Name, "config.yaml")
		if !fileutil.Exist(filename) {
			return &sink.Config{}, nil
		}

		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		defer f.Close()
		text, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		var c sink.Config
		err = yaml.Unmarshal(text, &c)
		if err != nil {
			return nil, err
		}
		return &c, nil
	}
	cfg, err := tryLoad()
	if err != nil {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"%s\n",
			errors.Wrapf(err, "error processing configuration file"),
		)
		return sink.Config{}
	}
	return *cfg
}

func setColorMode() {
	switch cli.ColorMode {
	case "auto":
		break
	case "dark":
		color.Mode = color.DarkMode
	case "light":
		color.Mode = color.LightMode
	}
}
