package layout

import (
	"fmt"
)

// NoPadding ...
var NoPadding = newPadding(0, 0)

// NoPaddingLayout ...
var NoPaddingLayout = ColumnLayout{Padding: NoPadding}

// Padding ,,,
//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Padding struct {
	Left  *uint `yaml:"left,omitempty"`
	Right *uint `yaml:"right,omitempty"`
}

// ColumnLayout ...
type ColumnLayout struct {
	Width   uint     `yaml:"width,omitempty"`
	Padding *Padding `yaml:"padding,omitempty"`
	Sigil   string   `yaml:"sigil,omitempty"`
	Align   string   `yaml:"align,omitempty"`
}

func (cl ColumnLayout) padding(padding *uint) string {
	var p = 1
	if padding != nil {
		p = int(*padding)
	}
	var s string
	for s = ""; len(s) < p; s += " " {
	}
	return s
}

func (cl ColumnLayout) leftPadding() string {
	width := cl.leftPaddingWidth()
	return cl.padding(&width)
}

func (cl ColumnLayout) rightPadding() string {
	if cl.Padding == nil || cl.Padding.Right == nil {
		return cl.padding(nil)
	}
	width := cl.rightPaddingWidth()
	return cl.padding(&width)
}

func newPadding(left uint, right uint) *Padding {
	return &Padding{Left: &left, Right: &right}
}

func (cl ColumnLayout) leftPaddingWidth() uint {
	p := cl.Padding
	if p == nil || p.Left == nil {
		return 1
	}
	return *p.Left
}

func (cl ColumnLayout) rightPaddingWidth() uint {
	p := cl.Padding
	if p == nil || p.Right == nil {
		return 1
	}
	return *p.Right
}

// EffectivePadding ...
func (cl ColumnLayout) EffectivePadding() Padding {
	return *newPadding(
		cl.leftPaddingWidth(),
		cl.rightPaddingWidth(),
	)
}

// Prefix ...
func (cl ColumnLayout) Prefix() string {
	var prefix = cl.leftPadding()
	if cl.Sigil != "" {
		prefix = cl.Sigil + prefix[1:]
	}
	return prefix
}

// Suffix ...
func (cl ColumnLayout) Suffix() string {
	return cl.rightPadding()
}

// Format ...
func (cl ColumnLayout) Format(name string) string {
	var sign = "-"
	if cl.Align == "right" {
		sign = ""
	}
	var valueFormat = "%s"
	if cl.Width > 0 {
		valueFormat = fmt.Sprintf("%%%[1]s%[2]d.%[2]ds", sign, cl.Width)
	}
	return fmt.Sprintf("{{%s%s%s}}::%s", cl.leftPadding(), valueFormat, cl.rightPadding(), name)
}
