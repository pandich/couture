package theme

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"github.com/pandich/couture/model/level"
	color2 "github.com/pandich/couture/theme/color"
)

// SplitComplementaryGenerator ...
func SplitComplementaryGenerator(mode color2.ContrastPolarity, c color2.AdaptorColor) Generator {
	a := gamut.Complementary(c.AsImageColor())
	b := gamut.SplitComplementary(gamut.HueOffset(c.AsImageColor(), 120))
	x := gamut.Analogous(c.AsImageColor())
	return Generator{
		Mode:             mode,
		ApplicationColor: color2.FromImageColor(x[1]).AdjustConstrast(mode, color2.LessContrast, 0.2),
		TimestampColor:   color2.FromImageColor(a).AdjustConstrast(mode, color2.MoreContrast, 0.2),
		EntityColor:      c,
		MessageColor:     color2.FromImageColor(b[1]),
	}
}

// Generator ...
type Generator struct {
	Mode             color2.ContrastPolarity
	ApplicationColor color2.AdaptorColor
	TimestampColor   color2.AdaptorColor
	EntityColor      color2.AdaptorColor
	MessageColor     color2.AdaptorColor
}

// AsTheme ...
func (p Generator) AsTheme() Theme {
	th := Theme{}
	p.apply(&th)
	return th
}

func (p Generator) apply(th *Theme) {
	p.applySources(th)
	p.applyHeader(th)
	p.applyEntity(th)
	p.applyLevels(th)
	p.applyMessages(th)
}

func (p Generator) applySources(th *Theme) {
	const sourceColorCount = 180

	var cp = colorful.SoftPalette
	if gamut.Warm(p.EntityColor.AsImageColor()) {
		cp = colorful.WarmPalette
	}
	paletteColors, _ := cp(sourceColorCount)

	for _, source := range paletteColors {
		th.Source = append(
			th.Source,
			Style{
				Bg: color2.FromImageColor(source).AsHexColor(),
				Fg: color2.FromImageColor(source).Contrast().AsHexColor(),
			},
		)
	}
}

func (p Generator) applyHeader(th *Theme) {
	th.Application = Style{
		Fg: p.ApplicationColor.AsHexColor(),
		Bg: p.ApplicationColor.AdjustConstrast(p.Mode, color2.MoreContrast, 0.9).AsHexColor(),
	}
	th.Timestamp = Style{
		Fg: p.TimestampColor.AsHexColor(),
		Bg: p.TimestampColor.AdjustConstrast(p.Mode, color2.MoreContrast, 0.8).AsHexColor(),
	}
}

func (p Generator) applyEntity(th *Theme) {
	entityFg := p.EntityColor
	entityBg := entityFg.AdjustConstrast(p.Mode, color2.MoreContrast, 0.8)

	th.Entity = Style{
		Fg: entityFg.AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
	th.Context = Style{
		Fg: entityFg.AsHexColor(),
		Bg: entityFg.AdjustConstrast(p.Mode, color2.MoreContrast, 0.7).AsHexColor(),
	}

	th.Line = Style{
		Fg: entityFg.AdjustConstrast(p.Mode, color2.MoreContrast, 0.2).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}

	th.LineDelimiter = Style{
		Fg: entityFg.AdjustConstrast(p.Mode, color2.MoreContrast, 0.2).AdjustConstrast(p.Mode, color2.LessContrast, 0.4).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}

	actionFg := entityFg.AdjustConstrast(p.Mode, color2.LessContrast, 0.2)
	th.Action = Style{
		Fg: actionFg.AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
	th.ActionDelimiter = Style{
		Fg: actionFg.AdjustConstrast(p.Mode, color2.LessContrast, 0.4).AsHexColor(),
		Bg: entityBg.AsHexColor(),
	}
}

func (p Generator) applyLevels(th *Theme) {
	styleForName := func(name string) Style {
		c := color2.MustByName(name).
			Blend(p.EntityColor, 5)
		return Style{Fg: c.Contrast().AsHexColor(), Bg: c.AsHexColor()}
	}
	th.Level = map[level.Level]Style{
		level.Trace: styleForName("Charcoal Gray"),
		level.Debug: styleForName("Gray"),
		level.Info:  styleForName("White"),
		level.Warn:  styleForName("Yellow"),
		level.Error: styleForName("Red"),
	}
}

func (p Generator) applyMessages(th *Theme) {
	var blend = color2.Hex("#000000")
	if p.Mode == color2.LightMode {
		blend = color2.Hex("#ffffff")
	}
	styleForName := func(name string) Style {
		bg := color2.MustByName(name).Blend(blend, 90)
		return Style{Fg: p.MessageColor.AsHexColor(), Bg: bg.AsHexColor()}
	}
	th.Message = map[level.Level]Style{
		level.Trace: styleForName("Charcoal Gray"),
		level.Debug: styleForName("Gray"),
		level.Info:  styleForName("White"),
		level.Warn:  styleForName("Yellow"),
		level.Error: styleForName("Red"),
	}
}
