package color

import (
	"github.com/muesli/termenv"
)

type contrastPolarity int

const (
	// DarkMode ...
	DarkMode contrastPolarity = iota
	// LightMode ...
	LightMode

	// MoreNoticable ...
	MoreNoticable = DarkMode
	// LessNoticable ...
	LessNoticable = LightMode
)

// Mode ...
var Mode contrastPolarity

func init() {
	Mode = DarkMode
	if !termenv.HasDarkBackground() {
		Mode = LightMode
	}
}
