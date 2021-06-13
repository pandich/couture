package cmd

import (
	"github.com/pandich/couture/sink"
	"github.com/pandich/couture/sink/layout"
	"github.com/pandich/couture/theme"
	"os"
	"time"
)

var enabled = true
var disabled = false

var defaultLayout = layout.Registry[layout.Default]
var defaultOut = os.Stdout
var defaultTheme *theme.Theme
var defaultTimeFormat = time.Stamp

var doricConfig = sink.Config{}
var doricConfigDefaults = sink.Config{
	AutoResize:       &enabled,
	Color:            &enabled,
	ConsistentColors: &enabled,
	Expand:           &disabled,
	Highlight:        &disabled,
	MultiLine:        &disabled,
	ShowSchema:       &disabled,
	Wrap:             &disabled,

	Layout:     &defaultLayout,
	Out:        defaultOut,
	Theme:      defaultTheme,
	TimeFormat: &defaultTimeFormat,
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
