package cmd

// The alias mechanism provides short URIs which expand into command-line arguments/

import (
	"context"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/hashicorp/go-multierror"
	"github.com/pandich/couture/couture"
	errors2 "github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// The alias mechaism allows for full command-line arguments to this application to be represented by a short
// alias name. Additionally, groups of aliases can be defined and expanded into their individual aliases.
// The user's configuration file (see Config) defines the aliases and groups. See expandAliases to see how
// the confgutation is loaded from the user's home.
// This is useful for shortening long URIs and for defining groups of common command arguments.
import (
	"path"
)

import (
	"gopkg.in/yaml.v3"
)

const (
	// aliasScheme is the URI scheme for aliases. Aliases look like: alias://...
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

	Context map[string]string `yaml:"-" env:"COUTURE_CONTEXT"`
}

// expandAliases evaluates each argument and expands any aliases or groups.
// Errors are returned if expanded arguments result in malformed URIs.
func (config *aliasConfig) expandAliases(args []string) ([]string, error) {
	var (
		err      error
		expanded []string
		aliasURI *url.URL
		value    string
	)

	// expand all the arguments
	for i := range args {
		expanded = append(expanded, config.expandArgument(args[i])...)
	}

	errs := &multierror.Error{}

	// final all arguments within the expanded list which are alias URIs
	for i := range expanded {
		// if the current arg doesn't parse as a URI, it is a literal: no action needed.
		if aliasURI, err = url.Parse(expanded[i]); err != nil {
			continue
		}

		// if the current arg is a URL but not an alias: no action needed.
		if aliasScheme != aliasURI.Scheme {
			continue
		}

		value, err = expandAliasURL(config, aliasURI)

		switch {
		// if the alias URI could not be examded, add the failure to the list of
		// expansion errors.
		case err != nil:
			errs = multierror.Append(errs, err)

			// If the alias was found, replace the expanded argument with
			// the alias lookup.
		case value != "":
			expanded[i] = value
		}
	}

	return expanded, errs.ErrorOrNil()
}

// loadAliasConfig loads the user's alias configuration file. Errors are returned if the file cannot be read or
// is malforemd. Config file is located in $HOME/.config/counture/aliases.yaml
func loadAliasConfig() (*aliasConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// if we can't find the aliase file under either yaml extensions
	// return without error.
	aliasFilename := path.Join(home, ".config", couture.Name, "aliases.yaml")
	if !fileutil.Exist(aliasFilename) {

		aliasFilename = path.Join(home, ".config", couture.Name, "aliases.yml")
		if !fileutil.Exist(aliasFilename) {
			return &aliasConfig{}, nil
		}
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

	ctx := context.Background()
	if err = envconfig.Process(ctx, config); err != nil {
		log.Fatal(err)
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
	case strings.HasPrefix(arg, groupNamePrefix):
		// then expand the group into its aliases

		groupName := arg[len(groupNamePrefix):]
		query := ""
		if strings.Contains(groupName, "?") {
			parts := strings.SplitN(groupName, "?", 2)
			groupName = parts[0]
			query = parts[1]
		}

		group, ok := config.Groups[groupName]
		if !ok {
			return nil
		}

		for _, alias := range group {
			v := aliasURIPrefix + alias
			if query != "" {
				v += "?" + query
			}
			args = append(args, v)
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
	m := aliasContext(aliasURL)
	for k, v := range config.Context {
		m[k] = []string{v}
	}
	return raymond.Render(alias, m)
}

// aliasContext sets up the global properties usable in an alias template.
func aliasContext(aliasURL *url.URL) map[string][]string {
	const century = 100

	now := time.Now()

	ctx := map[string][]string{
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

	for _, v := range aliasURL.Query() {
		ctx[v[0]] = v[1:]
	}

	if aliasURL.User != nil {
		ctx["_user"] = []string{aliasURL.User.Username()}
		if password, ok := aliasURL.User.Password(); ok {
			ctx["_password"] = []string{password}
		}
	}

	for k, v := range aliasURL.Query() {
		ctx[k] = v
	}

	for _, s := range os.Environ() {
		parts := strings.SplitN(s, "=", 2)
		k, v := parts[0], parts[1]
		if strings.HasPrefix(k, "COUTURE_CONTEXT_") {
			effectiveKey := strings.TrimPrefix(k, "COUTURE_CONTEXT_")
			ctx[effectiveKey] = []string{v}
		}
	}

	return ctx
}
