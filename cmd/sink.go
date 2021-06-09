package cmd

import (
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/sink/pretty/config"
	"os"
	"time"
)

var enabled = true
var disabled = false
var defaultTheme = theme.Registry[theme.Prince]
var defaultLayout = layout.Registry[layout.Default]
var defaultTimeFormat = time.Stamp

var prettyConfig = config.Config{}
var prettyConfigDefaults = config.Config{
	AutoResize:       &enabled,
	Color:            &enabled,
	ConsistentColors: &enabled,
	Expand:           &disabled,
	Highlight:        &disabled,
	Multiline:        &disabled,
	Out:              os.Stdout,
	ShowSchema:       &disabled,
	Theme:            &defaultTheme,
	Layout:           &defaultLayout,
	TimeFormat:       &defaultTimeFormat,
	Wrap:             &disabled,
}

func loadSinkConfig() error {
	userCfg, err := loadConfig()
	if err != nil {
		return err
	}

	if userCfg != nil {
		prettyConfig.FillMissing(*userCfg)
	}
	prettyConfig.FillMissing(prettyConfigDefaults)
	return nil
}
