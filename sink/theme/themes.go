package theme

import (
	"sort"
	"strings"
)

var themeColors = map[string]string{
	"halloween": "Burnt Orange",
	"land":      "Ochre",
	"prince":    "Logan",
	"sea":       "Ocean Blue",
	"sky":       "Sky Blue",
	"tango":     "Tangerine",
}

// Names ...
func Names() []string {
	var names []string
	for name := range themeColors {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool { return strings.Compare(names[i], names[j]) < 0 })
	return names
}
