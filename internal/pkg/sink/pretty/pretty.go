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
	"github.com/muesli/reflow/padding"
	"io"
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
	switch {
	case !cfg.Color:
		fallthrough
	case !isatty.IsTerminal(cfg.Out.Fd()) && !cfg.TTY:
		cfmt.DisableColors()
	default:
		cfmt.EnableColors()
	}
	if len(cfg.Columns) == 0 {
		cfg.Columns = column.DefaultColumns
	}
	var snk sink.Sink = &prettySink{
		terminalWidth: cfg.EffectiveTerminalWidth(),
		table:         column.NewTable(cfg),
		config:        cfg,
		out:           newOut(cfg.Out),
	}
	return &snk
}

// Init ...
func (snk *prettySink) Init(sources []*source.Source) {
	var sourceColors = map[model.SourceURL]string{}
	for _, src := range sources {
		sourceColors[(*src).URL()] = column.RegisterSource(snk.config.Theme, snk.config.ConsistentColors, *src)
	}
	if snk.config.Banner {
		snk.out <- snk.bannerLine(sources, sourceColors)
	}
}

// Accept ...
func (snk *prettySink) Accept(event model.SinkEvent) error {
	snk.out <- snk.table.RenderEvent(event)
	return nil
}

func newOut(writer io.WriteCloser) chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		writer := bufio.NewWriter(writer)
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

func (snk *prettySink) bannerLine(sources []*source.Source, sourceColors map[model.SourceURL]string) string {
	const oneHalf = 0.5
	const extraCharCount = uint(4)
	const minSourceWidth = uint(40)
	const maxSourceWidth = uint(float64(minSourceWidth) * 1.5)

	var width = uint(oneHalf * float64(config.TerminalWidth()))
	if width < minSourceWidth {
		width = minSourceWidth
	} else if width > maxSourceWidth {
		width = maxSourceWidth
	}
	actualWidth := width + extraCharCount

	var line = fmt.Sprintf("{{%s}}::bold|white|bgGray\n", padding.String("Legend:", actualWidth))
	for _, src := range sources {
		bg := sourceColors[(*src).URL()]
		fg, _ := colorful.MakeColor(gamut.Contrast(gamut.Hex(bg)))

		sigil := string((*src).Sigil())
		sourceURLFormat := fmt.Sprintf("{{%1.1[1]s ➥ %%-%[2]d.%[2]ds}}::%[3]s|bg%[4]s\n", sigil, width, fg.Hex(), bg)
		sourceURLString := (*src).URL().String()
		line += cfmt.Sprintf(sourceURLFormat, sourceURLString)
	}
	line += cfmt.Sprintf("{{%s}}::bold|white|bgGray\n", padding.String("꛳", actualWidth))
	return line
}
