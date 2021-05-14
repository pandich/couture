package main

import (
	"couture/cmd/cli"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink/pretty"
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
	options = append(options, pretty.New())

	mgr, err := manager.New(options...)
	if err != nil {
		die(err)
	}

	if err := (*mgr).Start(); err != nil {
		die(err)
	}

	os.Exit(0)
}
