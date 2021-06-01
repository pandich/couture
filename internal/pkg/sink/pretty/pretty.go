package pretty

import (
	"bufio"
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"couture/internal/pkg/sink/pretty/column"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/source"
	"fmt"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mattn/go-isatty"
	"github.com/muesli/gamut"
	"os"
)

// Name ...
const Name = "pretty"

// prettySink provides render output.
type prettySink struct {
	terminalWidth uint
	table         *column.Table
	config        config.Config
	out           chan string
}

// New provides a configured prettySink sink.
func New(cfg config.Config) *sink.Sink {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	isBlackOrWhite := cfg.Theme.BaseColor == "#ffffff"
	if !isTTY || isBlackOrWhite {
		cfmt.DisableColors()
	}

	if len(cfg.Columns) == 0 {
		cfg.Columns = column.DefaultColumns
	}

	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		table:         column.NewTable(cfg),
		config:        cfg,
		out:           newOut(),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []*source.Source) {
	const minSourceWidth = 30

	var sourceColors = map[model.SourceURL]string{}
	for _, src := range sources {
		sourceColors[(*src).URL()] = column.RegisterSource(snk.config.Theme, snk.config.ConsistentColors, *src)
	}
	if snk.config.Banner {
		_, _ = cfmt.Println("{{ Legend: }}::bold|bgWhite|black")
		for _, src := range sources {
			bg := sourceColors[(*src).URL()]
			fg, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(bg)))

			sigil := string((*src).Sigil())
			var width = int(float64(config.TerminalWidth()) * 0.5)
			if width < minSourceWidth {
				width = minSourceWidth
			}
			if width > minSourceWidth*2 {
				width = minSourceWidth * 2
			}
			sourceURLFormat := fmt.Sprintf("{{%1.1[1]s ➥ %%-%[2]d.%[2]ds}}::%[3]s|bg%[4]s", sigil, width, fg.Hex(), bg)
			sourceURLString := (*src).URL().String()
			sourceURLBanner := fmt.Sprintf(sourceURLFormat, sourceURLString)
			_, _ = cfmt.Println(sourceURLBanner)
		}
		_, _ = cfmt.Println("\n")
		os.Exit(0)
	}
}

// Accept ...
func (snk *prettySink) Accept(event model.SinkEvent) error {
	snk.out <- snk.table.RenderEvent(event)
	return nil
}

func newOut() chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		writer := bufio.NewWriter(os.Stdout)
		for {
			message := <-out
			_, err := writer.WriteString(message + "\n")
			if err != nil {
				panic(err)
			}
			err = writer.Flush()
			if err != nil {
				panic(err)
			}
		}
	}()
	return out
}
