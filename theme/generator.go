package theme

import (
	"github.com/muesli/gamut"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/theme/color"
	errors2 "github.com/pkg/errors"
)

// GenerateTheme ...
func GenerateTheme(colorName string, sourceStyle string) (*Theme, error) {
	if s, ok := themeColors[colorName]; ok {
		colorName = s
	}
	ac, ok := color.ByName(colorName)
	if !ok {
		return nil, errors2.Errorf("invalid theme color: %s", colorName)
	}
	return splitComplementaryGenerator(ac, sourceStyle).asTheme(), nil
}

func splitComplementaryGenerator(baseColor color.AdaptorColor, sourceStyle string) generator {
	const triadicDirectionCutoff = 0.5 // 180º
	var messageColorTriadicIndex = 1
	if h, _, _ := baseColor.AsColorfulColor().Hsl(); h < triadicDirectionCutoff {
		messageColorTriadicIndex = 0
	}

	return generator{
		SourceStyle: sourceStyle,
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
	SourceStyle      string
}

func (p generator) asTheme() *Theme {
	th := Theme{SourceStyle: p.SourceStyle}
	p.apply(&th)
	return &th
}

func (p generator) apply(th *Theme) {
	p.applySources(th)
	p.applyHeader(th)
	p.applyEntity(th)
	p.applyLevels(th)
	p.applyMessages(th)
}

func (p generator) applySources(th *Theme) {
	blendColor := color.ByHex(th.Entity.Fg).Complementary()
	for _, paletteColor := range p.newSourcePalette(th) {
		th.Source = append(th.Source, paletteColor.
			Blend(blendColor, 10).
			AsHexPair())
	}
}

func (p generator) newSourcePalette(th *Theme) []color.AdaptorColor {
	const sourceColorCount = 100
	var generator gamut.ColorGenerator
	switch th.SourceStyle {
	case "warm":
		generator = gamut.WarmGenerator{}
	case "happy":
		generator = gamut.HappyGenerator{}
	case "similar":
		generator = gamut.SimilarHueGenerator{Color: gamut.Hex(th.Entity.Fg)}
	case "pastel":
		fallthrough
	default:
		generator = gamut.PastelGenerator{}
	}
	paletteColors, _ := gamut.Generate(sourceColorCount, generator)
	var out []color.AdaptorColor
	for _, pc := range paletteColors {
		out = append(out, color.ByImageColor(pc))
	}
	return out
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
