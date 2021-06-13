package cmd

import (
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/pandich/couture/couture"
	"github.com/pandich/couture/manager"
	"github.com/pandich/couture/model"
	"github.com/pandich/couture/schema"
	"github.com/pandich/couture/sink"
	"github.com/pandich/couture/sink/doric"
	"github.com/pandich/couture/sink/doric/column"
	"github.com/pandich/couture/theme"
	"gopkg.in/multierror.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var managerConfig = manager.Config{}

// Run runs the manager using the CLI arguments.
func Run() {
	var err error

	options := parseInputs()
	managerConfig.Schemas, err = schema.LoadSchemas()
	parser.FatalIfErrorf(err)

	mgr, err := manager.New(managerConfig, options...)
	parser.FatalIfErrorf(err)

	err = (*mgr).Run()
	parser.FatalIfErrorf(err)
}

func parseInputs() []interface{} {
	var args = os.Args[1:]
	args, err := expandAliases(args)
	parser.FatalIfErrorf(err)

	_, err = parser.Parse(args)
	parser.FatalIfErrorf(err)

	err = loadSinkConfig()
	parser.FatalIfErrorf(err)

	options, err := getOptions()
	parser.FatalIfErrorf(err)

	return options
}

func getOptions() ([]interface{}, error) {
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

	if len(doricConfig.Columns) == 0 {
		var defaultColumnNames []string
		for i := range column.DefaultColumns {
			defaultColumnNames = append(defaultColumnNames, string(column.DefaultColumns[i]))
		}
		doricConfig.Columns = defaultColumnNames
	}
	if doricConfig.TimeFormat == nil {
		tf := time.Stamp
		doricConfig.TimeFormat = &tf
	}

	th, err := theme.GenerateTheme(string(cli.Theme), string(cli.SourceStyle))
	parser.FatalIfErrorf(err)
	doricConfig.Theme = th
	var options = []interface{}{
		doric.New(doricConfig),
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

func loadConfig() (*sink.Config, error) {
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
