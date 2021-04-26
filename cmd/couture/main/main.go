package main

import (
	"couture/internal/pkg/manager"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/source"
	"log"
	"time"
)

func main() {
	mgr := manager.NewManager()

	if err := (*mgr).RegisterSink(sink.Diagnostic, sink.Logrus); err != nil {
		log.Fatal(err)
	}

	(*mgr).RegisterSource(1*time.Second, source.FakeSource)

	if err := (*mgr).Start(); err != nil {
		log.Fatal(err)
	}

	(*mgr).Wait()
}
