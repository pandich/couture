package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/posener/complete/cmd/install"
	"github.com/willabides/kongplete"
	"os"
)

// FIXME completions are not working properly
func completionsHook() func(k *kong.Kong) error {
	return func(k *kong.Kong) error {
		commandName := k.Model.Name
		doInstall := os.Getenv("COMP_INSTALL") == "1"
		doUninstall := os.Getenv("COMP_UNINSTALL") == "1"
		if doInstall || doUninstall {
			kongplete.Complete(k)
			var err error
			if doInstall {
				if install.IsInstalled(commandName) {
					_ = install.Uninstall(commandName)
				}
				err = install.Install(commandName)
			} else {
				err = install.Uninstall(commandName)
			}
			if err != nil {
				return err
			}
			os.Exit(0)
		}
		return nil
	}
}
