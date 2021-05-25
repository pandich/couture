package cli

import (
	"github.com/aymerick/raymond"
	errors2 "github.com/pkg/errors"
	"net/url"
	"os"
	"regexp"
)

var simpleArgs = regexp.MustCompile(`@(?P<name>\w+)`)

func expandAliases() ([]string, error) {
	args := os.Args[1:]
	for i := range args {
		var arg = args[i]
		arg = expandSchemeShortcut(arg)
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

func expandSchemeShortcut(arg string) string {
	if len(arg) > 0 && arg[0] == '@' {
		arg = "alias://" + arg[1:]
	}
	return arg
}

func expandAlias(aliasURL *url.URL) (string, error) {
	aliases := aliasConfig()
	alias, ok := aliases[aliasURL.Host]
	if !ok {
		return "", errors2.Errorf("unknown alias: %s", aliasURL.Host)
	}

	// expand simple value placeholders
	alias = simpleArgs.ReplaceAllString(alias,
		"{{#if ${name}}}"+
			"${name}={{${name}.[0]}}"+
			"{{/if}}",
	)

	return raymond.Render(alias, aliasURL.Query())
}
