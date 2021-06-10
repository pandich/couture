package cmd

import (
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/layout"
	"os"
	"time"
)

var enabled = true
var disabled = false
var defaultTheme = sink.Registry[sink.Prince]
var defaultLayout = layout.Registry[layout.Default]
var defaultTimeFormat = time.Stamp

var doricConfig = sink.Config{}
var doricConfigDefaults = sink.Config{
	AutoResize:       &enabled,
	Color:            &enabled,
	ConsistentColors: &enabled,
	Expand:           &disabled,
	Highlight:        &disabled,
	MultiLine:        &disabled,
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
