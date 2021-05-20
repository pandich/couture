package cli

import (
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	errConfigNotFound = &viper.ConfigFileNotFoundError{}
)

func loadConfig() error {
	viper.SetConfigName("." + applicationName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil && !errors2.As(err, &errConfigNotFound) {
		return err
	}
	return nil
}

func aliasConfig() map[string]string {
	const aliasesConfigKey = "aliases"
	return viper.GetStringMapString(aliasesConfigKey)
}
