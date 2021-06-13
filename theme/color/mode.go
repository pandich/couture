package color

import (
	"github.com/muesli/termenv"
)

// ContrastPolarity ...
type ContrastPolarity int

const (
	// DarkMode ...
	DarkMode ContrastPolarity = iota
	// LightMode ...
	LightMode

	// MoreContrast ...
	MoreContrast = DarkMode
	// LessContrast ...
	LessContrast = LightMode
)

// Mode ...
var Mode ContrastPolarity

func init() {
	Mode = DarkMode
	if !termenv.HasDarkBackground() {
		Mode = LightMode
	}
}
