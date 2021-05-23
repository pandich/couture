package main

import (
	"couture/cmd/cli"
	"couture/internal/pkg/model"
	"couture/internal/pkg/tty"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mgr := cli.Run()
	trapSignals(mgr)
	(*mgr).Wait()
}

// FIXME exiting the program is very problematic - ctrl-c rarely works
func trapSignals(mgr *model.Manager) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	signal.Notify(signalChan, syscall.SIGTERM)

	go func() {
		const stopGracePeriod = 2 * time.Second
		if <-signalChan == syscall.SIGINT {
			if tty.IsTTY() {
				(*mgr).Stop()
				go func() { time.Sleep(stopGracePeriod); os.Exit(1) }()
				(*mgr).Wait()
			}
		}
		os.Exit(0)
	}()
}
