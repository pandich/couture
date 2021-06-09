package column

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/model/layout"
	"couture/internal/pkg/model/theme"
	"couture/internal/pkg/schema"
)

func newApplicationColumn(style theme.Style, layout layout.ColumnLayout) column {
	return newWeightedColumn(
		schema.Application,
		layout,
		style,
		stringValue(func(event model.SinkEvent) string { return string(event.Application) }),
	)
}
