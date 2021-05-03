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
