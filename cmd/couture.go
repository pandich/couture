package main

import (
	"couture/cmd/cli"
	"couture/internal/pkg/model"
	"github.com/muesli/termenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mgr := cli.Run()
	trapInterrupt(mgr)
	(*mgr).Wait()
}

func trapInterrupt(mgr *model.Manager) {
	const stopGracePeriod = 2 * time.Second

	die := func() {
		println()
		termenv.Reset()
		os.Exit(0)
	}

	interruption := make(chan os.Signal)
	signal.Notify(interruption, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruption
		(*mgr).Stop()
		go func() { time.Sleep(stopGracePeriod); die() }()
		(*mgr).Wait()
		die()
	}()
}
