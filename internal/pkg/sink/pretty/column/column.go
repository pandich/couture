package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink/pretty/config"
	"couture/internal/pkg/sink/pretty/theme"
	"couture/internal/pkg/source"
)

// TODO column widths should adapt to the terminal
// TODO handle config.ShowSigil
// TODO it is clumsy that the column Name has to be specified three times per column

type column interface {
	Name() string
	Register(theme.Theme)
	Formatter(source.Source, model.Event) string
	Renderer(config.Config, source.Source, model.Event) []interface{}
}

var columns = []column{
	sourceColumn{},
	timestampColumn{},
	applicationColumn{},
	threadColumn{},
	callerColumn{},
	levelColumn{},
	messageColumn{},
	stackTraceColumn{},
}
