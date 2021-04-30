package main

import (
	"couture/cmd/cli"
	"couture/cmd/config"
	"couture/internal/pkg/manager"
)

func main() {
	ctx := cli.MustLoad()
	config.MustLoad(ctx)

	mgr := manager.NewManager()

	(*mgr).MustRegister(cli.Options()...)
	(*mgr).MustRegister(config.Options()...)

	(*mgr).MustRegister(cli.Sources()...)
	(*mgr).MustRegister(config.Sources()...)

	(*mgr).MustRegister(cli.Sinks()...)
	(*mgr).MustRegister(config.Sinks()...)

	(*mgr).MustStart()
}
