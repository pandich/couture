package model

import (
	"strconv"
	"strings"
)

// NoLineNumber ...
const (
	// NoLineNumber indicates no line number is present.
	NoLineNumber LineNumber = 0
)

// ClassName ...
type (
	// ApplicationName the name of an application.
	ApplicationName string
	// ClassName a class name.
	ClassName string
	// MethodName a method name.
	MethodName string
	// LineNumber  a line number.
	LineNumber uint64
)

// UnmarshalJSON ...
func (l *LineNumber) UnmarshalJSON(bytes []byte) error {
	s := strings.Trim(string(bytes), `"`)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*l = LineNumber(n)
	return nil
}

// Abbreviate ...
func (c ClassName) Abbreviate(maxWidth int) ClassName {
	var s = string(c)
	var pieces = strings.Split(s, ".")
	var l = len(s)
	var changed = true
	for l > maxWidth && changed {
		changed = false
		for i := 0; i < len(pieces)-1; i++ {
			if len(pieces[i]) > 1 {
				l -= len(pieces[i]) - 1
				pieces[i] = string(pieces[i][0])
				if l > maxWidth {
					changed = true
				}
			}
		}
	}
	changed = true
	for l > maxWidth && changed {
		changed = false
		if len(pieces) > 1 {
			l -= len(pieces[0]) + 1
			pieces = pieces[1:]
			if l > maxWidth {
				changed = true
			}
		}
	}
	if l > maxWidth {
		pieces[0] = pieces[0][len(pieces[0])-maxWidth:]
	}
	return ClassName(strings.Join(pieces, "."))
}
