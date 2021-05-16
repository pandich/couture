package main

import (
	"couture/cmd/cli"
	"fmt"
	"os"
)

func main() {
	if err := cli.RunCommand(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
