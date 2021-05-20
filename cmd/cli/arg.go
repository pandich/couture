package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"gopkg.in/multierror.v1"
	"os"
)

func getArgs() ([]interface{}, error) {
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

func evaluateArgs() []string {
	aliases := aliasConfig()
	args := os.Args[1:]
	for i := range args {
		if alias, ok := aliases[args[i]]; ok {
			args[i] = alias
		}
	}
	return args
}
