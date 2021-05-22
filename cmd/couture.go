package main

import (
	"couture/cmd/cli"
	"couture/internal/pkg/model"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"time"
)

func main() {
	mgr := cli.Run()
	trapSignals(mgr)
	(*mgr).Wait()
}

func trapSignals(mgr *model.Manager) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, unix.SIGINT)
	signal.Notify(signalChan, unix.SIGKILL)

	go func() {
		const stopGracePeriod = 2 * time.Second
		if <-signalChan == unix.SIGINT {
			(*mgr).Stop()
			go func() {
				time.Sleep(stopGracePeriod)
				os.Exit(0)
			}()
			(*mgr).Wait()
		}
		os.Exit(0)
	}()
}
