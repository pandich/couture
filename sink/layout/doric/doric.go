package doric

import (
	"github.com/pandich/couture/event"
	"github.com/pandich/couture/sink"
	column2 "github.com/pandich/couture/sink/layout/doric/column"
	"github.com/pandich/couture/source"
	"github.com/i582/cfmt/cmd/cfmt"
)

const sinkName = "doric"

// doricSink provides render output.
type doricSink struct {
	terminalWidth uint
	table         *column2.Table
	config        sink.Config
	out           chan string
}

// New provides a configured doricSink sink.
func New(config sink.Config) sink.Sink {
	return &doricSink{
		terminalWidth: config.EffectiveTerminalWidth(),
		table:         column2.NewTable(config),
		config:        config,
		out:           sink.NewOut(sinkName, config.Out),
	}
}

// Init ...
func (snk doricSink) Init(sources []*source.Source) {
	switch {
	case snk.config.Color != nil && !*snk.config.Color:
		cfmt.DisableColors()
	case snk.config.EffectiveIsTTY():
		cfmt.EnableColors()
	default:
		cfmt.DisableColors()
	}
	var sourceColors = map[event.SourceURL]string{}
	for _, src := range sources {
		consistentColors := *snk.config.ConsistentColors
		style := snk.config.Theme.AsHexPair(consistentColors, *src)
		sourceColors[(*src).URL()] = style.Bg
		column2.RegisterSourceStyle(style, snk.config.Layout.Source, *src)
	}
}

// Accept ...
func (snk doricSink) Accept(event event.SinkEvent) error {
	snk.out <- snk.table.Render(event)
	return nil
}
