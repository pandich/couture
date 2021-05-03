package main

import (
	"couture/cmd/cli"
	"couture/internal/pkg/manager"
	"fmt"
	"os"
)

func main() {
	die := func(err error) {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	options, err := cli.MustParseOptions()
	if err != nil {
		die(err)
	}

	mgr, err := manager.New(options...)
	if err != nil {
		die(err)
	}

	if err := (*mgr).Start(); err != nil {
		die(err)
	}

	os.Exit(0)
}
