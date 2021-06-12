package main

import (
	"github.com/pandich/couture/theme"
	color2 "github.com/pandich/couture/theme/color"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

var definitions = map[string]theme.Generator{
	"prince__dark":  theme.SplitComplementaryGenerator(color2.DarkMode, color2.Hex("#9b99bf")),
	"prince__light": theme.SplitComplementaryGenerator(color2.LightMode, color2.Hex("#9b99bf")),
}

func writeThemeToFile(name string, th theme.Theme) error {
	filename := name + ".yaml"
	b, err := yaml.Marshal(th)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.WriteString(f, string(b))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	for name, pal := range definitions {
		err := writeThemeToFile(name, pal.AsTheme())
		if err != nil {
			log.Fatalln(err)
		}
	}
}
