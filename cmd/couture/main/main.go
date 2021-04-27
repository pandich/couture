package main

import (
	"couture/cmd/couture/cli"
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"log"
)

func main() {
	cli.ParseCommandLine()

	mgr := manager.NewManager()

	if err := (*mgr).Register(source.Fake, sink.Console); err != nil {
		log.Fatal(err)
	}

	if err := (*mgr).Start(); err != nil {
		log.Fatal(err)
	}

	(*mgr).Wait()
}
