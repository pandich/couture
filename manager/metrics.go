package manager

import (
	"bufio"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"os"
	"sync"
)

var metricsDumpLock = sync.Mutex{}

// dumpMetrics to STDERR and exit cleanly. This is for development/diagnostic purposes only.
func dumpMetrics() {
	defer os.Exit(0)

	metricsDumpLock.Lock()
	defer metricsDumpLock.Unlock()

	out := os.Stderr
	defer bufio.NewWriter(out).Flush()

	// loop over all metrics and print their name and value
	metrics.DefaultRegistry.Each(
		func(name string, _ interface{}) {
			switch metric := metrics.Get(name).(type) {

			case metrics.Counter:
				snapshot := metric.Snapshot()
				_, _ = fmt.Fprintf(
					out,
					"%s: count=%d\n",
					name,
					snapshot.Count(),
				)

			case metrics.Meter:
				snapshot := metric.Snapshot()
				_, _ = fmt.Fprintf(
					out,
					"%s: count=%d, rate(sec)=%0.2f\n",
					name,
					snapshot.Count(),
					snapshot.RateMean(),
				)

			case metrics.Timer:
				snapshot := metric.Snapshot()
				_, _ = fmt.Fprintf(
					out,
					"%s: count=%d, rate(sec)=%0.2f, time(sec)=%0.2f\n",
					name,
					snapshot.Count(),
					snapshot.RateMean(),
					snapshot.Mean(),
				)

			}
		},
	)
}
