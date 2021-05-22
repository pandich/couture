package cli

import (
	"couture/internal/pkg/couture"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
)

func loadConfig() error {
	errConfigNotFound := &viper.ConfigFileNotFoundError{}

	viper.SetConfigName("." + couture.Name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil && !errors2.As(err, &errConfigNotFound) {
		return err
	}
	return nil
}

func aliasConfig() map[string]string {
	const aliasesConfigKey = "aliases"
	return viper.GetStringMapString(aliasesConfigKey)
}
