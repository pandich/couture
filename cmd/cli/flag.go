package cli

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/json"
	"couture/internal/pkg/sink/pretty"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/gamut"
	errors2 "github.com/pkg/errors"
	"image/color"
	"io"
	"os"
	"os/exec"
	"strings"
)

const purpleRain = "#AE99BF"
const crystalBlue = "#A0D7D9"
const paleChestnut = "#E2A3B4"
const celadonGreen = "#A5E0B6"

var themeColors = map[string]color.Color{
	"prince":  gamut.Hex(purpleRain),
	"ocean":   gamut.Hex(crystalBlue),
	"warm":    gamut.Hex(paleChestnut),
	"forrest": gamut.Hex(celadonGreen),
}

func getFlags() ([]interface{}, error) {
	snk, err := sinkFlag()
	if err != nil {
		return nil, err
	}
	return []interface{}{
		manager.LogLevelOption(cli.Level),
		manager.FilterOption(cli.Include, cli.Exclude),
		manager.SinceOption(cli.Since),
		snk,
	}, nil
}

func sinkFlag() (interface{}, error) {
	paginator, err := paginatorFlag()
	if err != nil {
		return nil, err
	}
	switch cli.OutputFormat {
	case "json":
		return json.New(paginator), nil
	case "pretty":
		return pretty.New(paginator, cli.Wrap, themeFlag()), nil
	default:
		return nil, errors2.Errorf("unknown output format: %s\n", cli.OutputFormat)
	}
}

func paginatorFlag() (io.Writer, error) {
	if !cli.Paginate {
		return os.Stdout, nil
	}

	var pagerArgs = strings.Split(cli.Paginator, " \t\n")
	pager, pagerArgs := pagerArgs[0], pagerArgs[1:]
	pagerCmd := exec.Command(pager, pagerArgs...)

	// I/O
	pagerCmd.Stdout, pagerCmd.Stderr = os.Stdout, os.Stderr
	writer, err := pagerCmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err = pagerCmd.Start(); err != nil {
		return nil, err
	}
	return writer, nil
}

func themeFlag() color.Color {
	if cli.Theme == "none" {
		cfmt.DisableColors()
		return color.White
	}
	if c, ok := themeColors[cli.Theme]; ok {
		return c
	}
	return color.White
}
