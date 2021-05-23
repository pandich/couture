package cli

import (
	"bytes"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/model"
	"gopkg.in/multierror.v1"
	"html/template"
	"net/url"
	"os"
	"regexp"
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
		if len(arg) > 0 || arg[0] == '@' {
			arg = "alias:///" + arg[1:]
		}
		aliasURL, err := url.Parse(arg)
		if err == nil {
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
	var aliasShortFormPattern = regexp.MustCompile(`@(\w+)`)

	aliases := aliasConfig()

	alias, ok := aliases[aliasURL.Path[1:]]
	if !ok {
		// this isn't an alias
		return "", nil
	}

	alias = aliasShortFormPattern.ReplaceAllString(alias, "{{index (index .Query \"$1\") 0}}")
	tmpl, err := template.New("alias").Parse(alias)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, aliasURL)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
