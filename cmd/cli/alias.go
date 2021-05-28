package cli

// TODO cleanup alias code

import (
	"couture/internal/pkg/couture"
	"fmt"
	"github.com/aymerick/raymond"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/url"
	"regexp"
	"time"
)

func loadAliasConfig() error {
	errConfigNotFound := &viper.ConfigFileNotFoundError{}
	viper.SetConfigName("aliases")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/" + couture.Name)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil && !errors2.As(err, errConfigNotFound) {
		return err
	}
	return nil
}

func aliasConfig() map[string]string {
	const aliasesConfigKey = "aliases"
	return viper.GetStringMapString(aliasesConfigKey)
}

func expandAliases(args []string) ([]string, error) {
	for i := range args {
		var arg = args[i]
		arg = expandSchemeShortForm(arg)
		u, err := url.Parse(arg)
		if err == nil {
			if u.Scheme == "alias" {
				value, err := expandAlias(u)
				if err != nil {
					return nil, err
				}
				if value != "" {
					args[i] = value
				}
			}
		}
	}
	return args, nil
}

func expandAlias(aliasURL *url.URL) (string, error) {
	simpleArgs := regexp.MustCompile(`@(?P<name>\w+)`)

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

	return raymond.Render(alias, aliasContext(aliasURL))
}

func aliasContext(aliasURL *url.URL) map[string][]string {
	context := map[string][]string{}
	addURLAliasVars(context, aliasURL)
	addDateAliasVars(context)
	return context
}

func expandSchemeShortForm(arg string) string {
	if len(arg) > 0 && arg[0] == '@' {
		arg = "alias://" + arg[1:]
	}
	return arg
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
