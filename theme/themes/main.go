package main

import (
	"github.com/pandich/couture/color"
	"github.com/pandich/couture/theme"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

var definitions = map[string]theme.Template{
	"prince__dark":  theme.SplitComplementary(color.DarkMode, color.Hex("#9b99bf")),
	"prince__light": theme.SplitComplementary(color.LightMode, color.Hex("#9b99bf")),
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
