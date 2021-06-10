package cmd

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/doric"
	"couture/internal/pkg/sink/doric/config"
	"github.com/coreos/etcd/pkg/fileutil"
	"gopkg.in/multierror.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
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

func loadConfig() (*config.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filename := path.Join(home, ".config", couture.Name, "config.yaml")
	if !fileutil.Exist(filename) {
		return &config.Config{}, nil
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
	var c config.Config
	err = yaml.Unmarshal(text, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
