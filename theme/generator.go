package theme

import (
	"github.com/muesli/gamut"
	"github.com/pandich/couture/model/level"
	"github.com/pandich/couture/theme/color"
	errors2 "github.com/pkg/errors"
)

// GenerateTheme ...
func GenerateTheme(base string, sourceStyle string) (*Theme, error) {
	ac, ok := color.ByName(base)
	if !ok {
		return nil, errors2.Errorf("invalid theme color: %s", base)
	}
	return splitComplementaryGenerator(ac, sourceStyle).asTheme(), nil
}

func splitComplementaryGenerator(baseColor color.AdaptorColor, sourceStyle string) generator {
	//nolint: gomnd
	return generator{
		SourceStyle: sourceStyle,
		ApplicationColor: baseColor.
			Analogous()[1].
			AdjustConstrast(color.LessNoticable, 0.2),
		TimestampColor: baseColor.
			Complementary().
			Monochromatic()[0xB0],
		EntityColor: baseColor,
		MessageColor: baseColor.
			Triadic()[1].
			AdjustConstrast(color.LessNoticable, 0.4).
			Lighter(0.2),
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
	th.Application = color.HexPair{
		Fg: p.ApplicationColor.AsHexColor(),
		Bg: p.ApplicationColor.AdjustConstrast(color.MoreNoticable, 0.9).AsHexColor(),
	}
	th.Timestamp = color.HexPair{
		Fg: p.TimestampColor.AsHexColor(),
		Bg: p.TimestampColor.AdjustConstrast(color.MoreNoticable, 0.8).AsHexColor(),
	}
}

func (p generator) applyEntity(th *Theme) {
	entityFg := p.EntityColor
	entityBg := entityFg.AdjustConstrast(color.MoreNoticable, 0.8)

	th.Entity = color.HexPair{
		Fg: entityFg.AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
	th.Context = color.HexPair{
		Fg: entityFg.AsHexColor(),
		Bg: entityFg.AdjustConstrast(color.MoreNoticable, 0.7).AsHexColor(),
	}

	th.Line = color.HexPair{
		Fg: entityFg.AdjustConstrast(color.MoreNoticable, 0.2).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}

	th.LineDelimiter = color.HexPair{
		Fg: entityFg.AdjustConstrast(color.MoreNoticable, 0.2).AdjustConstrast(color.LessNoticable, 0.4).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}

	actionFg := entityFg.AdjustConstrast(color.LessNoticable, 0.2)
	th.Action = color.HexPair{
		Fg: actionFg.AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
	th.ActionDelimiter = color.HexPair{
		Fg: actionFg.AdjustConstrast(color.LessNoticable, 0.4).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
}

func (p generator) applyLevels(th *Theme) {
	styleForName := func(name string) color.HexPair {
		c := color.MustByName(name).
			Blend(p.EntityColor, 5)
		return color.HexPair{Fg: c.Contrast().AsHexColor(), Bg: c.AsHexColor()}
	}
	th.Level = map[level.Level]color.HexPair{
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
	styleForName := func(name string) color.HexPair {
		bg := color.MustByName(name).Blend(blend, 90)
		return color.HexPair{Fg: p.MessageColor.AsHexColor(), Bg: bg.AsHexColor()}
	}
	th.Message = map[level.Level]color.HexPair{
		level.Trace: styleForName("Charcoal Gray"),
		level.Debug: styleForName("Gray"),
		level.Info:  styleForName("White"),
		level.Warn:  styleForName("Yellow"),
		level.Error: styleForName("Red"),
	}
}
