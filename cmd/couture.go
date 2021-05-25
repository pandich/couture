package main

import (
	"couture/cmd/cli"
)

func main() {
	mgr := cli.Start()
	(*mgr).Wait()
}
