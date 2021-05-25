package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"github.com/aymerick/raymond"
	"gopkg.in/multierror.v1"
	"net/url"
	"os"
)

func sourceArgs() ([]interface{}, error) {
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

func evaluatedArgs() ([]string, error) {
	args := os.Args[1:]
	for i := range args {
		arg := args[i]
		if len(arg) > 0 && arg[0] == '@' {
			arg = "alias://" + arg[1:]
		}
		aliasURL, err := url.Parse(arg)
		if err == nil && aliasURL.Scheme == "alias" {
			value, err := expandAlias(aliasURL)
			if err != nil {
				return nil, err
			}
			if value != "" {
				args[i] = value
			}
		}
	}
	return args, nil
}

func expandAlias(aliasURL *url.URL) (string, error) {
	aliases := aliasConfig()
	alias, ok := aliases[aliasURL.Host]
	if !ok {
		return alias, nil
	}
	render, err := raymond.Render(alias, aliasURL.Query())
	return render, err
}
