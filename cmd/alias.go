package cmd

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/gagglepanda/couture/couture"
	"github.com/hashicorp/go-multierror"
	errors2 "github.com/pkg/errors"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"
)

// The alias mechaism allows for full command-line arguments to this application to be represented by a short
// alias name. Additionally, groups of aliases can be defined and expanded into their individual aliases.
// The user's configuration file (see Config) defines the aliases and groups. See expandAliases to see how
// the confgutation is loaded from the user's home.
import (
	"gopkg.in/yaml.v2"
)

const (
	// aliasScheme is the URI scheme for aliases.
	aliasScheme = "alias"

	// aliasNamePrefix allows a short-form alias to be specified. (e.g, @foo instead of alias://foo). The @ is
	// a pnemonic for "alias".
	aliasNamePrefix = "@"

	// groupNamePrefix allows for all configured groups (see Config.Groups) to be expanded into their individual
	// aliases. (e.g., @@foo expands to alias://bar and alias://baz give a config with a key 'foo' with entries
	// 'bar' and 'baz'). The '@@ pneumatic for more than one alias.
	groupNamePrefix = aliasNamePrefix + aliasNamePrefix

	// aliasURIPrefix is canonical the prefix for all alias URIs.
	aliasURIPrefix = aliasScheme + "://"
)

// aliasConfig is the structure of the YAML file that defines aliases.
type aliasConfig struct {
	// Groups is a map of group names to a list of aliases (see Aliases below).
	// The group name is used as the alias name.
	Groups map[string][]string `yaml:"groups,omitempty"`
	// Aliases is a
	Aliases map[string]string `yaml:"aliases,omitempty"`
}

func (config *aliasConfig) expandAliases(args []string) ([]string, error) {
	var (
		err      error
		expanded []string
		aliasURI *url.URL
		value    string
	)

	for i := range args {
		expanded = append(expanded, config.expandArgument(args[i])...)
	}

	errs := &multierror.Error{}

	for i := range expanded {
		aliasURI, err = url.Parse(expanded[i])
		if err == nil && aliasScheme == aliasURI.Scheme {
			value, err = expandAliasURL(config, aliasURI)

			switch {
			case err != nil:
				errs = multierror.Append(errs, err)

			case value != "":
				expanded[i] = value
			}
		}
	}

	return expanded, errs.ErrorOrNil()
}

func loadAliasConfig() (*aliasConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	aliasFilename := path.Join(home, ".config", couture.Name, "aliases.yaml")
	if !fileutil.Exist(aliasFilename) {
		return &aliasConfig{}, nil
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

	config := &aliasConfig{}
	err = yaml.Unmarshal(s, config)
	if err != nil {
		return nil, err
	}

	return config, err
}

// expandArgument expands an alias string into its component parts,
// replace group, and alias short form with full URL forms.
func (config *aliasConfig) expandArgument(arg string) []string {
	var args []string

	// perform alias expansion on the arguments.
	switch {

	// is this a group reference?
	case len(arg) > len(groupNamePrefix) && arg[:len(groupNamePrefix)] == groupNamePrefix:
		// then expand the group into its aliases

		groupName := arg[len(groupNamePrefix):]
		group, ok := config.Groups[groupName]
		if !ok {
			return nil
		}

		for _, alias := range group {
			args = append(args, aliasURIPrefix+alias)
		}

	// is this an alias reference?
	case len(arg) > len(aliasNamePrefix) && arg[:len(aliasNamePrefix)] == aliasNamePrefix:
		// then replace its short prefix with the URL prefix.
		args = append(args, aliasURIPrefix+arg[len(aliasNamePrefix):])

	// otherwise add the argument as-is.
	default:
		args = append(args, arg)

	}

	return args
}

func expandAliasURL(config *aliasConfig, aliasURL *url.URL) (string, error) {
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

	now := time.Now()

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
