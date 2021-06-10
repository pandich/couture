package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/source"
	"github.com/i582/cfmt/cmd/cfmt"
	"time"
)

// prettySink provides render output.
type prettySink struct {
	terminalWidth uint
	table         *column.Table
	config        config.Config
	out           chan string
}

// New provides a configured prettySink sink.
func New(cfg config.Config) *sink.Sink {
	switch {
	case cfg.Color != nil && !*cfg.Color:
		cfmt.DisableColors()
	case cfg.EffectiveIsTTY():
		cfmt.EnableColors()
	default:
		cfmt.DisableColors()
	}
	if len(cfg.Columns) == 0 {
		cfg.Columns = column.DefaultColumns
	}
	if cfg.TimeFormat == nil {
		tf := time.Stamp
		cfg.TimeFormat = &tf
	}
	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		table:         column.NewTable(cfg),
		config:        cfg,
		out:           sink.NewOut("pretty", cfg.Out),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []*source.Source) {
	var sourceColors = map[model.SourceURL]string{}
	for _, src := range sources {
		consistentColors := *snk.config.ConsistentColors
		style := snk.config.Theme.SourceStyle(consistentColors, *src)
		sourceColors[(*src).URL()] = style.Bg
		column.RegisterSourceStyle(style, snk.config.Layout.Source, *src)
	}
}

// Accept ...
func (snk *prettySink) Accept(event model.SinkEvent) error {
	snk.out <- snk.table.Render(event)
	return nil
}
