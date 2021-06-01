package cmd

// TODO alias groups
// TODO cleanup alias code

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

type aliasConfig struct {
	Groups  map[string][]string `json:"groups"`
	Aliases map[string]string   `json:"aliasConfig"`
}

func loadAliasConfig() (*aliasConfig, error) {
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
	return &cfg, nil
}

func expandAliases(args []string) ([]string, error) {
	cfg, err := loadAliasConfig()
	if err != nil {
		return nil, err
	}

	var expandedArgs []string
	for i := range args {
		var arg = args[i]
		expandedArgs = append(expandedArgs, expandSchemeShortForm(*cfg, arg)...)
	}
	for i := range expandedArgs {
		u, err := url.Parse(expandedArgs[i])
		if err == nil {
			if u.Scheme == "alias" {
				value, err := expandAlias(*cfg, u)
				if err != nil {
					return nil, err
				}
				if value != "" {
					expandedArgs[i] = value
				}
			}
		}
	}
	return expandedArgs, nil
}

func expandAlias(cfg aliasConfig, aliasURL *url.URL) (string, error) {
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
	context := map[string][]string{}
	addURLAliasVars(context, aliasURL)
	addDateAliasVars(context)
	return context
}

func expandSchemeShortForm(cfg aliasConfig, arg string) []string {
	const aliasURLPrefix = "alias://"
	if len(arg) > 1 && arg[0:2] == "@@" {
		group, ok := cfg.Groups[arg[2:]]
		if !ok {
			return nil
		}
		var aliases []string
		for _, alias := range group {
			aliases = append(aliases, aliasURLPrefix+alias)
		}
		return aliases
	} else if len(arg) > 0 && arg[0] == '@' {
		arg = aliasURLPrefix + arg[1:]
	}
	return []string{arg}
}

func addURLAliasVars(context map[string][]string, aliasURL *url.URL) {
	context["_name"] = []string{aliasURL.Host}
	context["_path"] = []string{aliasURL.Path}
	if aliasURL.User != nil {
		context["_user"] = []string{aliasURL.User.Username()}
		if password, ok := aliasURL.User.Password(); ok {
			context["_password"] = []string{password}
		}
	}
	for k, v := range aliasURL.Query() {
		context[k] = v
	}
}

func addDateAliasVars(context map[string][]string) {
	const century = 100

	now := time.Now()

	context["epoch"] = []string{fmt.Sprintf("%d", now.Unix())}

	context["yyyy"] = []string{fmt.Sprintf("%04d", now.Year())}
	context["yy"] = []string{fmt.Sprintf("%02d", now.Year()%century)}
	context["mm"] = []string{fmt.Sprintf("%02d", now.Month())}
	context["m"] = []string{fmt.Sprintf("%d", now.Month())}
	context["dd"] = []string{fmt.Sprintf("%02d", now.Day())}
	context["d"] = []string{fmt.Sprintf("%d", now.Day())}

	context["hh"] = []string{fmt.Sprintf("%02d", now.Hour())}
	context["h"] = []string{fmt.Sprintf("%d", now.Hour())}
	context["MM"] = []string{fmt.Sprintf("%02d", now.Minute())}
	context["M"] = []string{fmt.Sprintf("%d", now.Minute())}
	context["ss"] = []string{fmt.Sprintf("%02d", now.Second())}
	context["s"] = []string{fmt.Sprintf("%d", now.Second())}
}
