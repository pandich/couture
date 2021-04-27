package main

import (
	"couture/cmd/cli"
	"couture/cmd/config"
	"couture/internal/pkg/manager"
)

func main() {
	cli.Parse()

	mgr := manager.NewManager()
	(*mgr).MustRegister(cli.Sources()...)
	(*mgr).MustRegister(cli.Sinks()...)
	(*mgr).MustRegister(config.Sources()...)
	(*mgr).MustRegister(config.Sinks()...)
	(*mgr).MustStart()
}
