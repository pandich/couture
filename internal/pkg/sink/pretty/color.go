package pretty

import (
	"couture/internal/pkg/model"
	"couture/internal/pkg/sink"
	"crypto/sha256"
	"encoding/hex"
	"github.com/i582/cfmt/cmd/cfmt"
	"net/url"
)

type sourceColors chan string

func newSourceColors(theme Theme) sourceColors {
	return sink.NewColorCycle(theme.SourceColors, theme.DefaultColor)
}

func (p *sourceColors) registerSource(sourceURL model.SourceURL) {
	styleColor := <-*p
	cfmt.RegisterStyle(p.styleName(sourceURL),
		func(s string) string { return cfmt.Sprintf("{{/%-30.30s }}::"+styleColor, s) })
}

func (p *sourceColors) styleName(sourceURL model.SourceURL) string {
	u := url.URL(sourceURL)
	if s := u.String(); s != "" {
		hasher := sha256.New()
		hasher.Write([]byte(s))
		return "Source" + hex.EncodeToString(hasher.Sum(nil))
	}
	return "Default"
}
