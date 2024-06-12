package event

import (
	"strconv"
	"strings"
)

const (
	// NoLineNumber indicates no line number is present.
	NoLineNumber Line = 0
)

// Entity ...
type (
	// Application the name of an application.
	Application string
	// Context a context name.
	Context string
	// Entity an entity name.
	Entity string
	// Action an action name.
	Action string
	// Line  a line number.
	Line uint64
)

// UnmarshalJSON ...
func (l *Line) UnmarshalJSON(bytes []byte) error {
	s := strings.Trim(string(bytes), `"`)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*l = Line(n)
	return nil
}

// Abbreviate shortens the name of an entity by trimming characters from the least-significant parts of the caller
// forward. A string is shortened until it meets the maxWidth, or no further abbreviation can be made.
// Example: 'com.example.of.couture.SomeClass' would be 'c.e.o.couture.SomeClass' with a value of 23.
// Any value between 0 and 18 would generate: 'c.e.o.c.SomeClass'.
func (c Entity) Abbreviate(maxWidth int) Entity {
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
	return Entity(strings.Join(pieces, "."))
}

// String ...
func (contextName Context) String() string {
	return string(contextName)
}
