package theme

import (
	"github.com/muesli/gamut/palette"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/sink/color"
	errors2 "github.com/pkg/errors"
)

// GenerateTheme ...
func GenerateTheme(colorName string) (*Theme, error) {
	if s, ok := themeColors[colorName]; ok {
		colorName = s
	}
	ac, ok := color.ByName(colorName)
	if !ok {
		return nil, errors2.Errorf("invalid theme color: %s", colorName)
	}
	return splitComplementaryGenerator(ac).asTheme(), nil
}

func splitComplementaryGenerator(baseColor color.AdaptorColor) generator {
	const triadicDirectionCutoff = 0.5 // 180º
	var messageColorTriadicIndex = 1
	if h, _, _ := baseColor.AsColorfulColor().Hsl(); h < triadicDirectionCutoff {
		messageColorTriadicIndex = 0
	}

	return generator{
		ApplicationColor: baseColor.
			Analogous()[1].
			AdjustConstrast(color.MoreContrast, 20),
		TimestampColor: baseColor.
			Complementary().
			Monochromatic()[0xB0],
		EntityColor: baseColor,
		MessageColor: baseColor.
			Triadic()[messageColorTriadicIndex].
			AdjustConstrast(color.MoreContrast, 80),
	}
}

type generator struct {
	ApplicationColor color.AdaptorColor
	TimestampColor   color.AdaptorColor
	EntityColor      color.AdaptorColor
	MessageColor     color.AdaptorColor
}

func (p generator) asTheme() *Theme {
	th := Theme{}

	// order is important in this block
	p.applyEntity(&th)
	p.applyHeader(&th)
	p.applyLevels(&th)
	p.applyMessages(&th)
	p.applySources(&th)

	return &th
}

func (p generator) applySources(th *Theme) {
	const sourcePaletteSize = 100
	baseColor := th.base()
	complementaryColor := baseColor.Complementary()
	sourcePalette := baseColor.
		ToPleasingPalette(sourcePaletteSize).
		Clamped(palette.Crayola)
	for _, paletteColor := range sourcePalette {
		sourceColor := paletteColor.Blend(complementaryColor, 10)
		th.Source = append(th.Source, sourceColor.AsHexPair())
	}
}
func (p generator) applyHeader(th *Theme) {
	th.Application = color.FgBgTuple{
		Fg: p.ApplicationColor.AsHexColor(),
		Bg: p.ApplicationColor.AdjustConstrast(color.LessContrast, 90).AsHexColor(),
	}
	th.Timestamp = color.FgBgTuple{
		Fg: p.TimestampColor.AsHexColor(),
		Bg: p.TimestampColor.AdjustConstrast(color.LessContrast, 80).AsHexColor(),
	}
}

func (p generator) applyEntity(th *Theme) {
	entityFg := p.EntityColor
	entityBg := entityFg.AdjustConstrast(color.LessContrast, 80)

	th.Entity = color.FgBgTuple{
		Fg: entityFg.AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
	th.Context = color.FgBgTuple{
		Fg: entityFg.AsHexColor(),
		Bg: entityFg.AdjustConstrast(color.LessContrast, 70).AsHexColor(),
	}

	th.Line = color.FgBgTuple{
		Fg: entityFg.AdjustConstrast(color.LessContrast, 20).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}

	th.LineDelimiter = color.FgBgTuple{
		Fg: entityFg.AdjustConstrast(color.LessContrast, 20).AdjustConstrast(color.MoreContrast, 40).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}

	actionFg := entityFg.AdjustConstrast(color.MoreContrast, 20)
	th.Action = color.FgBgTuple{
		Fg: actionFg.AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
	th.ActionDelimiter = color.FgBgTuple{
		Fg: actionFg.AdjustConstrast(color.MoreContrast, 40).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
}

func (p generator) applyLevels(th *Theme) {
	styleForName := func(name string) color.FgBgTuple {
		c := color.MustByName(name).
			Blend(p.EntityColor, 5)
		return color.FgBgTuple{Fg: c.Contrast().AsHexColor(), Bg: c.AsHexColor()}
	}
	th.Level = map[level.Level]color.FgBgTuple{
		level.Trace: styleForName("Charcoal Gray"),
		level.Debug: styleForName("Gray"),
		level.Info:  styleForName("White"),
		level.Warn:  styleForName("Yellow"),
		level.Error: styleForName("Red"),
	}
}

func (p generator) applyMessages(th *Theme) {
	var blend color.AdaptorColor
	switch color.Mode {
	case color.LightMode:
		blend = color.White
	case color.DarkMode:
		fallthrough
	default:
		blend = color.Black
	}
	styleForName := func(name string) color.FgBgTuple {
		bg := color.MustByName(name).Blend(blend, 90)
		return color.FgBgTuple{Fg: p.MessageColor.AsHexColor(), Bg: bg.AsHexColor()}
	}
	th.Message = map[level.Level]color.FgBgTuple{
		level.Trace: styleForName("Charcoal Gray"),
		level.Debug: styleForName("Gray"),
		level.Info:  styleForName("White"),
		level.Warn:  styleForName("Yellow"),
		level.Error: styleForName("Red"),
	}
}
