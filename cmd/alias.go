package cmd

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gagglepanda/couture/couture"
	errors2 "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"
)

// TODO test alias behavior

var now = time.Now()

const aliasScheme = "alias"

type aliasConfig struct {
	Groups  map[string][]string `yaml:"groups,omitempty"`
	Aliases map[string]string   `yaml:"aliases,omitempty"`
}

func expandAliases(args []string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	aliasFilename := path.Join(home, ".config", couture.Name, "aliases.yaml")
	if !fileutil.Exist(aliasFilename) {
		return args, nil
	}
	aliasFile, err := os.Open(aliasFilename)
	if err != nil {
		return nil, err
	}
	defer aliasFile.Close()
	s, err := io.ReadAll(aliasFile)
	if err != nil {
		return nil, err
	}

	var config aliasConfig
	err = yaml.Unmarshal(s, &config)
	if err != nil {
		return nil, err
	}

	var expandedArgs []string
	for i := range args {
		expandedArgs = append(expandedArgs, expandAliasString(config, args[i])...)
	}
	for i := range expandedArgs {
		if u, err := url.Parse(expandedArgs[i]); err == nil && u.Scheme == "alias" {
			value, err := expandAliasURL(config, u)
			if err != nil {
				return nil, err
			}
			if value != "" {
				expandedArgs[i] = value
			}
		}
	}
	return expandedArgs, nil
}

func expandAliasString(config aliasConfig, arg string) []string {
	const aliasNamePrefix = "@"
	const groupNamePrefix = aliasNamePrefix + aliasNamePrefix
	const aliasURLPrefix = aliasScheme + "://"

	var args []string
	switch {
	case len(arg) > len(groupNamePrefix) && arg[:len(groupNamePrefix)] == groupNamePrefix:
		groupName := arg[len(groupNamePrefix):]
		group, ok := config.Groups[groupName]
		if !ok {
			return nil
		}
		for _, alias := range group {
			args = append(args, aliasURLPrefix+alias)
		}
	case len(arg) > len(aliasNamePrefix) && arg[:len(aliasNamePrefix)] == aliasNamePrefix:
		args = append(args, aliasURLPrefix+arg[len(aliasNamePrefix):])
	default:
		args = append(args, arg)
	}
	return args
}

func expandAliasURL(config aliasConfig, aliasURL *url.URL) (string, error) {
	simpleArgs := regexp.MustCompile(`@(?P<name>\w+)`)
	alias, ok := config.Aliases[aliasURL.Host]
	if !ok {
		return "", errors2.Errorf("unknown alias: %s", aliasURL.Host)
	}
	// expand simple value placeholders
	alias = simpleArgs.ReplaceAllString(
		alias,
		"{{#if ${name}}}"+
			"${name}={{${name}.[0]}}"+
			"{{/if}}",
	)
	return raymond.Render(alias, aliasContext(aliasURL))
}

func aliasContext(aliasURL *url.URL) map[string][]string {
	const century = 100
	context := map[string][]string{
		"epoch": {fmt.Sprintf("%d", now.Unix())},
		"yyyy":  {fmt.Sprintf("%04d", now.Year())},
		"yy":    {fmt.Sprintf("%02d", now.Year()%century)},
		"mm":    {fmt.Sprintf("%02d", now.Month())},
		"m":     {fmt.Sprintf("%d", now.Month())},
		"dd":    {fmt.Sprintf("%02d", now.Day())},
		"d":     {fmt.Sprintf("%d", now.Day())},
		"hh":    {fmt.Sprintf("%02d", now.Hour())},
		"h":     {fmt.Sprintf("%d", now.Hour())},
		"MM":    {fmt.Sprintf("%02d", now.Minute())},
		"M":     {fmt.Sprintf("%d", now.Minute())},
		"ss":    {fmt.Sprintf("%02d", now.Second())},
		"s":     {fmt.Sprintf("%d", now.Second())},
		"_name": {aliasURL.Host},
		"_path": {aliasURL.Path},
	}
	if aliasURL.User != nil {
		context["_user"] = []string{aliasURL.User.Username()}
		if password, ok := aliasURL.User.Password(); ok {
			context["_password"] = []string{password}
		}
	}
	for k, v := range aliasURL.Query() {
		context[k] = v
	}
	return context
}
