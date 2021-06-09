package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
)

func newContextColumn(style theme.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Context,
		layout,
		style,
		stringValue(func(event model.SinkEvent) string { return string(event.Context) }),
	)
}
