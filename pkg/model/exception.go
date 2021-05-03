package model

import (
	"fmt"
)

// StackTrace ...
type (
	// StackTrace a stack trace.
	StackTrace string
	// Exception an exception.
	Exception struct {
		// StackTrace the full text of the stack trace.
		StackTrace StackTrace `json:"stack_trace"`
	}
)

// NewException ...
func NewException(err error) *Exception {
	if err == nil {
		return nil
	}
	return &Exception{StackTrace: StackTrace(fmt.Sprintf("%+v", err))}
}
