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

	// LessContrast ...
	LessContrast = DarkMode
	// MoreContrast ...
	MoreContrast = LightMode
)

// Mode ...
var Mode contrastPolarity

func init() {
	Mode = DarkMode
	if !termenv.HasDarkBackground() {
		Mode = LightMode
	}
}
