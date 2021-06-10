package cmd

import (
	"couture/internal/pkg/sink/doric/config"
	layout2 "couture/internal/pkg/sink/layout"
	theme2 "couture/internal/pkg/sink/theme"
	"os"
	"time"
)

var enabled = true
var disabled = false
var defaultTheme = theme2.Registry[theme2.Prince]
var defaultLayout = layout2.Registry[layout2.Default]
var defaultTimeFormat = time.Stamp

var doricConfig = config.Config{}
var doricConfigDefaults = config.Config{
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
		doricConfig.FillMissing(*userCfg)
	}
	doricConfig.FillMissing(doricConfigDefaults)
	return nil
}
