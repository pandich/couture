package cmd

import (
	"couture/internal/pkg/couture"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"couture/internal/pkg/schema"
	"couture/internal/pkg/sink/pretty"
	"couture/internal/pkg/sink/pretty/config"
	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/pkg/fileutil"
	"gopkg.in/multierror.v1"
	"io/ioutil"
	"os"
	"path"
)

var prettyConfig = config.Config{
	Out: os.Stdout,
}

// Run runs the manager using the CLI arguments.
func Run() {
	var args = os.Args[1:]
	args, err := expandAliases(args)
	maybeDie(err)

	err = loadConfig()
	maybeDie(err)

	_, err = parser.Parse(args)
	maybeDie(err)

	schemas, err := schema.LoadSchemas()
	maybeDie(err)
	cfg := manager.Config{
		DumpMetrics: cli.Metrics,
		Level:       cli.Level,
		Since:       cli.Since,
		Filters:     cli.Filter,
		Schemas:     schemas,
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
		pretty.New(prettyConfig),
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

func loadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	aliasFilename := path.Join(home, ".config", couture.Name, "config.toml")
	if !fileutil.Exist(aliasFilename) {
		return nil
	}

	f, err := os.Open(aliasFilename)
	if err != nil {
		return err
	}

	defer f.Close()
	text, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return toml.Unmarshal(text, &prettyConfig)
}
