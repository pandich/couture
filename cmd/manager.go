package cmd

import (
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/pandich/couture/internal/pkg/couture"
	"github.com/pandich/couture/internal/pkg/manager"
	"github.com/pandich/couture/internal/pkg/model"
	"github.com/pandich/couture/internal/pkg/schema"
	"github.com/pandich/couture/internal/pkg/sink"
	"github.com/pandich/couture/internal/pkg/sink/doric"
	"github.com/pandich/couture/internal/pkg/sink/doric/column"
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

	managerConfig.Since = cli.Since

	options := parseInputs()
	managerConfig.Schemas, err = schema.LoadSchemas()
	maybeDie(err)

	mgr, err := manager.New(managerConfig, options...)
	maybeDie(err)

	maybeDie((*mgr).Run())
}

func parseInputs() []interface{} {
	var args = os.Args[1:]
	args, err := expandAliases(args)
	maybeDie(err)

	_, err = parser.Parse(args)
	maybeDie(err)

	err = loadSinkConfig()
	maybeDie(err)

	options, err := getOptions()
	maybeDie(err)

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
		doricConfig.Columns = column.DefaultColumns
	}
	if doricConfig.TimeFormat == nil {
		tf := time.Stamp
		doricConfig.TimeFormat = &tf
	}

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
