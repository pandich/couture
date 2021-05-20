package cli

import (
	"couture/internal/pkg/manager"
	errors2 "github.com/pkg/errors"
)

// RunApplication ...
func RunApplication() {
	parser := newParser()

	err := loadConfig()
	parser.FatalIfErrorf(err)

	ctx, err := parser.Parse(evaluateArgs())
	parser.FatalIfErrorf(err)

	err = runCommand(ctx.Command())
	parser.FatalIfErrorf(err)
}

func runCommand(command string) error {
	const logCommand = "log <url>"

	switch command {
	case logCommand:
		return runManagerCommand()
	default:
		return errors2.Errorf("command not implemented: %s\n", command)
	}
}

func runManagerCommand() error {
	mgrOptions, err := getFlags()
	if err != nil {
		return err
	}
	sources, err := getArgs()
	if err != nil {
		return err
	}
	options := append(mgrOptions, sources...)

	mgr, err := manager.New(options...)
	if err != nil {
		return err
	}
	return (*mgr).Start()
}
