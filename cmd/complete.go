package cmd

import (
	"couture/internal/pkg/model/level"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/theme"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/posener/complete"
	"github.com/posener/complete/cmd/install"
	"github.com/willabides/kongplete"
	"os"
)

// 	FIXME completions are not working properly
func completionsHook() func(k *kong.Kong) error {
	return func(ctx *kong.Kong) error {
		var levelNames []string
		for _, lvl := range level.Levels {
			levelNames = append(levelNames, string(lvl))
		}
		commandName := ctx.Model.Name
		doInstall := os.Getenv("COMP_INSTALL") == "1"
		doUninstall := os.Getenv("COMP_UNINSTALL") == "1"
		if doInstall || doUninstall {
			kongplete.Complete(
				ctx,
				kongplete.WithPredictor("sources", complete.PredictNothing),
				kongplete.WithPredictor("time_format", complete.PredictSet(timeFormatNames...)),
				kongplete.WithPredictor("column_names", complete.PredictSet(column.Names()...)),
				kongplete.WithPredictor("themes", complete.PredictSet(theme.Names...)),
				kongplete.WithPredictor("width", complete.PredictSet("72", "80", "120", "132")),
				kongplete.WithPredictor("level", complete.PredictSet(levelNames...)),
			)
			if !doInstall || (doInstall && install.IsInstalled(commandName)) {
				_ = install.Uninstall(commandName)
				fmt.Println("completions uninstalled")
			}
			if doInstall {
				err := install.Install(commandName)
				if err != nil {
					return err
				}
				fmt.Println("completions installed")
			}
			os.Exit(0)
		}
		return nil
	}
}
