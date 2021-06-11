package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/theme"
	"gopkg.in/yaml.v3"
)

// TODO https://github.com/lucasb-eyer/go-colorful
func main() {
	maybePanic := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	warm, err := colorful.WarmPalette(10)
	maybePanic(err)
	sources, err := colorful.WarmPalette(180)
	maybePanic(err)

	var sourceStlyes []theme.Style
	for _, source := range sources {
		sourceStlyes = append(sourceStlyes, theme.Style{Fg: source.Hex(), Bg: "#000000"})
	}

	th := theme.Theme{
		Timestamp:       theme.Style{Fg: warm[0].Hex(), Bg: "#0000000"},
		Application:     theme.Style{Fg: warm[1].Hex(), Bg: "#0000000"},
		Context:         theme.Style{Fg: warm[2].Hex(), Bg: "#0000000"},
		Entity:          theme.Style{Fg: warm[3].Hex(), Bg: "#0000000"},
		ActionDelimiter: theme.Style{Fg: warm[4].Hex(), Bg: "#0000000"},
		Action:          theme.Style{Fg: warm[5].Hex(), Bg: "#0000000"},
		LineDelimiter:   theme.Style{Fg: warm[6].Hex(), Bg: "#0000000"},
		Line:            theme.Style{Fg: warm[7].Hex(), Bg: "#0000000"},
		Level: map[level.Level]theme.Style{
			level.Trace: {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Debug: {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Info:  {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Warn:  {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Error: {Fg: warm[1].Hex(), Bg: "#000000"},
		},
		Message: map[level.Level]theme.Style{
			level.Trace: {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Debug: {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Info:  {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Warn:  {Fg: warm[1].Hex(), Bg: "#000000"},
			level.Error: {Fg: warm[1].Hex(), Bg: "#000000"},
		},
		Source: sourceStlyes,
	}
	b, err := yaml.Marshal(th)
	maybePanic(err)
	fmt.Println(string(b))
}
