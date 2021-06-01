package cmd

import (
	"couture/internal/pkg/couture"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/aymerick/raymond"
	"github.com/coreos/etcd/pkg/fileutil"
	errors2 "github.com/pkg/errors"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"
)

var now = time.Now()

const aliasScheme = "alias"

type aliasConfig struct {
	Groups  map[string][]string `json:"groups"`
	Aliases map[string]string   `json:"aliasConfig"`
}

func expandAliases(args []string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	aliasFilename := path.Join(home, ".config", couture.Name, "aliases.toml")
	if !fileutil.Exist(aliasFilename) {
		return nil, nil
	}
	aliasFile, err := os.Open(aliasFilename)
	if err != nil {
		return nil, err
	}
	defer aliasFile.Close()
	s, err := ioutil.ReadAll(aliasFile)
	if err != nil {
		return nil, err
	}

	var cfg aliasConfig
	err = toml.Unmarshal(s, &cfg)
	if err != nil {
		return nil, err
	}

	var expandedArgs []string
	for i := range args {
		expandedArgs = append(expandedArgs, expandAliasString(cfg, args[i])...)
	}
	for i := range expandedArgs {
		if u, err := url.Parse(expandedArgs[i]); err == nil && u.Scheme == "alias" {
			value, err := expandAliasURL(cfg, u)
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

func expandAliasString(cfg aliasConfig, arg string) []string {
	const aliasNamePrefix = "@"
	const groupNamePrefix = "@@"
	const aliasURLPrefix = aliasScheme + "://"

	var args []string
	switch {
	case len(arg) > len(groupNamePrefix) && arg[:len(groupNamePrefix)] == groupNamePrefix:
		groupName := arg[len(groupNamePrefix):]
		group, ok := cfg.Groups[groupName]
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

func expandAliasURL(cfg aliasConfig, aliasURL *url.URL) (string, error) {
	simpleArgs := regexp.MustCompile(`@(?P<name>\w+)`)
	alias, ok := cfg.Aliases[aliasURL.Host]
	if !ok {
		return "", errors2.Errorf("unknown alias: %s", aliasURL.Host)
	}
	// expand simple value placeholders
	alias = simpleArgs.ReplaceAllString(alias,
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
