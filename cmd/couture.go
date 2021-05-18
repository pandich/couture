package main

import (
	"couture/cmd/cli"
	"fmt"
	"os"
)

func main() {

	if err := cli.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
